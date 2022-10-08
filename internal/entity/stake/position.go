package stake

import "github.com/shopspring/decimal"

type InternalPosition struct {
	Symbol                 string
	Name                   string
	AvailableForTradingQTY decimal.Decimal
	AveragePrice           decimal.Decimal
	Cost                   decimal.Decimal
	OpenQTY                decimal.Decimal // the current volume
}

func (p InternalPosition) Equal(other InternalPosition) bool {
	return p.Symbol == other.Symbol &&
		p.Name == other.Name &&
		p.AvailableForTradingQTY.Equal(other.AvailableForTradingQTY) &&
		p.AveragePrice.Equal(other.AveragePrice) &&
		p.Cost.Equal(other.Cost) &&
		p.OpenQTY.Equal(other.OpenQTY)
}
