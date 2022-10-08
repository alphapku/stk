package response

type StakePositions struct {
	StakePositions []*StakePosition `json:"equityPositions"`
}

type StakePosition struct {
	Symbol                   string `json:"symbol"`
	Name                     string `json:"name"`
	OpenQTY                  string `json:"openQty"`
	AvailableForTradingQTY   string `json:"availableForTradingQty"`
	AveragePrice             string `json:"averagePrice"`
	MarketValue              string `json:"marketValue"`
	MarketPrice              string `json:"marketPrice"`
	PriorClose               string `json:"priorClose"`
	DayProfitOrLoss          string `json:"dayProfitOrLoss"`
	DayProfitOrLossPercent   string `json:"dayProfitOrLossPercent"`
	TotalProfitOrLoss        string `json:"totalProfitOrLoss"`
	TotalProfitOrLossPercent string `json:"totalProfitOrLossPercent"`
}
