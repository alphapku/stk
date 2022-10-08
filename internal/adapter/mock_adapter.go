package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	mk "StakeBackendGoTest/internal/entity/mock"
	intl "StakeBackendGoTest/internal/entity/stake"
	cvt "StakeBackendGoTest/internal/pkg/converters/mock"
	log "StakeBackendGoTest/pkg/log"

	"go.uber.org/zap"
)

const (
	maxMessageSent = 1000            // change this for longer mocking
	dataInterval   = 3 * time.Second // change this for faster mocking
)

type MockAdapter struct {
	mockMsgCount int
	msgInterval  time.Duration

	intlPositions []*intl.InternalPosition
	intlPrices    []*intl.InternalPrice
}

func NewMockAdapter(mockMsgCount int, msgInterval time.Duration) *MockAdapter {
	return &MockAdapter{
		mockMsgCount:  mockMsgCount,
		msgInterval:   msgInterval,
		intlPositions: make([]*intl.InternalPosition, 0),
		intlPrices:    make([]*intl.InternalPrice, 0),
	}
}

func NewDefaultMockAdapter() *MockAdapter {
	return NewMockAdapter(maxMessageSent, dataInterval)
}

func (m *MockAdapter) Close(ctx context.Context) {
}

func (m *MockAdapter) Start(ctx context.Context, dataChan chan interface{}) (<-chan struct{}, error) {
	m.intlPositions, m.intlPrices = LoadAndParseMockData("../../")

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
					dataChan <- m.intlPositions
				} else {
					dataChan <- m.intlPrices
				}
			}
		}
	}()

	return done, nil
}

func LoadAndParseMockData(dir string) ([]*intl.InternalPosition, []*intl.InternalPrice) {
	intlPositions := make([]*intl.InternalPosition, 0)
	intlPrices := make([]*intl.InternalPrice, 0)
	positions := ReadMockPositionData(dir)
	for _, pos := range positions.Positions {
		if o, err := cvt.ToStakePosition(pos); err == nil {
			intlPositions = append(intlPositions, o)
		} else {
			log.Logger.Debug("failed to convert", zap.String("symbol", pos.Security), zap.Error(err))
		}
	}

	prices := ReadMockPriceData(dir)
	for _, prx := range prices.Prices {
		intlPrices = append(intlPrices, cvt.ToStakePrice(prx))
	}

	return intlPositions, intlPrices
}

func ReadMockPositionData(dir string) mk.Positions {
	file, _ := ioutil.ReadFile(fmt.Sprintf("%sinternal/adapter/mockdata/mockpositions.json", dir))
	data := mk.Positions{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func ReadMockPriceData(dir string) mk.Prices {
	file, _ := ioutil.ReadFile(fmt.Sprintf("%sinternal/adapter/mockdata/mockprices.json", dir))
	data := mk.Prices{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}
