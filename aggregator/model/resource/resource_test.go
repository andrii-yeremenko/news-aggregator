package resource_test

import (
	"errors"
	"news-aggregator/aggregator/model/resource"
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
			format:      resource.JSON,
			content:     "Some news content",
			expectedErr: nil,
		},
		{
			name:        "Empty source",
			source:      "",
			format:      resource.JSON,
			content:     "Some news content",
			expectedErr: errors.New("source cannot be empty"),
		},
		{
			name:        "Empty content",
			source:      "CNN",
			format:      resource.JSON,
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

	r, err := resource.New("BBC", resource.RSS, "Some news content")
	if err != nil {
		t.Fatalf("Error creating r: %v", err)
	}

	expectedSource := resource.Source("BBC")
	if source := r.Source(); source != expectedSource {
		t.Errorf("Expected source: %s, but got: %s", expectedSource, source)
	}

	expectedFormat := resource.Format(resource.RSS)
	if format := r.Format(); format != expectedFormat {
		t.Errorf("Expected format: %d, but got: %d", expectedFormat, format)
	}

	expectedContent := resource.Content("Some news content")
	if content := r.Content(); content != expectedContent {
		t.Errorf("Expected content: %s, but got: %s", expectedContent, content)
	}
}
