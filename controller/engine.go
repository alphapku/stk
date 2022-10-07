package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"StakeBackendGoTest/api"
	cfg "StakeBackendGoTest/configs"
	log "StakeBackendGoTest/pkg/log"
)

const (
	timeout = 5 * time.Second
)

type Engine struct {
	*gin.Engine
	*http.Server

	cfg *cfg.Config
}

func NewEngine(cfg *cfg.Config) *Engine {
	e := &Engine{
		Engine: gin.New(),
		cfg:    cfg,
	}

	return e
}

func (e *Engine) Run() error {
	if err := e.validate(); err != nil {
		return err
	}

	api.AddRouters(e.cfg, e.Engine)

	e.Server = &http.Server{
		Addr:    e.cfg.Addr,
		Handler: e.Engine,
	}

	go func() {
		if err := e.ListenAndServe(); err != nil {
			log.Logger.Error("failed to start server: %v", zap.Error(err))
		}
	}()

	return nil
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
