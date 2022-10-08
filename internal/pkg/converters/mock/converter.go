package mock

import (
	"errors"

	"github.com/shopspring/decimal"

	mk "StakeBackendGoTest/internal/entity/mock"
	intl "StakeBackendGoTest/internal/entity/stake"
)

var (
	ErrZeroAveragePrice = errors.New("average is zero")
)

func ToStakePosition(p *mk.Position) (*intl.InternalPosition, error) {
	if p.AveragePrice.IsZero() {
		return nil, ErrZeroAveragePrice
	}

	return &intl.InternalPosition{
		Symbol:                 p.Security,
		Name:                   p.SecurityDescription,
		AvailableForTradingQTY: decimal.NewFromInt(int64(p.AvailableUnits)),
		AveragePrice:           p.AveragePrice,
		Cost:                   p.Cost,
		OpenQTY:                decimal.NewFromInt(int64(p.PortfolioUnits)),
	}, nil
}

func ToStakePrice(p *mk.Price) *intl.InternalPrice {
	return &intl.InternalPrice{
		Symbol:     p.Symbol,
		LastTrade:  p.LastTrade,
		Bid:        p.Bid,
		Ask:        p.Ask,
		PriorClose: p.PriorClose,
	}
}
