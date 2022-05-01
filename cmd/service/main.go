package main

import (
	"StakeBackendGoTest/internal"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	/* TODO: client and server error handling
	 *  - handle an invalid token being passed in
	 *  - handle case where invalid/null equityPositions are returned by the service class
	 */

	r.GET("/api/equityPositions", internal.GetEquityPositionsHandler)
	r.Run("localhost:8080")
}
