package utils

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"
)

// PaginationConfig holds pagination configuration
type PaginationConfig struct {
	DefaultLimit int
	MaxLimit     int
}

// DefaultPaginationConfig returns default pagination settings
func DefaultPaginationConfig() PaginationConfig {
	return PaginationConfig{
		DefaultLimit: 20,
		MaxLimit:     100,
	}
}

// OffsetPaginationMeta represents offset-based pagination metadata
type OffsetPaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// CursorPaginationMeta represents cursor-based pagination metadata
type CursorPaginationMeta struct {
	NextCursor string `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore"`
	Limit      int    `json:"limit"`
}

// ParsePage parses and validates page number from string
func ParsePage(pageStr string) int {
	if pageStr == "" {
		return 1
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

// ParseLimit parses and validates limit from string with default and max values
func ParseLimit(limitStr string, defaultLimit, maxLimit int) int {
	if limitStr == "" {
		return defaultLimit
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

// CalculateOffset calculates the offset for offset-based pagination
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// NewOffsetPagination creates offset pagination metadata
func NewOffsetPagination(page, limit, total int) *OffsetPaginationMeta {
	return &OffsetPaginationMeta{
		Page:  page,
		Limit: limit,
		Total: total,
	}
}

// CursorData represents cursor pagination data
type CursorData struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

// EncodeCursor encodes cursor data to base64 string
func EncodeCursor(id string, timestamp time.Time) (string, error) {
	data := CursorData{
		ID:        id,
		Timestamp: timestamp,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonData), nil
}

// DecodeCursor decodes base64 cursor string to cursor data
func DecodeCursor(cursor string) (*CursorData, error) {
	if cursor == "" {
		return nil, nil
	}
	jsonData, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}
	var data CursorData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// NewCursorPagination creates cursor pagination metadata
func NewCursorPagination(hasMore bool, limit int, nextCursor string) *CursorPaginationMeta {
	return &CursorPaginationMeta{
		HasMore:    hasMore,
		Limit:      limit,
		NextCursor: nextCursor,
	}
}

// ValidatePaginationParams validates pagination parameters
func ValidatePaginationParams(page, limit int, config PaginationConfig) error {
	if page < 1 {
		return ErrInvalidPage
	}
	if limit < 1 || limit > config.MaxLimit {
		return ErrInvalidLimit
	}
	return nil
}

// Pagination errors
var (
	ErrInvalidPage   = NewValidationError("page must be greater than 0")
	ErrInvalidLimit  = NewValidationError("limit must be between 1 and max limit")
	ErrInvalidCursor = NewValidationError("invalid cursor format")
)

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}
