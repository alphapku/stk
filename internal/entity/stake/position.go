package stake

import "github.com/shopspring/decimal"

type Position struct {
	Symbol                 string
	Name                   string
	OpenQty                decimal.Decimal
	AvailableForTradingQty decimal.Decimal
	AveragePrice           decimal.Decimal
	Cost                   decimal.Decimal
}
