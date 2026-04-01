package app

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestNormalizeUsername(t *testing.T) {
	if got := normalizeUsername("  Alice_User  "); got != "alice_user" {
		t.Fatalf("expected normalized username %q, got %q", "alice_user", got)
	}
}

func TestValidatePasswordComplexity(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{password: "Password123", valid: true},
		{password: "lettersonly", valid: false},
		{password: "123456789", valid: false},
		{password: "Alpha9", valid: true},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			if got := validatePasswordComplexity(tt.password); got != tt.valid {
				t.Fatalf("expected complexity check for %q to be %v, got %v", tt.password, tt.valid, got)
			}
		})
	}
}

func TestQuestionRowToQuestionDataUsesCatalogMetadata(t *testing.T) {
	assignedAt := time.Date(2026, time.April, 1, 9, 0, 0, 0, time.UTC)
	completedAt := time.Date(2026, time.April, 2, 9, 0, 0, 0, time.UTC)

	row := questionRow{
		ID:           "assignment-problem-1",
		Title:        "Array Rotation",
		Description:  "database description",
		Difficulty:   "easy",
		ExternalLink: catalogLink("gfg-dsa-360", "gfg-1"),
		AppProgress:  "completed",
		Notes:        "notes",
		Resources:    "resources",
		AssignedAt:   assignedAt,
		CompletedAt: pgtype.Timestamptz{
			Time:  completedAt,
			Valid: true,
		},
	}

	data := row.toQuestionData()

	if data.Description != "Practice array rotation techniques and in-place updates." {
		t.Fatalf("expected catalog description, got %q", data.Description)
	}
	if data.Topic != "Arrays" {
		t.Fatalf("expected catalog topic %q, got %q", "Arrays", data.Topic)
	}
	if data.Status != "completed" {
		t.Fatalf("expected completed status, got %q", data.Status)
	}
	if data.CompletedAt != completedAt.Format(time.RFC3339) {
		t.Fatalf("expected completedAt %q, got %q", completedAt.Format(time.RFC3339), data.CompletedAt)
	}
}

func TestQuestionRowToQuestionDataFallsBackToDatabaseFields(t *testing.T) {
	assignedAt := time.Date(2026, time.April, 1, 9, 0, 0, 0, time.UTC)
	row := questionRow{
		ID:           "assignment-problem-2",
		Title:        "Custom Problem",
		Description:  "database description",
		Difficulty:   "medium",
		AppProgress:  "",
		LegacyStatus: "attempted",
		AssignedAt:   assignedAt,
	}

	data := row.toQuestionData()

	if data.Description != "database description" {
		t.Fatalf("expected database description fallback, got %q", data.Description)
	}
	if data.Topic != "General" {
		t.Fatalf("expected default topic %q, got %q", "General", data.Topic)
	}
	if data.ProgressStatus != "revision_needed" {
		t.Fatalf("expected attempted legacy status to map to revision_needed, got %q", data.ProgressStatus)
	}
	if data.Status != "pending" {
		t.Fatalf("expected non-completed question to stay pending, got %q", data.Status)
	}
}

func TestMapProgressToLegacyStatus(t *testing.T) {
	tests := []struct {
		progress string
		expected string
	}{
		{progress: "completed", expected: "completed"},
		{progress: "discussion_needed", expected: "attempted"},
		{progress: "revision_needed", expected: "attempted"},
		{progress: "not_started", expected: "pending"},
		{progress: "unexpected", expected: "pending"},
	}

	for _, tt := range tests {
		t.Run(tt.progress, func(t *testing.T) {
			if got := mapProgressToLegacyStatus(tt.progress); got != tt.expected {
				t.Fatalf("expected legacy status %q, got %q", tt.expected, got)
			}
		})
	}
}
