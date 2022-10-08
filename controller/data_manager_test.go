package controller

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	stk "StakeBackendGoTest/internal/entity/stake"
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
	tests := []*stk.Position{{
		Symbol:                 "",
		Name:                   "",
		OpenQty:                decimal.Decimal{},
		AvailableForTradingQty: decimal.Decimal{},
		AveragePrice:           decimal.Decimal{},
		Cost:                   decimal.Decimal{},
	},
	}

	for _, test := range tests {
		s.d.onMarketPositions([]*stk.Position{test})
	}
	s.Equal(1, 2)
}

func TestDataManagerTestSuite(t *testing.T) {
	suite.Run(t, new(dataManagerTestSuite))
}
