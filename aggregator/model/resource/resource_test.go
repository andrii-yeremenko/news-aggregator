package resource_test

import (
	"NewsAggregator/aggregator/model/resource"
	"errors"
	"testing"
)

func TestNewResource(t *testing.T) {

	tests := []struct {
		name        string
		source      resource.Source
		format      resource.Format
		content     resource.Content
		expectedErr error
	}{
		{
			name:        "Valid parameters",
			source:      "CNN",
			format:      "json",
			content:     "Some news content",
			expectedErr: nil,
		},
		{
			name:        "Empty source",
			source:      "",
			format:      "json",
			content:     "Some news content",
			expectedErr: errors.New("source cannot be empty"),
		},
		{
			name:        "Empty format",
			source:      "CNN",
			format:      "",
			content:     "Some news content",
			expectedErr: errors.New("format cannot be empty"),
		},
		{
			name:        "Empty content",
			source:      "CNN",
			format:      "json",
			content:     "",
			expectedErr: errors.New("content cannot be empty"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := resource.New(test.source, test.format, test.content)
			if err != nil {
				if err.Error() != test.expectedErr.Error() {
					t.Errorf("Expected error: %v, but got: %v", test.expectedErr, err)
				}
			} else {
				if test.expectedErr != nil {
					t.Errorf("Expected error: %v, but got nil", test.expectedErr)
				}
			}
		})
	}
}

func TestResourceMethods(t *testing.T) {

	r, err := resource.New("BBC", "xml", "Some news content")
	if err != nil {
		t.Fatalf("Error creating r: %v", err)
	}

	expectedSource := resource.Source("BBC")
	if source := r.Source(); source != expectedSource {
		t.Errorf("Expected source: %s, but got: %s", expectedSource, source)
	}

	expectedFormat := resource.Format("xml")
	if format := r.Format(); format != expectedFormat {
		t.Errorf("Expected format: %s, but got: %s", expectedFormat, format)
	}

	expectedContent := resource.Content("Some news content")
	if content := r.Content(); content != expectedContent {
		t.Errorf("Expected content: %s, but got: %s", expectedContent, content)
	}
}
