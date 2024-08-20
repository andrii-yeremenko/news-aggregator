package feed

import (
	"fmt"
	"testing"
)

// TestFormatToString tests the FormatToString function.
func TestFormatToString(t *testing.T) {
	tests := []struct {
		format   Format
		expected string
	}{
		{RSS, "RSS"},
		{HTML, "HTML"},
		{JSON, "JSON"},
		{UNKNOWN, "UNKNOWN"},
		{Format(999), "UNKNOWN"}, // Test with an invalid format
	}

	for _, test := range tests {
		result := FormatToString(test.format)
		if result != test.expected {
			t.Errorf("FormatToString(%d) = %s; expected %s", test.format, result, test.expected)
		}
	}
}

// TestParseFormat tests the ParseFormat function.
func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
		err      error
	}{
		{"RSS", RSS, nil},
		{"rss", RSS, nil},
		{"HTML", HTML, nil},
		{"html", HTML, nil},
		{"JSON", JSON, nil},
		{"json", JSON, nil},
		{"UNKNOWN", UNKNOWN, fmt.Errorf("unknown format: UNKNOWN")},
		{"invalid", UNKNOWN, fmt.Errorf("unknown format: invalid")},
	}

	for _, test := range tests {
		result, err := ParseFormat(test.input)
		if result != test.expected || (err != nil && err.Error() != test.err.Error()) {
			t.Errorf("ParseFormat(%s) = %d, %v; expected %d, %v", test.input, result, err, test.expected, test.err)
		}
	}
}
