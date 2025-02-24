package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	visits map[string]int
	mu     sync.Mutex
}

var (
	limit    = 100
	duration = time.Minute
	rl       = &rateLimiter{
		visits: make(map[string]int),
	}
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		defer rl.mu.Unlock()

		if rl.visits[ip] >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		rl.visits[ip]++
		time.AfterFunc(duration, func() {
			rl.mu.Lock()
			defer rl.mu.Unlock()
			rl.visits[ip]--
		})
		c.Next()
	}
}
