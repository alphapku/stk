package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	cfg "StakeBackendGoTest/configs"
	"StakeBackendGoTest/controller"
	def "StakeBackendGoTest/pkg/const"
	log "StakeBackendGoTest/pkg/log"
)

func main() {
	// r := gin.Default()

	/* TODO: client and server error handling
	 *  - handle an invalid token being passed in
	 *  - handle case where invalid/null equityPositions are returned by the service class
	 */

	// r.GET("/api/equityPositions", internal.GetEquityPositionsHandler)
	// r.Run("localhost:8080")

	// TODO, load Config from file
	cfg := &cfg.Config{
		Adapter: &cfg.Adapter{
			AdapterType: def.MockAdapter,
		},
		Addr:    "localhost:8080",
		EnvMode: def.DevMode,
	}

	// initialize the logger
	if err := log.Init(cfg.EnvMode); err != nil {
		panic("failed to initialize the logger") // allowed to panic in `main``
	}

	e := controller.NewEngine(cfg)
	if err := e.Run(); err != nil {
		log.Logger.Panic("failed to start the server", zap.Error(err))
	}

	// install signal processor
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	if err := e.Stop(); err != nil {
		log.Logger.Panic("failed to stop the server", zap.Error(err))
	} else {
		log.Logger.Info("shutting down gracefully")
	}
}
