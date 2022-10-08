package response

type StakePositions struct {
	StakePositions []StakePosition
}

type StakePosition struct {
	Symbol                   string `json:"symbol"`
	Name                     string `json:"name"`
	OpenQty                  string `json:"openQty"`
	AvailableForTradingQty   string `json:"availableForTradingQty"`
	AveragePrice             string `json:"averagePrice"`
	MarketValue              string `json:"marketValue"` // = OpenQty * MarketPrice
	MarketPrice              string `json:"marketPrice"` // lastTrade price
	PriorClose               string `json:"priorClose"`
	DayProfitOrLoss          string `json:"dayProfitOrLoss"`
	DayProfitOrLossPercent   string `json:"dayProfitOrLossPercent"`
	TotalProfitOrLoss        string `json:"totalProfitOrLoss"`
	TotalProfitOrLossPercent string `json:"totalProfitOrLossPercent"`
}
