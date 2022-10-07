package mock

import (
	"github.com/shopspring/decimal"
)

type Positions struct {
	Positions []*Position `json:"equityPositions"`
}

type Position struct {
	Security            string          `json:"security"`
	SecurityDescription string          `json:"securityDescription"`
	Cost                decimal.Decimal `json:"cost"`
	AveragePrice        decimal.Decimal `json:"averagePrice"`
	AvailableUnits      int             `json:"availableUnits"`
	PortfolioUnits      int             `json:"portfolioUnits"`
}
