package middleware

import (
	"demo-service/internal/metrics"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		metrics.ActiveConnections.Inc()
		defer metrics.ActiveConnections.Dec()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		endpoint := c.FullPath()

		metrics.HTTPRequestDuration.WithLabelValues(method, endpoint, status).Observe(duration)
		metrics.HTTPRequestTotal.WithLabelValues(method, endpoint, status).Inc()
	}
}



