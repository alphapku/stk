package stake

import "github.com/shopspring/decimal"

type Price struct {
	Symbol     string
	LastTrade  decimal.Decimal
	Bid        decimal.Decimal
	Ask        decimal.Decimal
	PriorClose decimal.Decimal
}
