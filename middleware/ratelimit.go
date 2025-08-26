package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter stores rate limit information
type RateLimiter struct {
	mu     sync.Mutex
	visits map[string][]time.Time
}

// rateLimitConfig holds configuration for rate limiting
type rateLimitConfig struct {
	limit  int
	window time.Duration
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		visits: make(map[string][]time.Time),
	}
}

// LimitBy creates a middleware with dynamic limits
func (rl *RateLimiter) LimitBy(limit int, window time.Duration) gin.HandlerFunc {
	config := rateLimitConfig{
		limit:  limit,
		window: window,
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		rl.mu.Lock()
		defer rl.mu.Unlock()

		// Clean up old visits for this IP
		now := time.Now()
		validVisits := []time.Time{}
		for _, visitTime := range rl.visits[clientIP] {
			if now.Sub(visitTime) <= config.window {
				validVisits = append(validVisits, visitTime)
			}
		}

		// Check if rate limit exceeded
		if len(validVisits) >= config.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":      "error",
				"message":     "Rate limit exceeded. Please try again later.",
				"retry_after": config.window.Seconds(),
				"limit":       config.limit,
				"window":      config.window.String(),
			})
			c.Abort()
			return
		}

		// Add current visit
		validVisits = append(validVisits, now)
		rl.visits[clientIP] = validVisits

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", string(rune(config.limit)))
		c.Header("X-RateLimit-Remaining", string(rune(config.limit - len(validVisits))))
		c.Header("X-RateLimit-Reset", now.Add(config.window).Format(time.RFC1123))

		c.Next()
	}
}

// CleanupOldVisits periodically cleans up old entries
func (rl *RateLimiter) CleanupOldVisits() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, visits := range rl.visits {
			validVisits := []time.Time{}
			for _, visitTime := range visits {
				if now.Sub(visitTime) <= time.Hour { // Clean up anything older than 1 hour
					validVisits = append(validVisits, visitTime)
				}
			}
			if len(validVisits) == 0 {
				delete(rl.visits, ip)
			} else {
				rl.visits[ip] = validVisits
			}
		}
		rl.mu.Unlock()
	}
}