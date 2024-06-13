package parser

import (
	"NewsAggregator/aggregator/model/resource"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestRSSParser_Parse(t *testing.T) {
	path := filepath.Join("testdata/rss", "test.xml")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	mockResource, err := resource.New("Test Source", "rss", resource.Content(content))
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	parser := &RSSParser{}

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

func TestRSSParser_Parse_InvalidFormat(t *testing.T) {

	path := filepath.Join("testdata/rss", "invalid_format_test.xml")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	mockResource, err := resource.New("Test Source", "rss", resource.Content(content))
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	parser := &RSSParser{}

	_, err = parser.Parse(*mockResource)

	assert.Errorf(t, err, "Parser should return an error when the file is empty")
}

func TestRSSParser_Parse_InvalidArticles(t *testing.T) {

	path := filepath.Join("testdata/rss", "invalid_data_test.xml")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	mockResource, err := resource.New("Test Source", "rss", resource.Content(content))
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	parser := &RSSParser{}

	_, err = parser.Parse(*mockResource)

	assert.Errorf(t, err, "Parser should return an error when articles are invalid")
}

func TestRSSParser_Parse_CorruptedDate(t *testing.T) {

	path := filepath.Join("testdata/rss", "corrupted_date_test.xml")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	mockResource, err := resource.New("Test Source", "rss", resource.Content(content))
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	parser := &RSSParser{}

	_, err = parser.Parse(*mockResource)

	assert.Errorf(t, err, "Parser should return an error when article creation date format is invalid or unknown")
}
