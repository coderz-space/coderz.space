package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coderz-space/coderz.space/internal/common/utils"
	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
)

func TestResetPassword_BindError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/reset-password", bytes.NewReader([]byte(`invalid json`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewHandler(setupTestServiceWithMock(&MockQuerier{}))
	err := handler.ResetPassword(c)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestResetPassword_ServiceSuccess(t *testing.T) {
	e := echo.New()

	bodyBytes, _ := json.Marshal(map[string]string{
		"token":       "valid-token",
		"newPassword": "ValidPassword123",
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/reset-password", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock querier simulating a successful reset password
	userID, _ := utils.StringToUUID("11111111-1111-1111-1111-111111111111")
	mockQ := &MockQuerier{
		GetPasswordResetTokenFunc: func(ctx context.Context, tokenHash string) (db.PasswordResetToken, error) {
			return db.PasswordResetToken{
				UserID: userID,
			}, nil
		},
		UpdateUserPasswordFunc: func(ctx context.Context, arg db.UpdateUserPasswordParams) error {
			return nil
		},
		DeletePasswordResetTokenFunc: func(ctx context.Context, tokenHash string) error {
			return nil
		},
		DeleteUserRefreshTokensFunc: func(ctx context.Context, id pgtype.UUID) error {
			return nil
		},
	}

	handler := NewHandler(setupTestServiceWithMock(mockQ))
	err := handler.ResetPassword(c)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	var res GenericResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	if !res.Success {
		t.Errorf("Expected success response, got %v", res.Success)
	}
}

func TestResetPassword_ServiceFailure(t *testing.T) {
	e := echo.New()

	bodyBytes, _ := json.Marshal(map[string]string{
		"token":       "invalid-token",
		"newPassword": "ValidPassword123",
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/reset-password", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock querier simulating an expired/invalid token
	mockQ := &MockQuerier{
		GetPasswordResetTokenFunc: func(ctx context.Context, tokenHash string) (db.PasswordResetToken, error) {
			return db.PasswordResetToken{}, errors.New("not found")
		},
	}

	handler := NewHandler(setupTestServiceWithMock(mockQ))
	err := handler.ResetPassword(c)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
