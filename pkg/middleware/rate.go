// pkg/middleware/rate.go
package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"github.com/gin-gonic/gin"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Minute/10), 10)
		visitors[ip] = limiter
	}
	return limiter
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)

		if !getLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}