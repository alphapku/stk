package api

import (
	"github.com/gin-gonic/gin"

	mdw "StakeBackendGoTest/api/middleware"
	mdl "StakeBackendGoTest/internal/model"
)

func AddRouters(e *gin.Engine, d *mdl.DataManager) {
	e.Use(mdw.Auth, mdw.Trace)

	e.GET("/api/equityPositions", d.DoEquityPositions)
}
