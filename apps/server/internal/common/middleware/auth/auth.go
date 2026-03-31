package auth

import (
	"net/http"

	"github.com/DSAwithGautam/Coderz.space/internal/common/response"
	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

const (
	ClaimsKey = "claims"
)

// middleware to check if the user is authenticated
func AuthMiddleware(jwtSecret string, jwtExpiryTime string) echo.MiddlewareFunc {
	echojwtConfig := echojwt.Config{
		SigningKey: []byte(jwtSecret),
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return &utils.TokenPayload{}
		},
		TokenLookup: "header:Authorization:Bearer ,cookie:access_token",
		ErrorHandler: func(c *echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"message": "INVALID_TOKEN",
			})
		},
	}

	// create the wrapped middleware
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// First apply Echo JWT middleware
		jwtMiddleware := echojwt.WithConfig(echojwtConfig)

		return jwtMiddleware(func(c *echo.Context) error {
			// Extract token from context (set by Echo JWT middleware)
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return response.NewResponse(c, http.StatusUnauthorized, "STATUS_UNAUTHORIZED", "INVALID_TOKEN_FORMAT", nil, nil)
			}

			// Parse claims into our custom struct
			claims, ok := token.Claims.(*utils.TokenPayload)
			if !ok {
				return response.NewResponse(c, http.StatusUnauthorized, "STATUS_UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
			}

			// Store parsed claims in context for easy access
			c.Set(ClaimsKey, claims)

			// Continue to next handler
			return next(c)
		})
	}
}
