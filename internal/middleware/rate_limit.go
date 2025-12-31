package middleware

import (
	"demo-service/internal/config"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type rateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      rate.Limit
}

var globalLimiter *rateLimiter

func initRateLimiter() {
	rps := rate.Limit(config.AppConfig.RateLimitRPS)
	globalLimiter = &rateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rps,
	}
}

func (rl *rateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter, exists = rl.limiters[key]
		if !exists {
			limiter = rate.NewLimiter(rl.rps, int(rl.rps))
			rl.limiters[key] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter
}

func RateLimitMiddleware() gin.HandlerFunc {
	if globalLimiter == nil {
		initRateLimiter()
	}

	return func(c *gin.Context) {
		key := c.ClientIP()
		limiter := globalLimiter.getLimiter(key)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
