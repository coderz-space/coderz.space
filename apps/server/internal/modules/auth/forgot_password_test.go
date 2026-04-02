package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"

	"github.com/coderz-space/coderz.space/internal/config"
	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

type MockDBTX struct{}

func (m *MockDBTX) Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (m *MockDBTX) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	return nil, nil
}
func (m *MockDBTX) QueryRow(context.Context, string, ...interface{}) pgx.Row {
	return &mockRow{scanFn: func(dest ...any) error { return errors.New("user not found") }}
}

func TestForgotPasswordHandler(t *testing.T) {
	// Setup service with mocked DB queries and config
	mockDB := &MockDBTX{}
	queries := db.New(mockDB)
	cfg := &config.Config{}

	service := NewService(queries, cfg)
	handler := NewHandler(service)

	t.Run("valid email returns 200 success", func(t *testing.T) {
		reqBody := ForgotPasswordRequest{
			Email: "test@example.com",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/forgot-password", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.ForgotPassword(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		success, ok := resp["success"].(bool)
		assert.True(t, ok)
		assert.True(t, success)
	})

	t.Run("invalid email returns 400 validation error", func(t *testing.T) {
		reqBody := ForgotPasswordRequest{
			Email: "invalid-email",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/forgot-password", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.ForgotPassword(c)
		assert.NoError(t, err) // c.JSON returns nil
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "VALIDATION_FAILED", resp["message"])
		assert.Equal(t, "VALIDATION_ERROR", resp["status"])
		if val, exists := resp["success"]; exists {
			assert.Equal(t, false, val)
		} else {
			// omitempty might have stripped it, which is equivalent to false in this context
		}
	})

	t.Run("invalid request body returns 400 bad request", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/forgot-password", bytes.NewReader([]byte(`{"email": 123}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.ForgotPassword(c)
		assert.NoError(t, err) // c.JSON returns nil
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "INVALID_REQUEST_BODY", resp["message"])
		assert.Equal(t, "BAD_REQUEST", resp["status"])
		if val, exists := resp["success"]; exists {
			assert.Equal(t, false, val)
		} else {
			// omitempty might have stripped it, which is equivalent to false in this context
		}
	})
}
