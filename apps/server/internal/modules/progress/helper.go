package progress

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// CursorData represents the structure of a pagination cursor
type CursorData struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

// EncodeCursor encodes cursor data to a base64 string
func EncodeCursor(id pgtype.UUID, createdAt time.Time) (string, error) {
	idStr := ""
	if id.Valid {
		// Convert UUID bytes to hex string
		buf := id.Bytes
		idStr = fmt.Sprintf("%x-%x-%x-%x-%x",
			buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16])
	}
	cursor := CursorData{
		ID:        idStr,
		CreatedAt: createdAt,
	}

	jsonData, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(jsonData), nil
}

// DecodeCursor decodes a base64 cursor string to cursor data
func DecodeCursor(cursorStr string) (*CursorData, error) {
	if cursorStr == "" {
		return nil, nil
	}

	jsonData, err := base64.StdEncoding.DecodeString(cursorStr)
	if err != nil {
		return nil, err
	}

	var cursor CursorData
	if err := json.Unmarshal(jsonData, &cursor); err != nil {
		return nil, err
	}

	return &cursor, nil
}

// ParseLimit parses and validates the limit query parameter
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

// ParseBoolParam parses a boolean query parameter
func ParseBoolParam(param string) *bool {
	if param == "" {
		return nil
	}

	val := param == "true" || param == "1"
	return &val
}

// FormatTimestamp formats a pgtype.Timestamptz to RFC3339 string
func FormatTimestamp(t pgtype.Timestamptz) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(time.RFC3339)
}

// FormatNullableString formats a pgtype.Text to string
func FormatNullableString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// IsValidUUID checks if a string is a valid UUID format
func IsValidUUID(s string) bool {
	var u pgtype.UUID
	err := u.Scan(s)
	return err == nil
}

// BuildFilters creates a filter map from query parameters
func BuildFilters(assignmentProblemID string, resolved *bool) map[string]string {
	filters := make(map[string]string)

	if assignmentProblemID != "" && IsValidUUID(assignmentProblemID) {
		filters["assignment_problem_id"] = assignmentProblemID
	}

	if resolved != nil {
		if *resolved {
			filters["resolved"] = "true"
		} else {
			filters["resolved"] = "false"
		}
	}

	return filters
}
