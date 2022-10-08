package stake

import "github.com/shopspring/decimal"

type InternalPrice struct {
	Symbol     string
	LastTrade  decimal.Decimal
	Bid        decimal.Decimal
	Ask        decimal.Decimal
	PriorClose decimal.Decimal
}

func (p InternalPrice) Equal(other InternalPrice) bool {
	return p.Symbol == other.Symbol &&
		p.LastTrade.Equal(other.LastTrade) &&
		p.Bid.Equal(other.Bid) &&
		other.Ask.Equal(p.Ask) &&
		p.PriorClose.Equal(other.PriorClose)
}
