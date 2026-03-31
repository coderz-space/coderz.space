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

// Signup godoc
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body SignupRequest true "User registration details"
// @Success 201 {object} AuthResponse "User registered successfully"
// @Failure 400 {object} map[string]any "Bad request - validation error or email already exists"
// @Router /v1/auth/signup [post]

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

// Login godoc
// @Summary Authenticate user
// @Description Login with email and password to receive authentication tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse "Login successful"
// @Failure 401 {object} map[string]any "Unauthorized - invalid credentials"
// @Router /v1/auth/login [post]

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

// Refresh godoc
// @Summary Refresh access token
// @Description Generate a new access token using refresh token from cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} RefreshResponse "Token refreshed successfully"
// @Failure 401 {object} map[string]any "Unauthorized - missing or invalid refresh token"
// @Router /v1/auth/refresh [post]

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

// Logout godoc
// @Summary Logout user
// @Description Logout user and revoke refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} GenericResponse "Logout successful"
// @Router /v1/auth/logout [post]

func (h *Handler) Logout(c *echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err == nil {
		_ = h.service.Logout(c.Request().Context(), cookie.Value)
	}

	h.clearAuthCookies(c)

	return c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Data:    map[string]any{},
	})
}

// Me godoc
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserProfileResponse "User profile retrieved successfully"
// @Failure 401 {object} map[string]any "Unauthorized - invalid or missing token"
// @Failure 404 {object} map[string]any "Not found - user does not exist"
// @Router /v1/auth/me [get]

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

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset token (always returns success to prevent email enumeration)
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body ForgotPasswordRequest true "Email address"
// @Success 200 {object} GenericResponse "Password reset email sent (if email exists)"
// @Failure 400 {object} map[string]any "Bad request - validation error"
// @Router /v1/auth/forgot-password [post]
func (h *Handler) ForgotPassword(c *echo.Context, body ForgotPasswordRequest) error {
	// Always return success to prevent email enumeration
	_ = h.service.ForgotPassword(c.Request().Context(), body)

	return c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Data:    map[string]any{},
	})
}

// ResetPassword godoc
// @Summary Reset password with token
// @Description Reset user password using a valid reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} GenericResponse "Password reset successful"
// @Failure 400 {object} map[string]any "Bad request - validation error or invalid/expired token"
// @Router /v1/auth/reset-password [post]
func (h *Handler) ResetPassword(c *echo.Context, body ResetPasswordRequest) error {
	err := h.service.ResetPassword(c.Request().Context(), body)
	if err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	}

	return c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Data:    map[string]any{},
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
