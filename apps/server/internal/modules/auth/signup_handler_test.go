package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coderz-space/coderz.space/internal/common/email"
	"github.com/coderz-space/coderz.space/internal/config"
	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v5"
)

func TestHandler_Signup(t *testing.T) {
	// Common config
	cfg := &config.Config{
		JWTSecret:           "test-secret",
		JWTExpires:          "1h",
		RefreshTokenExpires: time.Hour * 24,
	}

	tests := []struct {
		name           string
		body           interface{}
		setupMockDB    func() db.DBTX
		expectedStatus int
		checkResponse  func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "invalid json body",
			body: "invalid json", // Sent as string to force parse error
			setupMockDB: func() db.DBTX {
				return &mockDB{} // Won't be called
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &resp)
				if resp["message"] != "INVALID_REQUEST_BODY" {
					t.Errorf("Expected INVALID_REQUEST_BODY, got %v", resp["message"])
				}
			},
		},
		{
			name: "validation failure (invalid email)",
			body: SignupRequest{
				Name:     "Test User",
				Email:    "invalid-email",
				Password: "Password123",
			},
			setupMockDB: func() db.DBTX {
				return &mockDB{} // Won't be called
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &resp)
				if resp["message"] != "VALIDATION_FAILED" {
					t.Errorf("Expected VALIDATION_FAILED, got %v", resp["message"])
				}
			},
		},
		{
			name: "service layer failure (password complexity)",
			body: SignupRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password", // missing number
			},
			setupMockDB: func() db.DBTX {
				return &mockDB{} // Won't be called
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &resp)
				if resp["message"] != "VALIDATION_FAILED" {
					t.Errorf("Expected VALIDATION_FAILED, got %v", resp["message"])
				}
			},
		},
		{
			name: "service layer failure (email already exists)",
			body: SignupRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "Password123",
			},
			setupMockDB: func() db.DBTX {
				return &mockDB{
					queryRowFn: func(_ context.Context, _ string, _ ...interface{}) pgx.Row {
						return &mockRow{
							scanFn: func(_ ...any) error {
								return errors.New("duplicate key value violates unique constraint")
							},
						}
					},
				}
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &resp)
				if resp["message"] != "duplicate key value violates unique constraint" {
					t.Errorf("Expected database error message, got %v", resp["message"])
				}
			},
		},
		{
			name: "successful signup",
			body: SignupRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "Password123",
			},
			setupMockDB: func() db.DBTX {
				return &mockDB{
					queryRowFn: func(_ context.Context, _ string, _ ...interface{}) pgx.Row {
						return &mockRow{
							scanFn: func(_ ...any) error {
								return nil
							},
						}
					},
					execFn: func(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
						return pgconn.CommandTag{}, nil
					},
				}
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp AuthResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if !resp.Success {
					t.Errorf("Expected success=true, got %v", resp.Success)
				}

				// Verify cookies were set
				cookies := rec.Result().Cookies()
				var foundAccess, foundRefresh bool
				for _, c := range cookies {
					if c.Name == "access_token" {
						foundAccess = true
					}
					if c.Name == "refresh_token" {
						foundRefresh = true
					}
				}
				if !foundAccess || !foundRefresh {
					t.Errorf("Expected access_token and refresh_token cookies to be set")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup dependencies
			mockDB := tt.setupMockDB()
			queries := db.New(mockDB)
			service := NewService(queries, cfg, email.NewService(cfg))
			handler := NewHandler(service)

			// Setup Echo
			e := echo.New()
			var reqBody []byte
			if str, ok := tt.body.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/auth/signup", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute handler
			_ = handler.Signup(c)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}
