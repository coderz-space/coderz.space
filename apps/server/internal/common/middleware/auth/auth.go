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
func AuthMiddleware(jwtSecret, jwtExpiryTime string) echo.MiddlewareFunc {
	echojwtConfig := echojwt.Config{
		SigningKey: []byte(jwtSecret),
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return &utils.TokenPayload{}
		},
		// Prioritize header over cookie by listing header first
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

// CookieConfig holds cookie configuration
type CookieConfig struct {
	Domain   string
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
	Path     string
}

// DefaultCookieConfig returns default cookie configuration
func DefaultCookieConfig() CookieConfig {
	return CookieConfig{
		Domain:   "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
}

// SetAccessTokenCookie sets the access token cookie with secure flags
func SetAccessTokenCookie(c *echo.Context, token string, maxAge int, config CookieConfig) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     config.Path,
		Domain:   config.Domain,
		MaxAge:   maxAge,
		Secure:   config.Secure,
		HttpOnly: config.HttpOnly,
		SameSite: config.SameSite,
	}
	c.SetCookie(cookie)
}

// SetRefreshTokenCookie sets the refresh token cookie with secure flags
func SetRefreshTokenCookie(c *echo.Context, token string, maxAge int, config CookieConfig) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     config.Path,
		Domain:   config.Domain,
		MaxAge:   maxAge,
		Secure:   config.Secure,
		HttpOnly: config.HttpOnly,
		SameSite: config.SameSite,
	}
	c.SetCookie(cookie)
}

// ClearAuthCookies clears both access and refresh token cookies
func ClearAuthCookies(c *echo.Context, config CookieConfig) {
	// Clear access token
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     config.Path,
		Domain:   config.Domain,
		MaxAge:   -1,
		Secure:   config.Secure,
		HttpOnly: config.HttpOnly,
		SameSite: config.SameSite,
	}
	c.SetCookie(accessCookie)

	// Clear refresh token
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     config.Path,
		Domain:   config.Domain,
		MaxAge:   -1,
		Secure:   config.Secure,
		HttpOnly: config.HttpOnly,
		SameSite: config.SameSite,
	}
	c.SetCookie(refreshCookie)
}
