package auth

import "github.com/jackc/pgx/v5/pgtype"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
}

type AuthUser struct {
	ID            pgtype.UUID `json:"id"`
	Name          string      `json:"name"`
	Email         string      `json:"email"`
	EmailVerified bool        `json:"emailVerified"`
}

type AuthResponseData struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	User         AuthUser `json:"user"`
}

type AuthResponse struct {
	Success bool             `json:"success"`
	Data    AuthResponseData `json:"data"`
}

type RefreshResponseData struct {
	AccessToken string `json:"accessToken"`
}

type RefreshResponse struct {
	Success bool                `json:"success"`
	Data    RefreshResponseData `json:"data"`
}

type UserProfileResponse struct {
	Success bool     `json:"success"`
	Data    AuthUser `json:"data"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=50"`
}

type GenericResponse struct {
	Success bool           `json:"success"`
	Data    map[string]any `json:"data"`
}
