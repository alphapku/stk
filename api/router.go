package api

import (
	"github.com/gin-gonic/gin"

	mdw "StakeBackendGoTest/api/middleware"
	cfg "StakeBackendGoTest/configs"
	itnl "StakeBackendGoTest/internal"
)

func AddRouters(cfg *cfg.Config, e *gin.Engine) {
	e.Use(mdw.Trace)

	e.GET("/api/equityPositions", itnl.NewAdapter(cfg.Adapter))
}
