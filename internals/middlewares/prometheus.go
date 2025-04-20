package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// RequestCount
// ErrorCount

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests count",
			Help: "counts the number of requests from the backend",
		},
		[]string{"path", "status"},
	)

	ErrRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "err requests count",
			Help: "counts the number of error requests from the backend",
		},
		[]string{"path", "status"},
	)
)

func Init() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(ErrRequestCount)
}

func TrackMetrics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		status := ctx.Writer.Status()
		ctx.Next()

		if status > 400 {
			ErrRequestCount.WithLabelValues(path, http.StatusText(status)).Inc()
		}
		RequestCount.WithLabelValues(path, http.StatusText(status))
	}
}
