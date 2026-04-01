package utils

import (
	"testing"
)

func TestGenerateTokenSignature(t *testing.T) {
	payload := TokenPayload{
		UserID:   "123",
		Email:    "test@example.com",
		Role:     "user",
		UserName: "testuser",
	}
	_, err := GenerateToken(payload, "1h", "secret")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
}
