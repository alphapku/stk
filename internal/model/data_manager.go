package model

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	resp "StakeBackendGoTest/api/response"
	stk "StakeBackendGoTest/internal/entity/stake"
	log "StakeBackendGoTest/pkg/log"
)

const (
	pctDecimalPlaces     = 2 // to show pct as xx.xx
	satoshiDecimalPlaces = 8
)

var (
	pctMultiplier = decimal.NewFromInt(100)
)

type DataManager struct {
	// positions map[string]map[string]*resp.StakePosition

	// stakePositions map[string]map[string]*stk.InternalPosition
	// stakePrices    map[string]map[string]*stk.InternalPrice

	// TODO, we do not save user info, so here let's assume the sytem is running for a dedicated user only.
	// Use the structures above with the outer map's key as account ID if running to support multiple users
	// By adding account info in internal positions, it could be scale easily
	positions []*resp.StakePosition

	internalPositions map[string]*stk.InternalPosition
	internalPrices    map[string]*stk.InternalPrice
}

func NewDataManager() *DataManager {
	return &DataManager{
		positions: make([]*resp.StakePosition, 0),

		internalPositions: make(map[string]*stk.InternalPosition),
		internalPrices:    make(map[string]*stk.InternalPrice),
	}
}

func (d *DataManager) OnMessage(msg interface{}) {
	switch m := msg.(type) {
	case []*stk.InternalPosition:
		d.onMarketPositions(m)
	case []*stk.InternalPrice:
		d.onMarketPrices(m)
	default:
		log.Logger.Warn("unexpected msg", zap.String("type", fmt.Sprintf("%T", m)))
	}
}

func (d *DataManager) onMarketPositions(positions []*stk.InternalPosition) {
	log.Logger.Debug("position(s) received", zap.Int("count", len(positions)))
	for _, pos := range positions {
		// save the position info, as we need them to calculate StakePosition when prices are updated
		d.internalPositions[pos.Symbol] = pos

		// update positions by combining internal price we have to calculate StakePosition
		if prx, ok := d.internalPrices[pos.Symbol]; ok {
			d.calcStakePosition(pos, prx)
		}
	}
}

func (d *DataManager) onMarketPrices(prices []*stk.InternalPrice) {
	log.Logger.Debug("price(s) received", zap.Int("count", len(prices)))
	for _, prx := range prices {
		// save the price info, as we need them to calculate StakePosition when positions are updated
		d.internalPrices[prx.Symbol] = prx

		// update positions by combining internal position we have to calculate StakePosition
		if pos, ok := d.internalPositions[prx.Symbol]; ok {
			d.calcStakePosition(pos, prx)
		}
	}
}

func (d *DataManager) calcStakePosition(pos *stk.InternalPosition, prx *stk.InternalPrice) {
	availForTrddingQty := pos.AvailableForTradingQty.StringFixed(satoshiDecimalPlaces)

	mktPrice := calcMarketPrice(prx)
	mktPriceStr := mktPrice.StringFixed(satoshiDecimalPlaces)
	mktValue := pos.OpenQty.Mul(mktPrice)
	mktValueStr := mktValue.StringFixed(satoshiDecimalPlaces)

	priorValue := prx.PriorClose.Mul(pos.AvailableForTradingQty)
	priorCloseStr := prx.PriorClose.StringFixed(satoshiDecimalPlaces)

	dayPNL := mktValue.Sub(priorValue)
	dayPNLStr := dayPNL.StringFixed(satoshiDecimalPlaces)
	dayPNLPCT := dayPNL.Div(priorValue).Mul(pctMultiplier)
	dayPNLPCTStr := dayPNLPCT.StringFixed(pctDecimalPlaces)

	totalPNL := mktValue.Sub(pos.Cost)
	totalPNLStr := totalPNL.StringFixed(satoshiDecimalPlaces)
	totalPNLPCT := totalPNL.Div(pos.Cost).Mul(pctMultiplier)
	totalPNLPCTStr := totalPNLPCT.StringFixed(pctDecimalPlaces)

	updated := false
	for _, p := range d.positions {
		if p.Symbol == prx.Symbol {
			updated = true

			// update the existed
			p.AvailableForTradingQty = availForTrddingQty
			p.MarketValue = mktValueStr
			p.MarketPrice = mktPriceStr
			p.DayProfitOrLoss = dayPNLStr
			p.DayProfitOrLossPercent = dayPNLPCTStr
			p.TotalProfitOrLoss = totalPNLStr
			p.TotalProfitOrLossPercent = totalPNLPCTStr
		}
	}

	if !updated {
		// create a new one
		d.positions = append(d.positions, &resp.StakePosition{
			Symbol:                   prx.Symbol,
			Name:                     pos.Name,
			OpenQty:                  pos.OpenQty.StringFixed(satoshiDecimalPlaces),
			AvailableForTradingQty:   availForTrddingQty,
			AveragePrice:             pos.AveragePrice.StringFixed(satoshiDecimalPlaces),
			MarketValue:              mktValueStr,
			MarketPrice:              mktPriceStr,
			PriorClose:               priorCloseStr,
			DayProfitOrLoss:          dayPNLStr,
			DayProfitOrLossPercent:   dayPNLPCTStr,
			TotalProfitOrLoss:        totalPNLStr,
			TotalProfitOrLossPercent: totalPNLPCTStr,
		})
	}
}

func (d *DataManager) DoEquityPositions(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, d.positions)
}

// Reset clears all the data DataManager caches
func (d *DataManager) Reset() {
	d.positions = make([]*resp.StakePosition, 0)
	d.internalPositions = make(map[string]*stk.InternalPosition)
	d.internalPrices = make(map[string]*stk.InternalPrice)
}

// calcMarketPrice returns the preferred market price from our agreement.
// Here, we use the lastTrade
func calcMarketPrice(prx *stk.InternalPrice) decimal.Decimal {
	return prx.LastTrade
}
