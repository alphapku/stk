package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	log "StakeBackendGoTest/pkg/log"
)

const (
	tokenKey  = "token"
	userIDKey = "user_id"

	testUserID = "f9e27aa9-727c-4fd8-907e-46522421089d"
	testToken  = "fJCoxhq8uR9GiUIgaIGfMgw7zCqxwDhQ"
)

func Auth(c *gin.Context) {
	userID := c.GetHeader("userIDKey")
	token := c.GetHeader(tokenKey)
	if token != testToken || userID != testUserID {
		ip, _ := c.RemoteIP()
		log.Logger.Debug("invalid token/user_id", zap.String("ip", ip.String()))
		c.JSON(http.StatusOK, gin.H{"error_msg": "invalid token/user_id"})
		c.Abort()
		return
	}

	c.Next()
}
