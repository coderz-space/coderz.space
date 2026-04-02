package auth

import (
	"context"

	"github.com/coderz-space/coderz.space/internal/config"
	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// MockQuerier implements db.Querier for testing
type MockQuerier struct {
	db.Querier
	GetPasswordResetTokenFunc    func(ctx context.Context, tokenHash string) (db.PasswordResetToken, error)
	UpdateUserPasswordFunc       func(ctx context.Context, arg db.UpdateUserPasswordParams) error
	DeletePasswordResetTokenFunc func(ctx context.Context, tokenHash string) error
	DeleteUserRefreshTokensFunc  func(ctx context.Context, userID pgtype.UUID) error
	GetUserByIdFunc              func(ctx context.Context, id pgtype.UUID) (db.User, error)
}

func (m *MockQuerier) GetPasswordResetToken(ctx context.Context, tokenHash string) (db.PasswordResetToken, error) {
	if m.GetPasswordResetTokenFunc != nil {
		return m.GetPasswordResetTokenFunc(ctx, tokenHash)
	}
	return db.PasswordResetToken{}, nil
}

func (m *MockQuerier) UpdateUserPassword(ctx context.Context, arg db.UpdateUserPasswordParams) error {
	if m.UpdateUserPasswordFunc != nil {
		return m.UpdateUserPasswordFunc(ctx, arg)
	}
	return nil
}

func (m *MockQuerier) DeletePasswordResetToken(ctx context.Context, tokenHash string) error {
	if m.DeletePasswordResetTokenFunc != nil {
		return m.DeletePasswordResetTokenFunc(ctx, tokenHash)
	}
	return nil
}

func (m *MockQuerier) DeleteUserRefreshTokens(ctx context.Context, userID pgtype.UUID) error {
	if m.DeleteUserRefreshTokensFunc != nil {
		return m.DeleteUserRefreshTokensFunc(ctx, userID)
	}
	return nil
}

func (m *MockQuerier) GetUserById(ctx context.Context, id pgtype.UUID) (db.User, error) {
	if m.GetUserByIdFunc != nil {
		return m.GetUserByIdFunc(ctx, id)
	}
	return db.User{}, nil
}

func setupTestServiceWithMock(q db.Querier) *Service {
	return NewService(q, &config.Config{})
}
