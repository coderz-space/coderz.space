package analytics

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"
)

// ParsePage parses the page query parameter with default value
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

// ParseLimit parses the limit query parameter with default and max values
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
