package adapter

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	mk "StakeBackendGoTest/internal/entity/mock"
	stk "StakeBackendGoTest/internal/entity/stake"
	cvt "StakeBackendGoTest/internal/pkg/converters/mock"
	log "StakeBackendGoTest/pkg/log"
)

const (
	maxMessageSent = 10
	dataInterval   = 3 * time.Second
)

type MockAdapter struct {
	file *os.File
}

func (f *MockAdapter) Close(ctx context.Context) {
	f.file.Close()
}

func (f *MockAdapter) Start(ctx context.Context, dataChan chan interface{}) (<-chan struct{}, error) {
	positions := readMockPositionData()
	stkPositions := make([]*stk.Position, len(positions.Positions))
	for i, pos := range positions.Positions {
		stkPositions[i] = cvt.ToStakePosition(pos)
	}

	prices := readMockPriceData()
	stkPrices := make([]*stk.Price, len(prices.Prices))
	for i, prx := range prices.Prices {
		stkPrices[i] = cvt.ToStakePrice(prx)
	}

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
					dataChan <- stkPositions
				} else {
					dataChan <- stkPrices
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
