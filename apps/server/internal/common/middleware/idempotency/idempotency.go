package idempotency

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/common/response"
	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
	"github.com/labstack/echo/v5"
)

const (
	// IdempotencyKeyHeader is the header name for idempotency key
	IdempotencyKeyHeader = "Idempotency-Key"

	// DefaultTTL is the default time-to-live for idempotency keys
	DefaultTTL = 24 * time.Hour
)

// CachedResponse represents a cached response for idempotency
type CachedResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
	Timestamp  time.Time
}

// IdempotencyStore manages idempotency keys and cached responses
type IdempotencyStore struct {
	cache map[string]*CachedResponse
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewIdempotencyStore creates a new idempotency store
func NewIdempotencyStore(ttl time.Duration) *IdempotencyStore {
	store := &IdempotencyStore{
		cache: make(map[string]*CachedResponse),
		ttl:   ttl,
	}

	// Start cleanup goroutine
	go store.cleanup()

	return store
}

// cleanup removes expired entries periodically
func (s *IdempotencyStore) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, resp := range s.cache {
			if now.Sub(resp.Timestamp) > s.ttl {
				delete(s.cache, key)
			}
		}
		s.mu.Unlock()
	}
}

// Get retrieves a cached response
func (s *IdempotencyStore) Get(key string) (*CachedResponse, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resp, exists := s.cache[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Since(resp.Timestamp) > s.ttl {
		return nil, false
	}

	return resp, true
}

// Set stores a cached response
func (s *IdempotencyStore) Set(key string, resp *CachedResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[key] = resp
}

// generateKey generates a scoped idempotency key
func generateKey(userID, endpoint, idempotencyKey string) string {
	data := fmt.Sprintf("%s:%s:%s", userID, endpoint, idempotencyKey)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// IdempotencyMiddleware creates an idempotency middleware
func IdempotencyMiddleware(store *IdempotencyStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Only apply to POST requests
			if c.Request().Method != http.MethodPost {
				return next(c)
			}

			// Get idempotency key from header
			idempotencyKey := c.Request().Header.Get(IdempotencyKeyHeader)
			if idempotencyKey == "" {
				// No idempotency key, proceed normally
				return next(c)
			}

			// Validate idempotency key format (should be UUID or similar)
			if len(idempotencyKey) < 16 || len(idempotencyKey) > 128 {
				return response.BadRequestError(&c, "INVALID_IDEMPOTENCY_KEY", "Idempotency key must be between 16 and 128 characters")
			}

			// Get user ID from context
			claims, ok := c.Get(auth.ClaimsKey).(*utils.TokenPayload)
			if !ok {
				// No auth context, can't scope idempotency key
				return next(c)
			}

			// Generate scoped key
			endpoint := c.Request().URL.Path
			scopedKey := generateKey(claims.UserID, endpoint, idempotencyKey)

			// Check if we have a cached response
			if cachedResp, exists := store.Get(scopedKey); exists {
				// Return cached response
				for key, value := range cachedResp.Headers {
					c.Response().Header().Set(key, value)
				}
				return c.JSONBlob(cachedResp.StatusCode, cachedResp.Body)
			}

			// Create a response recorder
			rec := &responseRecorder{
				ResponseWriter: c.Response().Writer,
				statusCode:     http.StatusOK,
				body:           []byte{},
				headers:        make(map[string]string),
			}
			c.Response().Writer = rec

			// Process request
			err := next(c)

			// Cache successful responses (2xx status codes)
			if rec.statusCode >= 200 && rec.statusCode < 300 {
				// Capture important headers
				for _, header := range []string{"Content-Type", "Content-Length"} {
					if value := c.Response().Header().Get(header); value != "" {
						rec.headers[header] = value
					}
				}

				store.Set(scopedKey, &CachedResponse{
					StatusCode: rec.statusCode,
					Body:       rec.body,
					Headers:    rec.headers,
					Timestamp:  time.Now(),
				})
			}

			return err
		}
	}
}

// responseRecorder records the response for caching
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
	headers    map[string]string
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}
