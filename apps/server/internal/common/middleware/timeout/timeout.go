package timeout

import (
	"context"
	"net/http"
	"time"

	"github.com/coderz-space/coderz.space/internal/common/response"
	"github.com/labstack/echo/v5"
)

// TimeoutMiddleware creates a middleware that enforces request timeout
// to prevent resource exhaustion from long-running requests
func TimeoutMiddleware(timeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Create a context with timeout
			ctx, cancel := context.WithTimeout((*c).Request().Context(), timeout)
			defer cancel()

			// Set the timeout context on the request
			req := (*c).Request().WithContext(ctx)
			(*c).SetRequest(req)

			// Channel to capture the result of the handler
			done := make(chan error, 1)

			// Run the handler in a goroutine
			go func() {
				done <- next(c)
			}()

			// Wait for either the handler to complete or timeout
			select {
			case err := <-done:
				return err
			case <-ctx.Done():
				// Timeout occurred
				if ctx.Err() == context.DeadlineExceeded {
					return response.NewResponse(c, http.StatusRequestTimeout, "REQUEST_TIMEOUT", "REQUEST_EXCEEDED_TIMEOUT", nil, nil)
				}
				return ctx.Err()
			}
		}
	}
}
