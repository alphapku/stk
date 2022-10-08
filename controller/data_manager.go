package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	stk "StakeBackendGoTest/internal/entity/stake"
	resp "StakeBackendGoTest/internal/response"
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

	// stakePositions map[string]map[string]*stk.Position
	// stakePrices    map[string]map[string]*stk.Price

	// TODO, we do not save user info, so here let's assume the sytem is running for a dedicated user only.
	// Use the structures above with the outer map's key as account ID if running to support multiple users
	positions map[string]*resp.StakePosition

	internalPositions map[string]*stk.Position
	internalPrices    map[string]*stk.Price
}

func NewDataManager() *DataManager {
	return &DataManager{
		positions: make(map[string]*resp.StakePosition),

		internalPositions: make(map[string]*stk.Position),
		internalPrices:    make(map[string]*stk.Price),
	}
}

func (d *DataManager) OnMessage(msg interface{}) {
	switch m := msg.(type) {
	case []*stk.Position:
		d.onMarketPositions(m)
	case []*stk.Price:
		d.onMarketPrices(m)
	default:
		log.Logger.Warn("unexpected msg", zap.String("type", fmt.Sprintf("%T", m)))
	}
}

func (d *DataManager) onMarketPositions(positions []*stk.Position) {
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

func (d *DataManager) onMarketPrices(prices []*stk.Price) {
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

func (d *DataManager) calcStakePosition(pos *stk.Position, prx *stk.Price) {
	availForTrddingQty := pos.AvailableForTradingQty.StringFixed(satoshiDecimalPlaces)

	mktPrice := calcMarketPrice(prx)
	mktPriceStr := mktPrice.StringFixed(satoshiDecimalPlaces)
	mktValue := pos.OpenQty.Mul(mktPrice)
	mktValueStr := mktValue.StringFixed(satoshiDecimalPlaces)

	priorValue := pos.OpenQty.Mul(prx.PriorClose)
	priorClose := prx.PriorClose.StringFixed(satoshiDecimalPlaces)

	dayPNL := mktValue.Sub(priorValue)
	dayPNLStr := mktValue.Sub(priorValue).StringFixed(satoshiDecimalPlaces)
	dayPNLPCT := dayPNL.Div(priorValue).Mul(pctMultiplier)
	dayPNLPCTStr := dayPNLPCT.StringFixed(pctDecimalPlaces)

	totalPNL := mktValue.Sub(pos.Cost)
	totalPNLStr := mktValue.Sub(pos.Cost).StringFixed(satoshiDecimalPlaces)
	totalPNLPCT := totalPNL.Div(pos.Cost).Mul(pctMultiplier)
	totalPNLPCTStr := totalPNLPCT.StringFixed(pctDecimalPlaces)

	if stkPos, ok := d.positions[prx.Symbol]; ok {
		// update the existed
		// stkPos.OpenQty = pos.OpenQty.String()
		stkPos.AvailableForTradingQty = availForTrddingQty
		// stkPos.AveragePrice = pos.AveragePrice.String()
		stkPos.MarketValue = mktValueStr
		stkPos.MarketPrice = mktPriceStr
		stkPos.PriorClose = priorClose
		stkPos.DayProfitOrLoss = dayPNLStr
		stkPos.DayProfitOrLossPercent = dayPNLPCTStr
		stkPos.TotalProfitOrLoss = totalPNLStr
		stkPos.TotalProfitOrLossPercent = totalPNLPCTStr
		a, _ := json.Marshal(stkPos)
		log.Logger.Debug("updated", zap.String("position", string(a)))
	} else {
		// create a new one
		d.positions[prx.Symbol] = &resp.StakePosition{
			Symbol:                   prx.Symbol,
			Name:                     pos.Name,
			OpenQty:                  pos.OpenQty.StringFixed(satoshiDecimalPlaces),
			AvailableForTradingQty:   availForTrddingQty,
			AveragePrice:             pos.AveragePrice.StringFixed(satoshiDecimalPlaces),
			MarketValue:              mktValueStr,
			MarketPrice:              mktPriceStr,
			PriorClose:               priorClose,
			DayProfitOrLoss:          dayPNLStr,
			DayProfitOrLossPercent:   dayPNLPCTStr,
			TotalProfitOrLoss:        totalPNLStr,
			TotalProfitOrLossPercent: totalPNLPCTStr,
		}
		a, _ := json.Marshal(d.positions[prx.Symbol])
		log.Logger.Debug("created", zap.String("position", string(a)))
	}
}

// calcMarketPrice returns the preferred market price from our agreement.
// Here, we use the lastTrade
func calcMarketPrice(prx *stk.Price) decimal.Decimal {
	return prx.LastTrade
}

func (d *DataManager) DoEquityPositions(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, []string{"hello", "world!"})
}
