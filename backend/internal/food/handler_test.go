package food

import (
	"net/http/httptest"
	"testing"
)

func TestParsePaginationDefaults(t *testing.T) {
	req := httptest.NewRequest("GET", "/food", nil)
	limit, offset, err := parsePagination(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if limit != defaultLimit {
		t.Fatalf("expected default limit %d, got %d", defaultLimit, limit)
	}
	if offset != 0 {
		t.Fatalf("expected default offset 0, got %d", offset)
	}
}

func TestParsePaginationClampAndNegative(t *testing.T) {
	req := httptest.NewRequest("GET", "/food?limit=1000&offset=-5", nil)
	limit, offset, err := parsePagination(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if limit != maxLimit {
		t.Fatalf("expected max limit %d, got %d", maxLimit, limit)
	}
	if offset != 0 {
		t.Fatalf("expected offset 0, got %d", offset)
	}
}

func TestParsePaginationInvalid(t *testing.T) {
	req := httptest.NewRequest("GET", "/food?limit=abc", nil)
	_, _, err := parsePagination(req)
	if err == nil {
		t.Fatalf("expected error for invalid limit")
	}
}

func TestNeedsEnrichment(t *testing.T) {
	category := "PANTRY"
	shelfLife := 10
	product := Product{Category: &category, ShelfLife: &shelfLife}
	if needsEnrichment(product) {
		t.Fatalf("expected no enrichment needed")
	}

	emptyCategory := ""
	product = Product{Category: &emptyCategory, ShelfLife: &shelfLife}
	if !needsEnrichment(product) {
		t.Fatalf("expected enrichment needed for empty category")
	}

	zeroShelf := 0
	product = Product{Category: &category, ShelfLife: &zeroShelf}
	if !needsEnrichment(product) {
		t.Fatalf("expected enrichment needed for zero shelf life")
	}
}

func TestCleanGeminiResponse(t *testing.T) {
	input := "```json\n[{\"food_name\":\"Banana\",\"category\":\"PRODUCE\",\"shelf_life\":5}]\n```"
	cleaned := cleanGeminiResponse(input)
	expected := "[{\"food_name\":\"Banana\",\"category\":\"PRODUCE\",\"shelf_life\":5}]"
	if cleaned != expected {
		t.Fatalf("expected cleaned response %q, got %q", expected, cleaned)
	}
}
