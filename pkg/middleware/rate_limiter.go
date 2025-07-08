package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter holds the rate limiter for each IP
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit // requests per second
	b   int        // burst size
}

// NewIPRateLimiter creates a new IP rate limiter
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// AddIP creates a new rate limiter and adds it to the ips map
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter

	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address if it exists.
// Otherwise calls AddIP to add IP address to the map
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddIP(ip)
	}

	return limiter
}

// RateLimitMiddleware is a middleware that limits the rate of incoming requests
func RateLimitMiddleware() gin.HandlerFunc {
	// Allow 100 requests per minute with a burst of 20
	limiter := NewIPRateLimiter(rate.Limit(100.0/60.0), 20)

	return func(c *gin.Context) {
		// Get the IP address for the current user.
		ip := c.ClientIP()

		// Get the limiter for the current IP.
		ipLimiter := limiter.GetLimiter(ip)

		// Check if the request is allowed.
		if !ipLimiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"message": "The service has received too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
