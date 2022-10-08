package adapter

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"time"

	mk "StakeBackendGoTest/internal/entity/mock"
	stk "StakeBackendGoTest/internal/entity/stake"
	cvt "StakeBackendGoTest/internal/pkg/converters/mock"
	log "StakeBackendGoTest/pkg/log"

	"go.uber.org/zap"
)

const (
	maxMessageSent = 20              // change this for longer mocking
	dataInterval   = 3 * time.Second // change this for faster mocking
)

type MockAdapter struct {
	mockMsgCount int
	msgInterval  time.Duration

	stkPositions []*stk.InternalPosition
	stkPrices    []*stk.InternalPrice
}

func NewMockAdapter(mockMsgCount int, msgInterval time.Duration) *MockAdapter {
	return &MockAdapter{
		mockMsgCount: mockMsgCount,
		msgInterval:  msgInterval,
		stkPositions: make([]*stk.InternalPosition, 0),
		stkPrices:    make([]*stk.InternalPrice, 0),
	}
}

func NewDefaultMockAdapter() *MockAdapter {
	return NewMockAdapter(maxMessageSent, dataInterval)
}

func (m *MockAdapter) Close(ctx context.Context) {
}

func (m *MockAdapter) loadAndParseMockData() {
	positions := readMockPositionData()
	for _, pos := range positions.Positions {
		if o, err := cvt.ToStakePosition(pos); err == nil {
			m.stkPositions = append(m.stkPositions, o)
		} else {
			log.Logger.Debug("failed to convert", zap.String("symbol", pos.Security), zap.Error(err))
		}
	}

	prices := readMockPriceData()
	for _, prx := range prices.Prices {
		m.stkPrices = append(m.stkPrices, cvt.ToStakePrice(prx))
	}
}

func (m *MockAdapter) Start(ctx context.Context, dataChan chan interface{}) (<-chan struct{}, error) {
	m.loadAndParseMockData()

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	done := make(chan struct{})
	go func() {
		sent := 0
		ticker := time.NewTicker(dataInterval)
		for {
			if sent >= maxMessageSent {
				log.Logger.Info("mocking data done")
				close(done)
				return
			}
			sent = sent + 1
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if r1.Intn(100) > 66 {
					dataChan <- m.stkPositions
				} else {
					dataChan <- m.stkPrices
				}
			}
		}
	}()

	return done, nil
}

func readMockPositionData() mk.Positions {
	file, _ := ioutil.ReadFile("../../internal/adapter/mockdata/mockpositions.json")
	data := mk.Positions{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func readMockPriceData() mk.Prices {
	file, _ := ioutil.ReadFile("../../internal/adapter/mockdata/mockprices.json")
	data := mk.Prices{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}
