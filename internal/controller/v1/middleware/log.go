package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Log(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("request", slog.String("method", c.Request.Method), slog.String("uri", c.Request.URL.Path))
		c.Next()
	}
}
