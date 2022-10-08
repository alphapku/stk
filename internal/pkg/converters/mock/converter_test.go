package mock

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	mk "StakeBackendGoTest/internal/entity/mock"
	stk "StakeBackendGoTest/internal/entity/stake"
)

type converterTestSuite struct {
	suite.Suite
}

func (s *converterTestSuite) SetupSuite() {
}

func (s *converterTestSuite) TestPositionConverter() {
	tests := []struct {
		pos         *mk.Position
		expectedPos *stk.Position
		expectedErr error
	}{{
		// error case
		pos: &mk.Position{
			Security:            "",
			SecurityDescription: "",
			Cost:                decimal.Decimal{},
			AveragePrice:        decimal.Decimal{},
			AvailableUnits:      0,
			PortfolioUnits:      0,
		},
		expectedPos: &stk.Position{
			Symbol:                 "",
			Name:                   "",
			OpenQty:                decimal.Decimal{},
			AvailableForTradingQty: decimal.Decimal{},
			AveragePrice:           decimal.Decimal{},
			Cost:                   decimal.Decimal{},
		},
		expectedErr: ErrZeroAveragePrice,
	},
		// normal case
		{
			pos: &mk.Position{
				Security:            "APT.ASX",
				SecurityDescription: "Afterpay Limited",
				Cost:                decimal.NewFromFloat(1025.0000),
				AveragePrice:        decimal.NewFromFloat(102.5000),
				AvailableUnits:      10,
				PortfolioUnits:      10,
			},
			expectedPos: &stk.Position{
				Symbol:                 "APT.ASX",
				Name:                   "Afterpay Limited",
				OpenQty:                decimal.NewFromFloat(10.0000),
				AvailableForTradingQty: decimal.NewFromFloat(10.0000),
				AveragePrice:           decimal.NewFromFloat(102.5000),
				Cost:                   decimal.NewFromFloat(1025.0000),
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		actual, err := ToStakePosition(test.pos)
		if test.expectedErr != nil {
			s.ErrorIs(err, test.expectedErr)
		} else {
			s.Equal(test.expectedPos.Symbol, actual.Symbol)
			s.Equal(test.expectedPos.Name, actual.Name)
			// kind of silly, as decimal equality comparison failed for something like with `s.Equal(d1, d2)`
			// 	d1 := decimal.New(100, 1)
			//  d2 := decimal.New(1000, 0)
			s.Assert().Equal(test.expectedPos.OpenQty.Equal(actual.OpenQty), true)
			s.Assert().Equal(test.expectedPos.AvailableForTradingQty.Equal(actual.AvailableForTradingQty), true)
			s.Assert().Equal(test.expectedPos.AveragePrice.Equal(actual.AveragePrice), true)
			s.Assert().Equal(test.expectedPos.Cost.Equal(actual.Cost), true)
		}
	}
}

func (s *converterTestSuite) TestPriceConverter() {
	// TODO,
}

func TestConverters(t *testing.T) {
	suite.Run(t, new(converterTestSuite))
}
