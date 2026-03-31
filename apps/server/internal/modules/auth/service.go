package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
	"github.com/DSAwithGautam/Coderz.space/internal/config"
	db "github.com/DSAwithGautam/Coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries *db.Queries
	config  *config.Config
}

func NewService(queries *db.Queries, config *config.Config) *Service {
	return &Service{queries: queries, config: config}
}

func (s *Service) Signup(ctx context.Context, req SignupRequest) (*AuthResponseData, error) {
	// Validate password complexity
	if !s.validatePasswordComplexity(req.Password) {
		return nil, errors.New("PASSWORD_MUST_CONTAIN_LETTER_AND_NUMBER")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Name:         req.Name,
		Email:        pgtype.Text{String: req.Email, Valid: true},
		PasswordHash: pgtype.Text{String: string(hashedPassword), Valid: true},
		Role:         db.UserRoleUser,
	})
	if err != nil {
		return nil, err
	}

	return s.generateAuthData(ctx, user)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponseData, error) {
	user, err := s.queries.GetUserByEmail(ctx, pgtype.Text{String: req.Email, Valid: true})
	if err != nil {
		return nil, errors.New("INVALID_CREDENTIALS")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password)); err != nil {
		return nil, errors.New("INVALID_CREDENTIALS")
	}

	return s.generateAuthData(ctx, user)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*AuthResponseData, error) {
	tokenHash := utils.HashString(refreshToken)
	rt, err := s.queries.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("INVALID_REFRESH_TOKEN")
	}

	if rt.ExpiresAt.Time.Before(time.Now()) {
		s.queries.DeleteRefreshToken(ctx, tokenHash)
		return nil, errors.New("EXPIRED_REFRESH_TOKEN")
	}

	user, err := s.queries.GetUserById(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	// Delete old refresh token (rotation)
	s.queries.DeleteRefreshToken(ctx, tokenHash)

	return s.generateAuthData(ctx, user)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := utils.HashString(refreshToken)
	return s.queries.DeleteRefreshToken(ctx, tokenHash)
}

func (s *Service) GetUserByID(ctx context.Context, userID pgtype.UUID) (*AuthUser, error) {
	user, err := s.queries.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &AuthUser{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email.String,
		EmailVerified: user.EmailVerified,
	}, nil
}

func (s *Service) generateAuthData(ctx context.Context, user db.User) (*AuthResponseData, error) {
	// Generate Access Token
	payload := utils.TokenPayload{
		UserID:   utils.UUIDToString(user.ID),
		Email:    user.Email.String,
		Role:     string(user.Role),
		UserName: user.Name,
	}

	accessToken, err := utils.GenerateToken(payload, s.config.JWT_EXPIRES)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token
	refreshToken, err := s.generateRandomString(32)
	if err != nil {
		return nil, err
	}

	tokenHash := utils.HashString(refreshToken)
	expiresAt := time.Now().Add(s.config.REFRESH_TOKEN_EXPIRES)

	_, err = s.queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &AuthResponseData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: AuthUser{
			ID:            user.ID,
			Name:          user.Name,
			Email:         user.Email.String,
			EmailVerified: user.EmailVerified,
		},
	}, nil
}

func (s *Service) generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *Service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	// Get user by email
	user, err := s.queries.GetUserByEmail(ctx, pgtype.Text{String: req.Email, Valid: true})
	if err != nil {
		// Silently fail to prevent email enumeration
		return nil
	}

	// Delete any existing password reset tokens for this user
	_ = s.queries.DeleteUserPasswordResetTokens(ctx, user.ID)

	// Generate reset token
	resetToken, err := s.generateRandomString(32)
	if err != nil {
		return err
	}

	// Hash the token before storing
	tokenHash := utils.HashString(resetToken)

	// Token expires in 1 hour
	expiresAt := time.Now().Add(1 * time.Hour)

	// Store the token
	_, err = s.queries.CreatePasswordResetToken(ctx, db.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return err
	}

	// Mock email sending for now (in production, integrate with SendGrid/SMTP)
	fmt.Printf("==== MOCK EMAIL SENDER ====\n")
	fmt.Printf("To: %s\n", req.Email)
	fmt.Printf("Subject: Password Reset Request\n")
	fmt.Printf("Body: Click here to reset your password: https://coderz.space/reset-password?token=%s\n", resetToken)
	fmt.Printf("===========================\n")

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	// Validate password complexity (at least 1 letter and 1 number)
	if !s.validatePasswordComplexity(req.NewPassword) {
		return errors.New("PASSWORD_MUST_CONTAIN_LETTER_AND_NUMBER")
	}

	// Hash the token to look it up
	tokenHash := utils.HashString(req.Token)

	// Get the reset token (only if not expired)
	resetToken, err := s.queries.GetPasswordResetToken(ctx, tokenHash)
	if err != nil {
		return errors.New("INVALID_OR_EXPIRED_TOKEN")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update user password
	err = s.queries.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           resetToken.UserID,
		PasswordHash: pgtype.Text{String: string(hashedPassword), Valid: true},
	})
	if err != nil {
		return err
	}

	// Delete the used reset token
	err = s.queries.DeletePasswordResetToken(ctx, tokenHash)
	if err != nil {
		return err
	}

	// Delete all refresh tokens for this user (force re-login)
	_ = s.queries.DeleteUserRefreshTokens(ctx, resetToken.UserID)

	return nil
}

func (s *Service) validatePasswordComplexity(password string) bool {
	hasLetter := false
	hasNumber := false

	for _, char := range password {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			hasLetter = true
		}
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
		if hasLetter && hasNumber {
			return true
		}
	}

	return hasLetter && hasNumber
}
