package stake

import "github.com/shopspring/decimal"

type Price struct {
	Symbol     string
	LastTrade  decimal.Decimal
	Bid        decimal.Decimal
	Ask        decimal.Decimal
	PriorClose decimal.Decimal
}

func (p Price) Equal(other Price) bool {
	return p.Symbol == other.Symbol &&
		p.LastTrade.Equal(other.LastTrade) &&
		p.Bid.Equal(other.Bid) &&
		other.Ask.Equal(p.Ask) &&
		p.PriorClose.Equal(other.PriorClose)
}
