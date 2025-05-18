package metrics

import (
    "github.com/gin-gonic/gin"
    "strconv"
    "time"
)

func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())

        path := c.FullPath()
        if path == "" {
            path = c.Request.URL.Path
        }

        RequestCount.WithLabelValues(c.Request.Method, path, status).Inc()
        RequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)

        if c.Writer.Status() >= 500 {
            RequestErrors.WithLabelValues(c.Request.Method, path, status).Inc()
        }
    }
}
