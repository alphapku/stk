package stake

import "github.com/shopspring/decimal"

type InternalPosition struct {
	Symbol                 string
	Name                   string
	AvailableForTradingQty decimal.Decimal
	AveragePrice           decimal.Decimal
	Cost                   decimal.Decimal
	OpenQty                decimal.Decimal // the current volume
}

func (p InternalPosition) Equal(other InternalPosition) bool {
	return p.Symbol == other.Symbol &&
		p.Name == other.Name &&
		p.AvailableForTradingQty.Equal(other.AvailableForTradingQty) &&
		p.AveragePrice.Equal(other.AveragePrice) &&
		p.Cost.Equal(other.Cost) &&
		p.OpenQty.Equal(other.OpenQty)
}
