package api

import (
	"github.com/gin-gonic/gin"

	mdw "StakeBackendGoTest/api/middleware"
	// cfg "StakeBackendGoTest/configs"
	ctrl "StakeBackendGoTest/controller"
	// itnl "StakeBackendGoTest/internal"
)

func AddRouters(e *gin.Engine, d *ctrl.DataManager) {
	e.Use(mdw.Auth, mdw.Trace)

	e.GET("/api/equityPositions", d.DoEquityPositions)
}
