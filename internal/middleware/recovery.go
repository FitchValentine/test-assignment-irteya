package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		fields := []zap.Field{
			zap.Any("error", recovered),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		}

		if traceID, exists := c.Get("trace_id"); exists {
			if traceIDStr, ok := traceID.(string); ok {
				fields = append(fields, zap.String("trace_id", traceIDStr))
			}
		}

		logger.Error("panic recovered", fields...)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	})
}

