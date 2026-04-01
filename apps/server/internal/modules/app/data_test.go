package app

import "testing"

func TestListSheetsPreservesSupportedOrder(t *testing.T) {
	sheets := listSheets()

	if len(sheets) != len(orderedSheetKeys) {
		t.Fatalf("expected %d sheets, got %d", len(orderedSheetKeys), len(sheets))
	}

	for index, key := range orderedSheetKeys {
		if sheets[index].Key != key {
			t.Fatalf("expected sheet %d to be %q, got %q", index, key, sheets[index].Key)
		}
		if len(sheets[index].Questions) == 0 {
			t.Fatalf("expected sheet %q to expose questions", key)
		}
	}
}

func TestFindSheetQuestionByLinkRoundTrip(t *testing.T) {
	link := catalogLink("gfg-dsa-360", "gfg-1")

	question, ok := findSheetQuestionByLink(link)
	if !ok {
		t.Fatalf("expected link %q to resolve", link)
	}

	if question.Title != "Array Rotation" {
		t.Fatalf("expected resolved question title %q, got %q", "Array Rotation", question.Title)
	}
	if question.Topic != "Arrays" {
		t.Fatalf("expected resolved question topic %q, got %q", "Arrays", question.Topic)
	}
}

func TestFindSheetQuestionByLinkRejectsUnknownLinks(t *testing.T) {
	tests := []string{
		"",
		"https://example.com",
		"app-sheet:missing-parts",
		"app-sheet:unknown-sheet:gfg-1",
		"app-sheet:gfg-dsa-360:missing-question",
	}

	for _, link := range tests {
		t.Run(link, func(t *testing.T) {
			if _, ok := findSheetQuestionByLink(link); ok {
				t.Fatalf("expected link %q to be rejected", link)
			}
		})
	}
}
