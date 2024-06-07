package parser

import (
	"NewsAggregator/aggregator/model/resource"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestRSSParser_Parse(t *testing.T) {
	path := filepath.Join("testdata", "test.xml")
	content, err := ioutil.ReadFile(path)
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
