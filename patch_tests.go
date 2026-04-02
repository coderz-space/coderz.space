package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	path := "apps/server/internal/modules/auth/handler_test.go"
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	newImports := `package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coderz-space/coderz.space/internal/common/middleware/auth"
	"github.com/coderz-space/coderz.space/internal/common/utils"
	"github.com/coderz-space/coderz.space/internal/config"
	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

// MockQuerier implements db.Querier for testing
type MockQuerier struct {
	db.Querier
	GetUserByIdFunc func(ctx context.Context, id pgtype.UUID) (db.User, error)
}

func (m *MockQuerier) GetUserById(ctx context.Context, id pgtype.UUID) (db.User, error) {
	if m.GetUserByIdFunc != nil {
		return m.GetUserByIdFunc(ctx, id)
	}
	return db.User{}, nil
}
`

	strContent := string(content)
	strContent = strings.Replace(strContent, "package auth\n\nimport (\n\t\"testing\"\n)", newImports, 1)

	// Replace TestMeAuthentication
	testAuth := `func TestMeAuthentication(t *testing.T) {
	e := echo.New()

	validUUIDStr := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name           string
		scenario       string
		setupContext   func(c echo.Context)
		expectedError  string
		expectedStatus int
	}{
		{
			name:           "missing claims fails",
			scenario:       "no claims in context",
			setupContext:   func(c echo.Context) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "INVALID_TOKEN_CLAIMS",
		},
		{
			name:           "invalid type for claims fails",
			scenario:       "claims is not *utils.TokenPayload",
			setupContext:   func(c echo.Context) {
				c.Set(auth.ClaimsKey, "invalid claims")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "INVALID_TOKEN_CLAIMS",
		},
		{
			name:           "invalid user ID in claims fails",
			scenario:       "claims has invalid UUID format",
			setupContext:   func(c echo.Context) {
				payload := &utils.TokenPayload{UserID: "invalid-uuid"}
				c.Set(auth.ClaimsKey, payload)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "INVALID_USER_ID",
		},
		{
			name:           "authenticated user can get profile",
			scenario:       "valid JWT token with claims",
			setupContext:   func(c echo.Context) {
				payload := &utils.TokenPayload{UserID: validUUIDStr}
				c.Set(auth.ClaimsKey, payload)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/auth/me", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tt.setupContext(c)

			mockQuerier := &MockQuerier{
				GetUserByIdFunc: func(ctx context.Context, id pgtype.UUID) (db.User, error) {
					return db.User{
						ID:            id,
						Name:          "Test User",
						Email:         pgtype.Text{String: "test@example.com", Valid: true},
						EmailVerified: true,
					}, nil
				},
			}
			service := NewService(mockQuerier, &config.Config{})
			handler := NewHandler(service)

			err := handler.Me(&c)
			if err != nil {
				// if Echo error handling returns an error
				t.Fatalf("Unexpected error: %v", err)
			}

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var resp map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.Equal(t, tt.expectedError, resp["message"])
			}
		})
	}
}`

	// Remove old TestMeAuthentication and replace
	startIndex := strings.Index(strContent, "func TestMeAuthentication(t *testing.T) {")
	if startIndex != -1 {
		endIndex := strings.Index(strContent[startIndex:], "}\n}\n")
		if endIndex != -1 {
			strContent = strContent[:startIndex] + testAuth + "\n" + strContent[startIndex+endIndex+4:]
		}
	}


	// Replace TestMeUserNotFound
	testNotFound := `func TestMeUserNotFound(t *testing.T) {
	e := echo.New()
	validUUIDStr := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name           string
		scenario       string
		setupQuerier   func(m *MockQuerier)
		expectedError  string
		expectedStatus int
	}{
		{
			name:           "existing user returns profile",
			scenario:       "user_id from token exists in database",
			setupQuerier: func(m *MockQuerier) {
				m.GetUserByIdFunc = func(ctx context.Context, id pgtype.UUID) (db.User, error) {
					return db.User{
						ID:            id,
						Name:          "Test User",
						Email:         pgtype.Text{String: "test@example.com", Valid: true},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "deleted user returns 404",
			scenario:       "user_id from token does not exist",
			setupQuerier: func(m *MockQuerier) {
				m.GetUserByIdFunc = func(ctx context.Context, id pgtype.UUID) (db.User, error) {
					return db.User{}, errors.New("not found")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "USER_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/auth/me", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			payload := &utils.TokenPayload{UserID: validUUIDStr}
			c.Set(auth.ClaimsKey, payload)

			mockQuerier := &MockQuerier{}
			tt.setupQuerier(mockQuerier)

			service := NewService(mockQuerier, &config.Config{})
			handler := NewHandler(service)

			err := handler.Me(&c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var resp map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.Equal(t, tt.expectedError, resp["message"])
			}
		})
	}
}`

	startIndex = strings.Index(strContent, "func TestMeUserNotFound(t *testing.T) {")
	if startIndex != -1 {
		endIndex := strings.Index(strContent[startIndex:], "}\n}\n")
		if endIndex != -1 {
			strContent = strContent[:startIndex] + testNotFound + "\n" + strContent[startIndex+endIndex+4:]
		}
	}


	// Replace TestMeResponseStructure
	testRespStr := `func TestMeResponseStructure(t *testing.T) {
	t.Run("response includes user profile", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/v1/auth/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		validUUIDStr := "550e8400-e29b-41d4-a716-446655440000"
		validUUID, _ := utils.StringToUUID(validUUIDStr)
		payload := &utils.TokenPayload{UserID: validUUIDStr}
		c.Set(auth.ClaimsKey, payload)

		mockQuerier := &MockQuerier{
			GetUserByIdFunc: func(ctx context.Context, id pgtype.UUID) (db.User, error) {
				return db.User{
					ID:            id,
					Name:          "Test User",
					Email:         pgtype.Text{String: "test@example.com", Valid: true},
					EmailVerified: true,
				}, nil
			},
		}

		service := NewService(mockQuerier, &config.Config{})
		handler := NewHandler(service)

		err := handler.Me(&c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp UserProfileResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.True(t, resp.Success)
		assert.Equal(t, "Test User", resp.Data.Name)
		assert.Equal(t, "test@example.com", resp.Data.Email)
		assert.True(t, resp.Data.EmailVerified)

		// Assert ID matches
		assert.Equal(t, validUUID.Bytes, resp.Data.ID.Bytes)
	})
}`

	startIndex = strings.Index(strContent, "func TestMeResponseStructure(t *testing.T) {")
	if startIndex != -1 {
		endIndex := strings.Index(strContent[startIndex:], "}\n}\n")
		if endIndex != -1 {
			strContent = strContent[:startIndex] + testRespStr + "\n" + strContent[startIndex+endIndex+4:]
		}
	}

	os.WriteFile(path, []byte(strContent), 0644)
}
