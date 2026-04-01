package app

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestHandleAppErrorMapsKnownErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "access denied", err: errors.New("ACCESS_DENIED"), wantStatus: http.StatusForbidden},
		{name: "not found", err: errors.New("QUESTION_NOT_FOUND"), wantStatus: http.StatusNotFound},
		{name: "conflict", err: errors.New("USERNAME_ALREADY_EXISTS"), wantStatus: http.StatusConflict},
		{name: "bad request", err: errors.New("INVALID_CURRENT_PASSWORD"), wantStatus: http.StatusBadRequest},
		{name: "internal", err: errors.New("SOMETHING_ELSE"), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/v1/app/context", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			if err := handleAppError(ctx, tt.err); err != nil {
				t.Fatalf("expected response to be written, got error %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}
