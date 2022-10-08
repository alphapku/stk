package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	resp "StakeBackendGoTest/api/response"
	log "StakeBackendGoTest/pkg/log"
)

const (
	tokenKey = "token"
	// userIDKey = "user_id"

	// testUserID = "f9e27aa9-727c-4fd8-907e-46522421089d"
	testToken = "fJCoxhq8uR9GiUIgaIGfMgw7zCqxwDhQ"
)

// Auth aborts the request if failing to authenticate the request.
// TODO, combing database and jwt-go to generate a dynamic token with specific user in Login endpoint
func Auth(c *gin.Context) {
	ip, _ := c.RemoteIP()
	// userID, ok := c.GetQuery("userIDKey")
	// if !ok || userID != testUserID {
	// 	log.Logger.Debug("invalid user_id", zap.String("ip", ip.String()))
	// 	c.JSON(http.StatusOK, gin.H{"error_msg": "invalid user_id"})
	// 	c.Abort()
	// 	return
	// }
	token, ok := c.GetQuery(tokenKey)
	if !ok || token != testToken {
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
