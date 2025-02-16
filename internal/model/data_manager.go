package model

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	resp "StakeBackendGoTest/api/response"
	intl "StakeBackendGoTest/internal/entity/stake"
	log "StakeBackendGoTest/pkg/log"
)

const (
	pctDecimalPlaces     = 2 // to show pct as xx.xx
	satoshiDecimalPlaces = 8

	na = "N/A"
)

var (
	pctMultiplier = decimal.NewFromInt(100)
)

type DataManager struct {
	// positions map[string]map[string]*resp.StakePosition

	// stakePositions map[string]map[string]*intl.InternalPosition
	// stakePrices    map[string]map[string]*intl.InternalPrice

	// TODO, we do not save user info, so here let's assume the sytem is running for a dedicated user only.
	// Use the structures above with the outer map's key as account ID if running to support multiple users
	// By adding account info in internal positions, it could be scale easily
	positions []*resp.StakePosition

	internalPositions map[string]*intl.InternalPosition
	internalPrices    map[string]*intl.InternalPrice
}

func NewDataManager() *DataManager {
	return &DataManager{
		positions: make([]*resp.StakePosition, 0),

		internalPositions: make(map[string]*intl.InternalPosition),
		internalPrices:    make(map[string]*intl.InternalPrice),
	}
}

func (d *DataManager) OnMessage(msg interface{}) {
	switch m := msg.(type) {
	case []*intl.InternalPosition:
		d.onMarketPositions(m)
	case []*intl.InternalPrice:
		d.onMarketPrices(m)
	default:
		log.Logger.Warn("unexpected msg", zap.String("type", fmt.Sprintf("%T", m)))
	}
}

func (d *DataManager) onMarketPositions(positions []*intl.InternalPosition) {
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

func (d *DataManager) onMarketPrices(prices []*intl.InternalPrice) {
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

func (d *DataManager) calcStakePosition(pos *intl.InternalPosition, prx *intl.InternalPrice) {
	availForTrddingQTY := pos.AvailableForTradingQTY.StringFixed(satoshiDecimalPlaces)

	mktPrice := calcMarketPrice(prx)
	mktPriceStr := mktPrice.StringFixed(satoshiDecimalPlaces)
	mktValue := pos.OpenQTY.Mul(mktPrice)
	mktValueStr := mktValue.StringFixed(satoshiDecimalPlaces)

	priorValue := prx.PriorClose.Mul(pos.AvailableForTradingQTY)
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
			p.AvailableForTradingQTY = availForTrddingQTY
			p.MarketValue = mktValueStr
			p.MarketPrice = mktPriceStr
			p.DayProfitOrLoss = dayPNLStr
			p.DayProfitOrLossPercent = dayPNLPCTStr
			p.TotalProfitOrLoss = totalPNLStr
			p.TotalProfitOrLossPercent = totalPNLPCTStr

			break
		}
	}

	if !updated {
		// create a new one
		d.positions = append(d.positions, &resp.StakePosition{
			Symbol:                   prx.Symbol,
			Name:                     pos.Name,
			OpenQTY:                  pos.OpenQTY.StringFixed(satoshiDecimalPlaces),
			AvailableForTradingQTY:   availForTrddingQTY,
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
	ctx.JSON(http.StatusOK, resp.Response{
		ErrorCode: 0,
		Data: resp.StakePositions{
			StakePositions: d.getStakePositions(),
		},
	})
}

// getStakePositions returns Stake positions.
// NB, it's not necessary to be from the real-time data as normally we would use cache
func (d *DataManager) getStakePositions() []*resp.StakePosition {
	r := []*resp.StakePosition{}

	for k, v := range d.internalPositions {
		updated := false
		for _, p := range d.positions {
			if p.Symbol == k {
				r = append(r, p)
				updated = true
				break
			}
		}

		if !updated {
			// create an StakePosition from internal position without PNL info
			n := &resp.StakePosition{
				Symbol:                   v.Symbol,
				Name:                     v.Name,
				OpenQTY:                  v.OpenQTY.StringFixed(satoshiDecimalPlaces),
				AvailableForTradingQTY:   v.AvailableForTradingQTY.StringFixed(satoshiDecimalPlaces),
				AveragePrice:             v.AveragePrice.StringFixed(satoshiDecimalPlaces),
				MarketValue:              na,
				MarketPrice:              na,
				PriorClose:               na,
				DayProfitOrLoss:          na,
				DayProfitOrLossPercent:   na,
				TotalProfitOrLoss:        na,
				TotalProfitOrLossPercent: na,
			}

			r = append(r, n)
		}
	}

	return r
}

// Reset clears all the data DataManager caches
func (d *DataManager) Reset() {
	d.positions = make([]*resp.StakePosition, 0)
	d.internalPositions = make(map[string]*intl.InternalPosition)
	d.internalPrices = make(map[string]*intl.InternalPrice)
}

// calcMarketPrice returns the preferred market price from our agreement.
// Here, we use the lastTrade
func calcMarketPrice(prx *intl.InternalPrice) decimal.Decimal {
	return prx.LastTrade
}
