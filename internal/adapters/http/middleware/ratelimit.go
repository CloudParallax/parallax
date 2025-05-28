package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

// RateLimiter represents a rate limiter for a specific client
type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate int
	lastRefill time.Time
	mutex      sync.Mutex
}

// RateLimitMiddleware provides rate limiting functionality
type RateLimitMiddleware struct {
	limiters    map[string]*RateLimiter
	maxRequests int
	window      time.Duration
	keyFunc     func(fiber.Ctx) string
	mutex       sync.RWMutex
}

// RateLimitConfig holds rate limit configuration
type RateLimitConfig struct {
	MaxRequests int                    // Maximum requests allowed
	Window      time.Duration          // Time window for rate limiting
	KeyFunc     func(fiber.Ctx) string // Function to extract client identifier
	Message     string                 // Custom message for rate limit exceeded
}

// DefaultRateLimitConfig returns default rate limit configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		MaxRequests: 100,
		Window:      time.Minute,
		KeyFunc:     defaultKeyFunc,
		Message:     "Rate limit exceeded",
	}
}

// NewRateLimitMiddleware creates a new rate limiting middleware
func NewRateLimitMiddleware(config ...RateLimitConfig) *RateLimitMiddleware {
	cfg := DefaultRateLimitConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return &RateLimitMiddleware{
		limiters:    make(map[string]*RateLimiter),
		maxRequests: cfg.MaxRequests,
		window:      cfg.Window,
		keyFunc:     cfg.KeyFunc,
		mutex:       sync.RWMutex{},
	}
}

// Handler returns the rate limiting middleware handler
func (r *RateLimitMiddleware) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		key := r.keyFunc(c)
		
		// Get or create rate limiter for this key
		limiter := r.getRateLimiter(key)
		
		// Check if request is allowed
		if !limiter.Allow() {
			// Set rate limit headers
			c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", r.maxRequests))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(r.window).Unix()))
			
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    fiber.StatusTooManyRequests,
					"message": "Rate limit exceeded",
				},
			})
		}
		
		// Set rate limit headers for successful requests
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", r.maxRequests))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.tokens))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(r.window).Unix()))
		
		return c.Next()
	}
}

// getRateLimiter gets or creates a rate limiter for the given key
func (r *RateLimitMiddleware) getRateLimiter(key string) *RateLimiter {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	limiter, exists := r.limiters[key]
	if !exists {
		limiter = &RateLimiter{
			tokens:     r.maxRequests,
			maxTokens:  r.maxRequests,
			refillRate: r.maxRequests,
			lastRefill: time.Now(),
		}
		r.limiters[key] = limiter
	}
	
	return limiter
}

// Allow checks if a request is allowed under the rate limit
func (rl *RateLimiter) Allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	
	// Calculate tokens to add based on time elapsed
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed.Minutes()) * rl.refillRate
	
	if tokensToAdd > 0 {
		rl.tokens = min(rl.maxTokens, rl.tokens+tokensToAdd)
		rl.lastRefill = now
	}
	
	// Check if we have tokens available
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	
	return false
}

// defaultKeyFunc extracts client IP as the rate limiting key
func defaultKeyFunc(c fiber.Ctx) string {
	return c.IP()
}

// KeyFuncByIP creates a key function that uses client IP
func KeyFuncByIP() func(fiber.Ctx) string {
	return func(c fiber.Ctx) string {
		return c.IP()
	}
}

// KeyFuncByUserID creates a key function that uses authenticated user ID
func KeyFuncByUserID() func(fiber.Ctx) string {
	return func(c fiber.Ctx) string {
		userID := c.Locals("user_id")
		if userID != nil {
			return fmt.Sprintf("user_%s", userID)
		}
		return c.IP() // Fallback to IP if no user ID
	}
}

// KeyFuncByHeader creates a key function that uses a specific header
func KeyFuncByHeader(headerName string) func(fiber.Ctx) string {
	return func(c fiber.Ctx) string {
		header := c.Get(headerName)
		if header != "" {
			return header
		}
		return c.IP() // Fallback to IP if header not present
	}
}

// CleanupExpiredLimiters removes old rate limiters (should be called periodically)
func (r *RateLimitMiddleware) CleanupExpiredLimiters() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	cutoff := time.Now().Add(-2 * r.window)
	
	for key, limiter := range r.limiters {
		limiter.mutex.Lock()
		if limiter.lastRefill.Before(cutoff) {
			delete(r.limiters, key)
		}
		limiter.mutex.Unlock()
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}