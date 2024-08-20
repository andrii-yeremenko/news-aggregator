package feed

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		source        Source
		format        Format
		link          Link
		expectedError error
	}{
		{
			name:          "valid feed",
			source:        "valid-source",
			format:        RSS,
			link:          "http://valid-link.com",
			expectedError: nil,
		},
		{
			name:          "empty source",
			source:        "",
			format:        RSS,
			link:          "http://valid-link.com",
			expectedError: errors.New("source cannot be empty"),
		},
		{
			name:          "unknown format",
			format:        0,
			source:        "valid-source",
			link:          "http://valid-link.com",
			expectedError: errors.New("format cannot be unknown"),
		},
		{
			name:          "empty link",
			source:        "valid-source",
			format:        RSS,
			link:          "",
			expectedError: errors.New("link cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feed, err := New(tt.source, tt.format, tt.link)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
				if feed.Source() != tt.source {
					t.Errorf("expected source: %v, got: %v", tt.source, feed.Source())
				}
				if feed.Format() != tt.format {
					t.Errorf("expected format: %v, got: %v", tt.format, feed.Format())
				}
				if feed.Link() != tt.link {
					t.Errorf("expected link: %v, got: %v", tt.link, feed.Link())
				}
			}
		})
	}
}

func TestFeedMethods(t *testing.T) {
	source := Source("test-source")
	link := Link("http://test-link.com")
	feed, err := New(source, RSS, link)
	if err != nil {
		t.Fatalf("unexpected error creating feed: %v", err)
	}

	t.Run("Source", func(t *testing.T) {
		if got := feed.Source(); got != source {
			t.Errorf("expected source: %v, got: %v", source, got)
		}
	})

	t.Run("Format", func(t *testing.T) {
		if got := feed.Format(); FormatToString(got) != FormatToString(RSS) {
			t.Errorf("expected format: %v, got: %v", RSS, got)
		}
	})

	t.Run("Link", func(t *testing.T) {
		if got := feed.Link(); got != link {
			t.Errorf("expected link: %v, got: %v", link, got)
		}
	})
}
