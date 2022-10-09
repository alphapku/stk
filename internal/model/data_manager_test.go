package model

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	resp "StakeBackendGoTest/api/response"
	intl "StakeBackendGoTest/internal/entity/stake"
	def "StakeBackendGoTest/pkg/const"
	log "StakeBackendGoTest/pkg/log"
)

type dataManagerTestSuite struct {
	suite.Suite

	d *DataManager
}

func (s *dataManagerTestSuite) SetupSuite() {
	_ = log.Init(def.DevMode)

	s.d = NewDataManager()
}

func (s *dataManagerTestSuite) TestDataManager() {
	var (
		symbol = "APT.ASX"
		name   = "Afterpay Limited"
	)

	// case #1
	// only internal position, no price, so the StakePositions will be delivered but some fields will be "N/A" as they depend on prices
	p := s.OnInternalPosition(symbol, name)
	s.Equal(len(s.d.positions), 0)
	s.Equal(len(s.d.internalPositions), 1)
	intlPosition, ok := s.d.internalPositions[symbol]
	s.True(ok)
	s.True(p.Equal(*intlPosition))

	stkPositions := s.d.getStakePositions()
	s.Equal(len(stkPositions), 1)
	pos := stkPositions[0]

	s.Equal(pos.Symbol, p.Symbol)
	s.Equal(pos.Name, p.Name)
	s.Equal(pos.OpenQTY, p.OpenQTY.StringFixed(satoshiDecimalPlaces))
	s.Equal(pos.AvailableForTradingQTY, p.AvailableForTradingQTY.StringFixed(satoshiDecimalPlaces))
	s.Equal(pos.AveragePrice, p.AveragePrice.StringFixed(satoshiDecimalPlaces))

	s.Equal(pos.MarketValue, na)
	s.Equal(pos.MarketPrice, na)
	s.Equal(pos.PriorClose, na)
	s.Equal(pos.DayProfitOrLoss, na)
	s.Equal(pos.DayProfitOrLossPercent, na)
	s.Equal(pos.TotalProfitOrLoss, na)
	s.Equal(pos.TotalProfitOrLossPercent, na)

	s.Assert().Equal(s.d.internalPositions[p.Symbol].Equal(*p), true)

	// case #2
	// only internal price, no position, so the StakePositions should not be created
	s.d.Reset()

	stkPrice := s.OnInternalPrice(symbol)
	s.Equal(len(s.d.positions), 0)
	s.Equal(len(s.d.internalPrices), 1)
	s.Assert().Equal(s.d.internalPrices[stkPrice.Symbol].Equal(*stkPrice), true)

	// case #3
	// fill position we are supposed to have Response position now
	p = s.OnInternalPosition(symbol, name)

	mktPrice := calcMarketPrice(stkPrice)
	curVolume := p.Cost.Div(p.AveragePrice)
	mktValue := mktPrice.Mul(curVolume)

	priorValue := stkPrice.PriorClose.Mul(p.AvailableForTradingQTY)
	priorClose := stkPrice.PriorClose

	dayPNL := mktValue.Sub(priorValue)
	dayPNLPCT := dayPNL.Div(priorValue).Mul(pctMultiplier)

	totalPNL := mktValue.Sub(p.Cost)
	totalPNLPCT := totalPNL.Div(p.Cost).Mul(pctMultiplier)

	expectedRespPosition := resp.StakePosition{
		Symbol:                   symbol,
		Name:                     name,
		OpenQTY:                  p.OpenQTY.StringFixed(satoshiDecimalPlaces),
		AvailableForTradingQTY:   p.AvailableForTradingQTY.StringFixed(satoshiDecimalPlaces),
		AveragePrice:             p.AveragePrice.StringFixed(satoshiDecimalPlaces),
		MarketValue:              mktValue.StringFixed(satoshiDecimalPlaces),
		MarketPrice:              mktPrice.StringFixed(satoshiDecimalPlaces),
		PriorClose:               priorClose.StringFixed(satoshiDecimalPlaces),
		DayProfitOrLoss:          dayPNL.StringFixed(satoshiDecimalPlaces),
		DayProfitOrLossPercent:   dayPNLPCT.StringFixed(pctDecimalPlaces),
		TotalProfitOrLoss:        totalPNL.StringFixed(satoshiDecimalPlaces),
		TotalProfitOrLossPercent: totalPNLPCT.StringFixed(pctDecimalPlaces),
	}

	s.Equal(len(s.d.positions), 1)

	s.Equal(expectedRespPosition, *s.d.positions[0])
}

func (s *dataManagerTestSuite) OnInternalPosition(symbol, name string) *intl.InternalPosition {
	p := &intl.InternalPosition{
		Symbol:                 symbol,
		Name:                   name,
		OpenQTY:                decimal.NewFromFloat(10.0000),
		AvailableForTradingQTY: decimal.NewFromFloat(10.0000),
		AveragePrice:           decimal.NewFromFloat(102.5000),
		Cost:                   decimal.NewFromFloat(1025.0000),
	}

	s.d.onMarketPositions([]*intl.InternalPosition{p})

	return p
}

func (s *dataManagerTestSuite) OnInternalPrice(symbol string) *intl.InternalPrice {
	p := &intl.InternalPrice{
		Symbol:     symbol,
		LastTrade:  decimal.NewFromFloat(114.9800),
		Bid:        decimal.NewFromFloat(114.98),
		Ask:        decimal.NewFromFloat(114.99),
		PriorClose: decimal.NewFromFloat(119.8700),
	}
	s.d.onMarketPrices([]*intl.InternalPrice{p})

	return p
}

func TestDataManagerTestSuite(t *testing.T) {
	suite.Run(t, new(dataManagerTestSuite))
}
