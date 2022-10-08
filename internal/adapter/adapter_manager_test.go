package adapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	cfg "StakeBackendGoTest/configs"
	def "StakeBackendGoTest/pkg/const"
	log "StakeBackendGoTest/pkg/log"
)

type adapterManagerTestSuite struct {
	suite.Suite
}

func (s *adapterManagerTestSuite) SetupSuite() {
	_ = log.Init(def.DevMode)
}

func (s *adapterManagerTestSuite) TestAdapterManagerCancelled() {
	cfg := &cfg.Adapter{
		AdapterType:     def.CoinbaseAdapter,
		MockMSGCount:    2,
		MSGIntervalSecs: 3, // 6 seconds is enough for us to test all here
	}

	m, err := NewAdapterManager(cfg)
	s.ErrorIs(err, ErrAdapterUnknown)
	s.Nil(m)

	cfg.AdapterType = def.MockAdapter
	m, err = NewAdapterManager(cfg)
	s.Nil(err)
	s.NotNil(m)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	dataChan := make(chan interface{}, 10)
	done, _ := m.Start(ctx, dataChan)

	msg := <-dataChan
	s.NotNil(msg) // do received messages

	cancel() // this would shutdown the manager
	<-done
}

func (s *adapterManagerTestSuite) TestAdapterManagerQuitting() {
	cfg := &cfg.Adapter{
		AdapterType:     def.MockAdapter,
		MockMSGCount:    0,
		MSGIntervalSecs: 0,
	}

	m, _ := NewAdapterManager(cfg)

	ctx := context.Background()
	ctx, _ = context.WithCancel(ctx)
	dataChan := make(chan interface{}, 10)
	done, _ := m.Start(ctx, dataChan)

	// we don't call cancel, the mocker quits with the zero MockMSGCount
	<-done
}

func TestAdapterManagerTestSuite(t *testing.T) {
	suite.Run(t, new(adapterManagerTestSuite))
}
