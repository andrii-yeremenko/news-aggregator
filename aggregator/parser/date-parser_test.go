package parser_test

import (
	"NewsAggregator/aggregator/parser"
	"testing"
	"time"
)

func TestNewDateParser(t *testing.T) {
	t.Run("Creating a new DateParser instance", func(t *testing.T) {
		dp := parser.NewDateParser()
		if dp == nil {
			t.Fatal("NewDateParser returned nil")
		}
	})
}

func TestDateParser_Parse(t *testing.T) {
	dp := parser.NewDateParser()

	tests := []struct {
		name     string
		input    string
		expected time.Time
		hasError bool
	}{
		{"Valid RFC1123Z format", "Thu, 28 May 2020 14:15:22 +0000", time.Date(2020, 5, 28, 14, 15, 22, 0, time.UTC), false},
		{"Valid RFC3339 format", "2020-05-28T14:15:22Z", time.Date(2020, 5, 28, 14, 15, 22, 0, time.UTC), false},
		{"Valid GMT format", "Mon, 02 Jan 2006 15:04:05 GMT", time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), false},
		{"Invalid date format", "invalid date", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := dp.Parse(tt.input)
			if (err != nil) != tt.hasError {
				t.Fatalf("expected error: %v, got: %v", tt.hasError, err)
			}
			if !tt.hasError && !parsed.Equal(tt.expected) {
				t.Fatalf("expected: %v, got: %v", tt.expected, parsed)
			}
		})
	}
}

func TestDateParser_ParseDefaultDateFormat(t *testing.T) {
	dp := parser.NewDateParser()

	tests := []struct {
		name     string
		input    string
		expected time.Time
		hasError bool
	}{
		{"Valid default format", "2024-02-01", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), false},
		{"Invalid date format", "invalid date", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := dp.ParseDefaultDateFormat(tt.input)
			if (err != nil) != tt.hasError {
				t.Fatalf("expected error: %v, got: %v", tt.hasError, err)
			}
			if !tt.hasError && !parsed.Equal(tt.expected) {
				t.Fatalf("expected: %v, got: %v", tt.expected, parsed)
			}
		})
	}
}
