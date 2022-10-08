package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	resp "StakeBackendGoTest/api/response"
	log "StakeBackendGoTest/pkg/log"
)

const (
	TokenKey  = "token"
	TestToken = "fJCoxhq8uR9GiUIgaIGfMgw7zCqxwDhQ"
)

// Auth aborts the request if failing to authenticate the request.
// TODO, combing database and jwt-go to generate a dynamic token with specific user in Login endpoint
func Auth(c *gin.Context) {
	ip, _ := c.RemoteIP()
	token := c.GetHeader(TokenKey)
	if token != TestToken {
		log.Logger.Debug("invalid token", zap.String("ip", ip.String()))
		c.JSON(http.StatusOK, resp.Response{
			ErrorCode: resp.InvalidToken,
			Data:      "invalid token",
		})
		c.Abort()
		return
	}

	c.Next()
}
