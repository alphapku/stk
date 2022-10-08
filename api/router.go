package api

import (
	"github.com/gin-gonic/gin"

	mdw "StakeBackendGoTest/api/middleware"
	mdl "StakeBackendGoTest/internal/model"
)

const (
	equityPositionURI = "/api/equityPositions"
)

func AddRouters(e *gin.Engine, d *mdl.DataManager) {
	e.Use(mdw.Auth, mdw.Trace)

	e.POST(equityPositionURI, d.DoEquityPositions)
}
