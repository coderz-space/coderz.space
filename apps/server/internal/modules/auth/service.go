package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
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
