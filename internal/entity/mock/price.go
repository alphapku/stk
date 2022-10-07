package mock

import (
	"github.com/shopspring/decimal"
)

type Prices struct {
	Prices []*Price `json:"priceData"`
}

type Price struct {
	MarketStatus string          `json:"marketStatus"`
	Symbol       string          `json:"symbol"`
	LastTrade    decimal.Decimal `json:"lastTrade"`
	Bid          decimal.Decimal `json:"bid"`
	Ask          decimal.Decimal `json:"ask"`
	PriorClose   decimal.Decimal `json:"priorClose"`
}
