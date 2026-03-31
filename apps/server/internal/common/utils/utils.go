package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// HashString returns a SHA256 hash of the input string
func HashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// UUIDToString converts a pgtype.UUID to its string representation
func UUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	buf := u.Bytes
	return hex.EncodeToString(buf[:4]) + "-" +
		hex.EncodeToString(buf[4:6]) + "-" +
		hex.EncodeToString(buf[6:8]) + "-" +
		hex.EncodeToString(buf[8:10]) + "-" +
		hex.EncodeToString(buf[10:16])
}

// StringToUUID converts a string to a pgtype.UUID
func StringToUUID(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	err := u.Scan(s)
	return u, err
}

// StringToInt converts a string to an integer
func StringToInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
