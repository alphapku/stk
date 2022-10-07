package mock

import (
	"errors"
	"strconv"

	"github.com/shopspring/decimal"

	em "StakeBackendGoTest/internal/entity/mock"
	resp "StakeBackendGoTest/internal/response"
)

var (
	ErrSymbolNotFoundInPrices = errors.New("symbol not found in prices")
)

func ToPositions(positions *em.Positions, prices *em.Prices) (*resp.StakePositions, error) {
	stakes := make([]*resp.StakePosition, len(positions.Positions))

	// turn prices to be a map based on symbol
	prx_map := make(map[string]*em.Price, len(prices.Prices))
	for _, prx := range prices.Prices {
		prx_map[prx.Symbol] = prx
	}

	for i, p := range positions.Positions {
		prx, ok := prx_map[p.Security]
		if !ok {
			return nil, ErrSymbolNotFoundInPrices
		}

		sp, err := toPosition(p, prx)
		if err != nil {
			return nil, err
		}
		stakes[i] = sp
	}

	return &resp.StakePositions{
		StakePositions: stakes,
	}, nil
}

func toPosition(p *em.Position, prx *em.Price) (*resp.StakePosition, error) {
	qty := p.Cost.Div(p.AveragePrice)

	open_value := p.Cost

	prio_value := qty.Mul(prx.PriorClose)

	mkt_price := calcMarketPrice(prx)
	mkt_value := qty.Mul(mkt_price)

	day_pnl := mkt_value.Sub(prio_value)
	day_pnl_pct := day_pnl.Div(prio_value)

	total_pnl := mkt_value.Sub(open_value)
	total_pnl_pct := total_pnl.Div(open_value)

	// TODO: return an error if something goes wrong after validating input
	return &resp.StakePosition{
		Symbol:                   p.Security,
		Name:                     p.SecurityDescription,
		OpenQty:                  qty.String(),
		AvailableForTradingQty:   strconv.Itoa(p.AvailableUnits),
		AveragePrice:             p.AveragePrice.String(),
		MarketValue:              mkt_value.String(),
		MarketPrice:              mkt_price.String(),
		PriorClose:               prx.PriorClose.String(),
		DayProfitOrLoss:          day_pnl.String(),
		DayProfitOrLossPercent:   day_pnl_pct.String(),
		TotalProfitOrLoss:        total_pnl.String(),
		TotalProfitOrLossPercent: total_pnl_pct.String(),
	}, nil
}

func calcMarketPrice(prx *em.Price) decimal.Decimal {
	return prx.LastTrade
}
