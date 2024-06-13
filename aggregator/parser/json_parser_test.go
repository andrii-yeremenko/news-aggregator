package parser

import (
	"NewsAggregator/aggregator/model/resource"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestJSONParser_Parse(t *testing.T) {
	path := filepath.Join("testdata/json", "test.json")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	mockResource, err := resource.New("Test Source", "json", resource.Content(content))
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	parser := &JSONParser{}

	articles, err := parser.Parse(*mockResource)
	assert.NoError(t, err, "Parser should not return an error")
	assert.NotEmpty(t, articles, "Parsed articles should not be empty")

	expectedTitle := "Test Title"
	expectedDescription := "Test Description"
	expectedAuthor := "John Doe"
	expectedLink := "http://example.com"
	expectedSource := "Test Source"

	assert.Equal(t, expectedTitle, string(articles[0].Title()), "Article title mismatch")
	assert.Equal(t, expectedDescription, string(articles[0].Description()), "Article description mismatch")
	assert.Equal(t, expectedAuthor, string(articles[0].Author()), "Article author mismatch")
	assert.Equal(t, expectedLink, string(articles[0].Link()), "Article link mismatch")
	assert.Equal(t, expectedSource, string(articles[0].Source()), "Article source mismatch")
}

func TestJSONParser_Parse_EmptyArticles(t *testing.T) {
	path := filepath.Join("testdata/json", "invalid_data_test.json")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	mockResource, err := resource.New("Test Source", "json", resource.Content(content))
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	parser := &JSONParser{}

	_, err = parser.Parse(*mockResource)

	assert.Errorf(t, err, "Parser should return an error when the articles are empty")
}

func TestJSONParser_Parse_InvalidDate(t *testing.T) {
	path := filepath.Join("testdata/json", "corrupted_date_test.json")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	mockResource, err := resource.New("Test Source", "json", resource.Content(content))
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	parser := &JSONParser{}

	_, err = parser.Parse(*mockResource)

	assert.Errorf(t, err, "Parser should return an error when the creation date is in invalid or unknown format")
}

func TestJSONParser_Parse_InvalidFormat(t *testing.T) {
	path := filepath.Join("testdata/json", "invalid_format_test.json")
	content, err := os.ReadFile(path)

	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	mockResource, err := resource.New("Test Source", "json", resource.Content(content))
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	parser := &JSONParser{}

	_, err = parser.Parse(*mockResource)

	assert.Errorf(t, err, "Parser should return an error when the file is in invalid or unknown json format")
}
