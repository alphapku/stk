package adapter

import (
	"context"
	"errors"

	cfg "StakeBackendGoTest/configs"
	def "StakeBackendGoTest/pkg/const"
	log "StakeBackendGoTest/pkg/log"
)

var (
	ErrAdapterUnknown = errors.New("unknown adapter")
)

type AdapterManager struct {
	adapter Adapter
}

func NewAdapterManager(cfg *cfg.Adapter) (*AdapterManager, error) {
	a := &AdapterManager{}

	switch cfg.AdapterType {
	case def.MockAdapter:
		a.adapter = NewMockAdapter(cfg.MockMSGCount, cfg.MSGIntervalSecs)
	default:
		return nil, ErrAdapterUnknown
	}

	return a, nil
}

// Start returns a channel to indicate that it stops
func (a *AdapterManager) Start(ctx context.Context, dataChan chan interface{}) (<-chan struct{}, error) {
	done, err := a.adapter.Start(ctx, dataChan)
	if err != nil {
		return nil, err
	}

	c := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			log.Logger.Debug("cancelling adapters")
			a.adapter.Close(ctx)
			c <- struct{}{}
			return
		case <-done:
			log.Logger.Debug("adapters closed")
			c <- struct{}{}
			return
		}

	}()

	return c, nil
}
