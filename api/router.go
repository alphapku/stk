package api

import (
	"github.com/gin-gonic/gin"

	mdw "StakeBackendGoTest/api/middleware"
	ctrl "StakeBackendGoTest/controller"
)

func AddRouters(e *gin.Engine, d *ctrl.DataManager) {
	e.Use(mdw.Auth, mdw.Trace)

	e.GET("/api/equityPositions", d.DoEquityPositions)
}
