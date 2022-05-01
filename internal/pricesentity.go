package internal

type Prices struct {
	Prices []Price `json:"priceData"`
}

type Price struct {
	MarketStatus string  `json:"marketStatus"`
	Symbol       string  `json:"symbol"`
	LastTrade    float64 `json:"lastTrade"`
	Bid          float64 `json:"bid"`
	Ask          float64 `json:"ask"`
	PriorClose   float64 `json:"priorClose"`
}
