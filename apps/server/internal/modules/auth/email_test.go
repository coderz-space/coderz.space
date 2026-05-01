package auth

import (
	"context"
	"testing"

	"github.com/coderz-space/coderz.space/internal/common/email"
	"github.com/coderz-space/coderz.space/internal/config"
	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

// MockEmailService implements email.Service for testing
type MockEmailService struct {
	SentEmails []EmailRecord
	ErrorOnSend bool
}

type EmailRecord struct {
	To         string
	ResetToken string
}

func (m *MockEmailService) SendPasswordResetEmail(to string, resetToken string) error {
	if m.ErrorOnSend {
		return assert.AnError
	}
	m.SentEmails = append(m.SentEmails, EmailRecord{To: to, ResetToken: resetToken})
	return nil
}

// TestForgotPasswordSendsEmail verifies that ForgotPassword calls the email service
func TestForgotPasswordSendsEmail(t *testing.T) {
	mockQueries := &MockQueries{
		GetUserByEmailFunc: func(ctx context.Context, email pgtype.Text) (db.User, error) {
			return db.User{
				ID:    pgtype.UUID{Bytes: [16]byte{1, 2, 3, 4}, Valid: true},
				Email: email,
				Name:  "Test User",
			}, nil
		},
		DeleteUserPasswordResetTokensFunc: func(ctx context.Context, userID pgtype.UUID) error {
			return nil
		},
		CreatePasswordResetTokenFunc: func(ctx context.Context, params db.CreatePasswordResetTokenParams) (db.PasswordResetToken, error) {
			return db.PasswordResetToken{
				UserID:    params.UserID,
				TokenHash: params.TokenHash,
				ExpiresAt: params.ExpiresAt,
			}, nil
		},
	}

	mockEmailService := &MockEmailService{}

	cfg := &config.Config{
		FrontendOrigin: "http://localhost:3000",
		Environment:   config.Development,
	}

	service := NewService(mockQueries, cfg, mockEmailService)

	err := service.ForgotPassword(context.Background(), ForgotPasswordRequest{
		Email: "test@example.com",
	})

	assert.NoError(t, err)
	assert.Len(t, mockEmailService.SentEmails, 1)
	assert.Equal(t, "test@example.com", mockEmailService.SentEmails[0].To)
	assert.NotEmpty(t, mockEmailService.SentEmails[0].ResetToken)
}

// TestForgotPasswordWithoutEmail verifies email enumeration prevention
func TestForgotPasswordWithoutEmail(t *testing.T) {
	mockQueries := &MockQueries{
		GetUserByEmailFunc: func(ctx context.Context, email pgtype.Text) (db.User, error) {
			return db.User{}, assert.AnError // User not found
		},
	}

	mockEmailService := &MockEmailService{}

	cfg := &config.Config{
		FrontendOrigin: "http://localhost:3000",
		Environment:   config.Development,
	}

	service := NewService(mockQueries, cfg, mockEmailService)

	err := service.ForgotPassword(context.Background(), ForgotPasswordRequest{
		Email: "nonexistent@example.com",
	})

	// Should not error to prevent email enumeration
	assert.NoError(t, err)
	// Email service should not be called
	assert.Len(t, mockEmailService.SentEmails, 0)
}

// TestForgotPasswordEmailServiceFailure verifies graceful handling of email service errors
func TestForgotPasswordEmailServiceFailure(t *testing.T) {
	mockQueries := &MockQueries{
		GetUserByEmailFunc: func(ctx context.Context, email pgtype.Text) (db.User, error) {
			return db.User{
				ID:    pgtype.UUID{Bytes: [16]byte{1, 2, 3, 4}, Valid: true},
				Email: email,
				Name:  "Test User",
			}, nil
		},
		DeleteUserPasswordResetTokensFunc: func(ctx context.Context, userID pgtype.UUID) error {
			return nil
		},
		CreatePasswordResetTokenFunc: func(ctx context.Context, params db.CreatePasswordResetTokenParams) (db.PasswordResetToken, error) {
			return db.PasswordResetToken{
				UserID:    params.UserID,
				TokenHash: params.TokenHash,
				ExpiresAt: params.ExpiresAt,
			}, nil
		},
	}

	mockEmailService := &MockEmailService{ErrorOnSend: true}

	cfg := &config.Config{
		FrontendOrigin: "http://localhost:3000",
		Environment:   config.Development,
	}

	service := NewService(mockQueries, cfg, mockEmailService)

	err := service.ForgotPassword(context.Background(), ForgotPasswordRequest{
		Email: "test@example.com",
	})

	// Should not fail even if email service fails - reset token should still be created
	assert.NoError(t, err)
	// Email service should have attempted to send
	assert.Len(t, mockEmailService.SentEmails, 1)
}
