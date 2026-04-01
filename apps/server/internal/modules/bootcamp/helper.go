package bootcamp

import (
	"time"
)

// ValidateDateRange checks if start_date is less than or equal to end_date
func ValidateDateRange(startDate, endDate string) bool {
	if startDate == "" || endDate == "" {
		return true // If either is empty, skip validation
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return false
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return false
	}

	return start.Before(end) || start.Equal(end)
}

// ParseDate converts a date string to pgtype.Date
func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", dateStr)
}

// FormatDate converts a time.Time to ISO date string
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}
