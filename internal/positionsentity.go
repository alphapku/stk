package internal

type Positions struct {
	Positions []Position `json:"equityPositions"`
}

type Position struct {
	Security            string  `json:"security"`
	SecurityDescription string  `json:"securityDescription"`
	Cost                float64 `json:"cost"`
	AveragePrice        float64 `json:"averagePrice"`
	AvailableUnits      int     `json:"availableUnits"`
	PortfolioUnits      int     `json:"portfolioUnits"`
}
