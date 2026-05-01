package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
)

type mockDB struct {
	execFn     func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	queryFn    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	queryRowFn func(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func (m *mockDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	if m.execFn != nil {
		return m.execFn(ctx, sql, args...)
	}
	return pgconn.CommandTag{}, nil
}

func (m *mockDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if m.queryFn != nil {
		return m.queryFn(ctx, sql, args...)
	}
	return nil, nil
}

func (m *mockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if m.queryRowFn != nil {
		return m.queryRowFn(ctx, sql, args...)
	}
	return &mockRow{}
}

type mockRow struct {
	scanFn func(dest ...any) error
}

func (m *mockRow) Scan(dest ...any) error {
	if m.scanFn != nil {
		return m.scanFn(dest...)
	}
	return nil
}

// MockQueries implements db.Queries interface for testing
type MockQueries struct {
	GetUserByEmailFunc                   func(ctx context.Context, email pgtype.Text) (db.User, error)
	DeleteUserPasswordResetTokensFunc    func(ctx context.Context, userID pgtype.UUID) error
	CreatePasswordResetTokenFunc         func(ctx context.Context, params db.CreatePasswordResetTokenParams) (db.PasswordResetToken, error)
	GetPasswordResetTokenFunc            func(ctx context.Context, tokenHash string) (db.PasswordResetToken, error)
	UpdateUserPasswordFunc               func(ctx context.Context, params db.UpdateUserPasswordParams) error
	DeletePasswordResetTokenFunc         func(ctx context.Context, tokenHash string) error
	DeleteUserRefreshTokensFunc          func(ctx context.Context, userID pgtype.UUID) error
	CreateUserFunc                       func(ctx context.Context, params db.CreateUserParams) (db.User, error)
	GetUserByIdFunc                      func(ctx context.Context, id pgtype.UUID) (db.User, error)
	CreateRefreshTokenFunc               func(ctx context.Context, params db.CreateRefreshTokenParams) (db.RefreshToken, error)
	GetRefreshTokenFunc                  func(ctx context.Context, tokenHash string) (db.RefreshToken, error)
	DeleteRefreshTokenFunc               func(ctx context.Context, tokenHash string) error
}

func (m *MockQueries) GetUserByEmail(ctx context.Context, email pgtype.Text) (db.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return db.User{}, nil
}

func (m *MockQueries) DeleteUserPasswordResetTokens(ctx context.Context, userID pgtype.UUID) error {
	if m.DeleteUserPasswordResetTokensFunc != nil {
		return m.DeleteUserPasswordResetTokensFunc(ctx, userID)
	}
	return nil
}

func (m *MockQueries) CreatePasswordResetToken(ctx context.Context, params db.CreatePasswordResetTokenParams) (db.PasswordResetToken, error) {
	if m.CreatePasswordResetTokenFunc != nil {
		return m.CreatePasswordResetTokenFunc(ctx, params)
	}
	return db.PasswordResetToken{}, nil
}

func (m *MockQueries) GetPasswordResetToken(ctx context.Context, tokenHash string) (db.PasswordResetToken, error) {
	if m.GetPasswordResetTokenFunc != nil {
		return m.GetPasswordResetTokenFunc(ctx, tokenHash)
	}
	return db.PasswordResetToken{}, nil
}

func (m *MockQueries) UpdateUserPassword(ctx context.Context, params db.UpdateUserPasswordParams) error {
	if m.UpdateUserPasswordFunc != nil {
		return m.UpdateUserPasswordFunc(ctx, params)
	}
	return nil
}

func (m *MockQueries) DeletePasswordResetToken(ctx context.Context, tokenHash string) error {
	if m.DeletePasswordResetTokenFunc != nil {
		return m.DeletePasswordResetTokenFunc(ctx, tokenHash)
	}
	return nil
}

func (m *MockQueries) DeleteUserRefreshTokens(ctx context.Context, userID pgtype.UUID) error {
	if m.DeleteUserRefreshTokensFunc != nil {
		return m.DeleteUserRefreshTokensFunc(ctx, userID)
	}
	return nil
}

func (m *MockQueries) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, params)
	}
	return db.User{}, nil
}

func (m *MockQueries) GetUserById(ctx context.Context, id pgtype.UUID) (db.User, error) {
	if m.GetUserByIdFunc != nil {
		return m.GetUserByIdFunc(ctx, id)
	}
	return db.User{}, nil
}

func (m *MockQueries) CreateRefreshToken(ctx context.Context, params db.CreateRefreshTokenParams) (db.RefreshToken, error) {
	if m.CreateRefreshTokenFunc != nil {
		return m.CreateRefreshTokenFunc(ctx, params)
	}
	return db.RefreshToken{}, nil
}

func (m *MockQueries) GetRefreshToken(ctx context.Context, tokenHash string) (db.RefreshToken, error) {
	if m.GetRefreshTokenFunc != nil {
		return m.GetRefreshTokenFunc(ctx, tokenHash)
	}
	return db.RefreshToken{}, nil
}

func (m *MockQueries) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	if m.DeleteRefreshTokenFunc != nil {
		return m.DeleteRefreshTokenFunc(ctx, tokenHash)
	}
	return nil
}
