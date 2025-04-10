package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var newlimiter = rate.NewLimiter(1, 5)

func RateLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !newlimiter.Allow() {
			ctx.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				gin.H{
					"message": "Not Allowed",
				},
			)
			return
		}
		ctx.Next()
	}
}
