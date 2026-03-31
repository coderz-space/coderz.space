package auth

import "github.com/jackc/pgx/v5/pgtype"

// LoginRequest represents the login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8,max=50,password_complexity" example:"Password123"`
}

// SignupRequest represents the user registration data
type SignupRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8,max=50,password_complexity" example:"Password123"`
	Name     string `json:"name" validate:"required,min=2,max=100" example:"John Doe"`
}

// AuthUser represents the authenticated user data
type AuthUser struct {
	ID            pgtype.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	EmailVerified bool        `json:"emailVerified" example:"false"`
	Name          string      `json:"name" example:"John Doe"`
	Email         string      `json:"email" example:"user@example.com"`
}

// AuthResponseData contains authentication tokens and user data
type AuthResponseData struct {
	AccessToken  string   `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string   `json:"refreshToken" example:"a1b2c3d4e5f6..."`
	User         AuthUser `json:"user"`
}

// AuthResponse is the response for login and signup
type AuthResponse struct {
	Data    AuthResponseData `json:"data"`
	Success bool             `json:"success" example:"true"`
}

// RefreshResponseData contains the new access token
type RefreshResponseData struct {
	AccessToken string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RefreshResponse is the response for token refresh
type RefreshResponse struct {
	Data    RefreshResponseData `json:"data"`
	Success bool                `json:"success" example:"true"`
}

// UserProfileResponse is the response for user profile
type UserProfileResponse struct {
	Data    AuthUser `json:"data"`
	Success bool     `json:"success" example:"true"`
}

// ForgotPasswordRequest represents the forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
}

// ResetPasswordRequest represents the password reset request
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required" example:"a1b2c3d4e5f6g7h8i9j0"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=50,password_complexity" example:"NewPassword123"`
}

// GenericResponse is a generic success response
type GenericResponse struct {
	Data    map[string]any `json:"data"`
	Success bool           `json:"success" example:"true"`
}
