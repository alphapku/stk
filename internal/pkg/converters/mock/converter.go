package mock

import (
	"errors"

	"github.com/shopspring/decimal"

	mk "StakeBackendGoTest/internal/entity/mock"
	stk "StakeBackendGoTest/internal/entity/stake"
)

var (
	ErrZeroAveragePrice = errors.New("average is zero")
)

func ToStakePosition(p *mk.Position) (*stk.InternalPosition, error) {
	if p.AveragePrice.IsZero() {
		return nil, ErrZeroAveragePrice
	}

	return &stk.InternalPosition{
		Symbol:                 p.Security,
		Name:                   p.SecurityDescription,
		AvailableForTradingQty: decimal.NewFromInt(int64(p.AvailableUnits)),
		AveragePrice:           p.AveragePrice,
		Cost:                   p.Cost,
		OpenQty:                decimal.NewFromInt(int64(p.PortfolioUnits)),
	}, nil
}

func ToStakePrice(p *mk.Price) *stk.InternalPrice {
	return &stk.InternalPrice{
		Symbol:     p.Symbol,
		LastTrade:  p.LastTrade,
		Bid:        p.Bid,
		Ask:        p.Ask,
		PriorClose: p.PriorClose,
	}
}
