package mock

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	mk "StakeBackendGoTest/internal/entity/mock"
	intl "StakeBackendGoTest/internal/entity/stake"
)

type converterTestSuite struct {
	suite.Suite
}

func (s *converterTestSuite) SetupSuite() {
}

func (s *converterTestSuite) TestPositionConverter() {
	tests := []struct {
		pos         *mk.Position
		expectedPos *intl.InternalPosition
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
		expectedPos: &intl.InternalPosition{
			Symbol:                 "",
			Name:                   "",
			OpenQTY:                decimal.Decimal{},
			AvailableForTradingQTY: decimal.Decimal{},
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
			expectedPos: &intl.InternalPosition{
				Symbol:                 "APT.ASX",
				Name:                   "Afterpay Limited",
				AvailableForTradingQTY: decimal.NewFromFloat(10.0000),
				AveragePrice:           decimal.NewFromFloat(102.5000),
				Cost:                   decimal.NewFromFloat(1025.0000),
				OpenQTY:                decimal.NewFromFloat(10.0000),
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		actual, err := ToStakePosition(test.pos)
		if test.expectedErr != nil {
			s.ErrorIs(err, test.expectedErr)
		} else {
			s.Equal(test.expectedPos.Equal(*actual), true)
		}
	}
}

func (s *converterTestSuite) TestPriceConverter() {
	// TODO,
}

func TestConverters(t *testing.T) {
	suite.Run(t, new(converterTestSuite))
}
