package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	// "StakeBackendGoTest/api"
	cfg "StakeBackendGoTest/configs"
	adt "StakeBackendGoTest/internal/adapter"
	mdl "StakeBackendGoTest/internal/model"
	log "StakeBackendGoTest/pkg/log"
)

const (
	timeout     = 5 * time.Second
	dataChanLen = 8192
)

type Engine struct {
	*gin.Engine
	*http.Server

	cfg *cfg.Config

	DataManager    *mdl.DataManager
	adapterManager *adt.AdapterManager
}

func NewEngine(cfg *cfg.Config) *Engine {
	e := &Engine{
		Engine:      gin.New(),
		DataManager: mdl.NewDataManager(),
		cfg:         cfg,
	}

	return e
}

func (e *Engine) Run(ctx context.Context) (<-chan struct{}, error) {
	if err := e.validate(); err != nil {
		return nil, err
	}

	dataChan := make(chan interface{}, dataChanLen)

	// api.AddRouters(e.cfg, e.Engine)
	am, err := adt.NewAdapterManager(e.cfg.Adapter)
	if err != nil {
		return nil, err
	}

	e.adapterManager = am
	e.Server = &http.Server{
		Addr:    e.cfg.Addr,
		Handler: e.Engine,
	}

	done := make(chan struct{}, 2)

	go func() {
		// for ... range is the good option as we only have one data channel at present
		for msg := range dataChan {
			e.DataManager.OnMessage(msg)
		}

	}()

	go func() {
		adapterManagerDone, err := e.adapterManager.Start(ctx, dataChan)
		if err == nil {
			<-adapterManagerDone
		}

		done <- struct{}{}
	}()

	go func() {
		if err := e.ListenAndServe(); err != nil {
			log.Logger.Error("failed to start server or it was closed", zap.Error(err))
			done <- struct{}{}
		}
	}()

	return done, nil
}

func (e *Engine) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return e.Shutdown(ctx)
}

func (e *Engine) validate() error {
	// TODO, return an error if failing to validate e.cfg
	return nil
}
