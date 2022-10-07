package mock

import (
	"errors"

	"github.com/shopspring/decimal"

	mk "StakeBackendGoTest/internal/entity/mock"
	stk "StakeBackendGoTest/internal/entity/stake"
)

var (
	ErrSymbolNotFoundInPrices = errors.New("symbol not found in prices")
)

func ToStakePosition(p *mk.Position) *stk.Position {
	return &stk.Position{
		Symbol:                 p.Security,
		Name:                   p.SecurityDescription,
		OpenQty:                p.Cost.Div(p.AveragePrice),
		AvailableForTradingQty: decimal.NewFromInt(int64(p.AvailableUnits)),
		AveragePrice:           p.AveragePrice,
		Cost:                   p.Cost,
	}
}

func ToStakePrice(p *mk.Price) *stk.Price {
	return &stk.Price{
		Symbol:     p.Symbol,
		LastTrade:  p.LastTrade,
		Bid:        p.Bid,
		Ask:        p.Ask,
		PriorClose: p.PriorClose,
	}
}
