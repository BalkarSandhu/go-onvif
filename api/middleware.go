package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *APIServer) rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !s.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (s *APIServer) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		s.logger.Info().
			Str("client_ip", c.ClientIP()).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("latency", time.Since(start)).
			Msg("Request processed")
	}
}
