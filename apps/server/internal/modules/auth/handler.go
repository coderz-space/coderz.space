package auth

import (
	"net/http"

	"github.com/DSAwithGautam/Coderz.space/internal/common/middleware/auth"
	"github.com/DSAwithGautam/Coderz.space/internal/common/response"
	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Signup(c *echo.Context, body SignupRequest) error {
	data, err := h.service.Signup(c.Request().Context(), body)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	h.setAuthCookies(c, data.AccessToken, data.RefreshToken)

	return c.JSON(http.StatusCreated, AuthResponse{
		Success: true,
		Data:    *data,
	})
}

func (h *Handler) Login(c *echo.Context, body LoginRequest) error {
	data, err := h.service.Login(c.Request().Context(), body)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil, nil)
	}

	h.setAuthCookies(c, data.AccessToken, data.RefreshToken)

	return c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Data:    *data,
	})
}

func (h *Handler) Refresh(c *echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "MISSING_REFRESH_TOKEN", nil, nil)
	}

	data, err := h.service.Refresh(c.Request().Context(), cookie.Value)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil, nil)
	}

	h.setAuthCookies(c, data.AccessToken, data.RefreshToken)

	return c.JSON(http.StatusOK, RefreshResponse{
		Success: true,
		Data: RefreshResponseData{
			AccessToken: data.AccessToken,
		},
	})
}

func (h *Handler) Logout(c *echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err == nil {
		h.service.Logout(c.Request().Context(), cookie.Value)
	}

	h.clearAuthCookies(c)

	return c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Data:    map[string]any{},
	})
}

func (h *Handler) Me(c *echo.Context) error {
	claims, ok := (*c).Get(auth.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_TOKEN_CLAIMS", nil, nil)
	}

	userID, err := utils.StringToUUID(claims.UserID)
	if err != nil {
		return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "INVALID_USER_ID", nil, nil)
	}

	user, err := h.service.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", "USER_NOT_FOUND", nil, nil)
	}

	return c.JSON(http.StatusOK, UserProfileResponse{
		Success: true,
		Data:    *user,
	})
}

func (h *Handler) setAuthCookies(c *echo.Context, accessToken, refreshToken string) {
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   900, // 15 minutes
	}
	c.SetCookie(accessCookie)

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.service.config.REFRESH_TOKEN_EXPIRES.Seconds()),
	}
	c.SetCookie(refreshCookie)
}

func (h *Handler) clearAuthCookies(c *echo.Context) {
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(accessCookie)

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(refreshCookie)
}
