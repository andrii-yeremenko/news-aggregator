package parser

import (
	"NewsAggregator/aggregator/model/resource"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUSATodayHTMLParser_Parse(t *testing.T) {
	path := filepath.Join("testdata/usatodayhtml", "usa_today_test.html")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	mockResource, err := resource.New("Test Source", "html", resource.Content(content))
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	parser := &USATodayHTMLParser{}

	articles, err := parser.Parse(*mockResource)
	assert.NoError(t, err, "Parser should not return an error")
	assert.NotEmpty(t, articles, "Parsed articles should not be empty")

	assert.Equal(t, "Test title", string(articles[0].Title()), "Article title mismatch")
	assert.Equal(t, "Test description", string(articles[0].Description()), "Article description mismatch")
	assert.Equal(t, "https://usatoday.com/article_url", string(articles[0].Link()), "Article link mismatch")
	assert.Equal(t, "Test Source", string(articles[0].Source()), "Article source mismatch")
}

func TestUSATodayHTMLParser_Parse_CorruptedDate(t *testing.T) {
	path := filepath.Join("testdata/usatodayhtml", "usa_today_corrupted_date_test.html")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	mockResource, err := resource.New("Test Source", "html", resource.Content(content))
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	parser := &USATodayHTMLParser{}

	_, err = parser.Parse(*mockResource)
	assert.Errorf(t, err, "Parser should return an error when the date is invalid")
}
