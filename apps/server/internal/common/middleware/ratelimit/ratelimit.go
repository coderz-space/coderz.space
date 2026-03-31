package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/common/response"
	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
	"github.com/labstack/echo/v5"
)

// TokenBucket represents a token bucket for rate limiting
type TokenBucket struct {
	tokens         int
	capacity       int
	refillRate     int
	refillInterval time.Duration
	lastRefill     time.Time
	mu             sync.Mutex
}

// RateLimiter manages rate limiting for users
type RateLimiter struct {
	buckets map[string]*TokenBucket
	mu      sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*TokenBucket),
	}
}

// getBucket retrieves or creates a token bucket for a user
func (rl *RateLimiter) getBucket(userID string, capacity, refillRate int, refillInterval time.Duration) *TokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[userID]
	if !exists {
		bucket = &TokenBucket{
			tokens:         capacity,
			capacity:       capacity,
			refillRate:     refillRate,
			refillInterval: refillInterval,
			lastRefill:     time.Now(),
		}
		rl.buckets[userID] = bucket
	}

	return bucket
}

// Allow checks if a request is allowed based on rate limit
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	refills := int(elapsed / tb.refillInterval)

	if refills > 0 {
		tb.tokens += refills * tb.refillRate
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}

	// Check if we have tokens available
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// RateLimitMiddleware creates a rate limiting middleware
// capacity: maximum number of requests
// refillRate: number of tokens to add per interval
// refillInterval: time interval for refilling tokens
func RateLimitMiddleware(capacity, refillRate int, refillInterval time.Duration) echo.MiddlewareFunc {
	limiter := NewRateLimiter()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Extract user ID from auth context
			claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
			if !ok {
				// If no auth context, allow the request (auth middleware will handle it)
				return next(c)
			}

			// Get or create bucket for this user
			bucket := limiter.getBucket(claims.UserID, capacity, refillRate, refillInterval)

			// Check if request is allowed
			if !bucket.Allow() {
				return response.NewResponse(c, http.StatusTooManyRequests, "TOO_MANY_REQUESTS", "RATE_LIMIT_EXCEEDED", nil, nil)
			}

			return next(c)
		}
	}
}
