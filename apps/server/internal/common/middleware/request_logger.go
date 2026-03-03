package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
)

// ZapLogger creates an Echo middleware that logs requests using Zap
// Example usage:
//
//	e.Use(middleware.ZapLogger())
var logger, _ = zap.NewProduction()

func ZapLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogRemoteIP:  true,
		LogRequestID: true,
		LogURIPath:   true,
		LogMethod:    true,
		LogLatency:   true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("Incoming request",
				zap.Int("status", v.Status),
				zap.String("method", v.Method),
				zap.String("URI", v.URI),
				zap.Duration("latency", v.Latency),
				zap.String("request_id", v.RequestID),
				zap.String("remote_ip", v.RemoteIP),
			)

			return nil
		},
	})
}

// ZapRecovery creates an Echo middleware that recovers from panics and logs them using Zap
// Example usage:
//
//	e.Use(middleware.Recovery())
func Recovery() echo.MiddlewareFunc {
	return middleware.Recover()
}
