package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"StakeBackendGoTest/pkg/log"
)

func Trace(c *gin.Context) {
	s := time.Now()
	defer func() {
		e := time.Since(s)
		ip, _ := c.RemoteIP()
		log.Logger.Debug("tracing", zap.String("client", ip.String()), zap.String("host", c.Request.Host), zap.String("url", c.Request.URL.Path), zap.String("latency", e.String()))
	}()

	c.Next()
}
