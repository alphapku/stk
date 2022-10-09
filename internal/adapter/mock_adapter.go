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

type MockAdapter struct {
	mockMSGCount int
	msgInterval  time.Duration

	intlPositions []*intl.InternalPosition
	intlPrices    []*intl.InternalPrice
}

func NewMockAdapter(mockMSGCount int, msgIntervalSecs int) *MockAdapter {
	return &MockAdapter{
		mockMSGCount:  mockMSGCount,
		msgInterval:   time.Duration(msgIntervalSecs) * time.Second,
		intlPositions: make([]*intl.InternalPosition, 0),
		intlPrices:    make([]*intl.InternalPrice, 0),
	}
}

func (m *MockAdapter) Close(ctx context.Context) {
}

func (m *MockAdapter) Start(ctx context.Context, dataChan chan interface{}) (<-chan struct{}, error) {
	var err error
	m.intlPositions, m.intlPrices, err = LoadAndParseMockData("../../")
	if err != nil {
		return nil, err
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	done := make(chan struct{})
	go func() {
		sent := 0
		ticker := time.NewTicker(m.msgInterval)
		for {
			if sent >= m.mockMSGCount {
				log.Logger.Debug("mocking data done")
				close(done)
				return
			}
			sent = sent + 1
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if r1.Intn(100) > 33 {
					dataChan <- m.intlPositions
				} else {
					dataChan <- m.intlPrices
				}
			}
		}
	}()

	return done, nil
}

func LoadAndParseMockData(dir string) ([]*intl.InternalPosition, []*intl.InternalPrice, error) {
	positions, err := ReadMockPositionData(dir)
	if err != nil {
		return nil, nil, err
	}

	intlPositions := make([]*intl.InternalPosition, 0)
	for _, pos := range positions.Positions {
		if o, err := cvt.ToStakePosition(pos); err == nil {
			intlPositions = append(intlPositions, o)
		} else {
			log.Logger.Warn("failed to convert", zap.String("symbol", pos.Security), zap.Error(err))
		}
	}

	prices, err := ReadMockPriceData(dir)
	if err != nil {
		return nil, nil, err
	}

	intlPrices := make([]*intl.InternalPrice, 0)
	for _, prx := range prices.Prices {
		intlPrices = append(intlPrices, cvt.ToStakePrice(prx))
	}

	return intlPositions, intlPrices, nil
}

func ReadMockPositionData(dir string) (*mk.Positions, error) {
	file, err := ioutil.ReadFile(fmt.Sprintf("%sinternal/adapter/mockdata/mockpositions.json", dir))
	if err != nil {
		return nil, err
	}

	data := &mk.Positions{}
	err = json.Unmarshal([]byte(file), data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ReadMockPriceData(dir string) (*mk.Prices, error) {
	file, err := ioutil.ReadFile(fmt.Sprintf("%sinternal/adapter/mockdata/mockprices.json", dir))
	if err != nil {
		return nil, err
	}

	data := &mk.Prices{}
	err = json.Unmarshal([]byte(file), data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
