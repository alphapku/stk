package controller

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	stk "StakeBackendGoTest/internal/entity/stake"
	resp "StakeBackendGoTest/internal/response"
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
	// only Stake position, no price, so the Response position should not be created
	stkPosition := s.OnStkPosition(symbol, name)
	s.Equal(len(s.d.positions), 0)
	s.Equal(len(s.d.internalPositions), 1)
	s.Assert().Equal(s.d.internalPositions[stkPosition.Symbol].Equal(*stkPosition), true)

	// case #2
	// only Stake price, no position, , so the Response position should not be created
	s.d.Reset()

	stkPrice := s.OnStkPrice(symbol)
	s.Equal(len(s.d.positions), 0)
	s.Equal(len(s.d.internalPrices), 1)
	s.Assert().Equal(s.d.internalPrices[stkPrice.Symbol].Equal(*stkPrice), true)

	// case #3
	// fill position we are supposed to have Response position now
	stkPosition = s.OnStkPosition(symbol, name)

	mktPrice := calcMarketPrice(stkPrice)
	curVolume := stkPosition.Cost.Div(stkPosition.AveragePrice)
	mktValue := mktPrice.Mul(curVolume)

	priorValue := stkPrice.PriorClose.Mul(stkPosition.AvailableForTradingQty)
	priorClose := stkPrice.PriorClose

	dayPNL := mktValue.Sub(priorValue)
	dayPNLPCT := dayPNL.Div(priorValue).Mul(pctMultiplier)

	totalPNL := mktValue.Sub(stkPosition.Cost)
	totalPNLPCT := totalPNL.Div(stkPosition.Cost).Mul(pctMultiplier)

	expectedRespPosition := resp.StakePosition{
		Symbol:                   symbol,
		Name:                     name,
		OpenQty:                  stkPosition.OpenQty.StringFixed(satoshiDecimalPlaces),
		AvailableForTradingQty:   stkPosition.AvailableForTradingQty.StringFixed(satoshiDecimalPlaces),
		AveragePrice:             stkPosition.AveragePrice.StringFixed(satoshiDecimalPlaces),
		MarketValue:              mktValue.StringFixed(satoshiDecimalPlaces),
		MarketPrice:              mktPrice.StringFixed(satoshiDecimalPlaces),
		PriorClose:               priorClose.StringFixed(satoshiDecimalPlaces),
		DayProfitOrLoss:          dayPNL.StringFixed(satoshiDecimalPlaces),
		DayProfitOrLossPercent:   dayPNLPCT.StringFixed(pctDecimalPlaces),
		TotalProfitOrLoss:        totalPNL.StringFixed(satoshiDecimalPlaces),
		TotalProfitOrLossPercent: totalPNLPCT.StringFixed(pctDecimalPlaces),
	}

	s.Equal(len(s.d.positions), 1)

	s.Equal(expectedRespPosition, *s.d.positions[stkPosition.Symbol])
}

func (s *dataManagerTestSuite) OnStkPosition(symbol, name string) *stk.InternalPosition {
	stkPosition := &stk.InternalPosition{
		Symbol:                 symbol,
		Name:                   name,
		OpenQty:                decimal.NewFromFloat(10.0000),
		AvailableForTradingQty: decimal.NewFromFloat(10.0000),
		AveragePrice:           decimal.NewFromFloat(102.5000),
		Cost:                   decimal.NewFromFloat(1025.0000),
	}

	s.d.onMarketPositions([]*stk.InternalPosition{stkPosition})

	return stkPosition
}

func (s *dataManagerTestSuite) OnStkPrice(symbol string) *stk.InternalPrice {
	stkPrice := &stk.InternalPrice{
		Symbol:     symbol,
		LastTrade:  decimal.NewFromFloat(114.9800),
		Bid:        decimal.NewFromFloat(114.98),
		Ask:        decimal.NewFromFloat(114.99),
		PriorClose: decimal.NewFromFloat(119.8700),
	}
	s.d.onMarketPrices([]*stk.InternalPrice{stkPrice})

	return stkPrice
}

func TestDataManagerTestSuite(t *testing.T) {
	suite.Run(t, new(dataManagerTestSuite))
}
