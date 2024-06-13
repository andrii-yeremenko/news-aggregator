package aggregator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"news-aggregator/aggregator/model/article"
	"news-aggregator/aggregator/model/resource"
	"testing"
)

type mockParser struct{}

func (m *mockParser) Parse(resource resource.Resource) ([]article.Article, error) {
	return nil, nil
}

func TestNewParserFactory(t *testing.T) {
	factory := NewParserFactory()

	tests := []struct {
		format    resource.Format
		publisher resource.Source
		expected  bool
	}{
		{"json", "nbc-news", true},
		{"rss", "abc-news", true},
		{"rss", "washington-times", true},
		{"rss", "bbc-world", true},
		{"html", "usa-today", true},
		{"xml", "non-existing", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s-%s", tt.format, tt.publisher), func(t *testing.T) {
			key := parserProperties{format: tt.format, publisher: tt.publisher}
			_, exists := factory.parsers[key]
			assert.Equal(t, tt.expected, exists)
		})
	}
}

func TestAddNewParser(t *testing.T) {
	factory := NewParserFactory()
	mockParser := &mockParser{}
	factory.AddNewParser("xml", "new-publisher", mockParser)

	key := parserProperties{format: "xml", publisher: "new-publisher"}
	p, exists := factory.parsers[key]

	assert.True(t, exists)
	assert.Equal(t, mockParser, p)
}

func TestGetParser(t *testing.T) {
	factory := NewParserFactory()

	tests := []struct {
		format    resource.Format
		publisher resource.Source
		expected  bool
	}{
		{"json", "nbc-news", true},
		{"rss", "abc-news", true},
		{"rss", "washington-times", true},
		{"rss", "bbc-world", true},
		{"html", "usa-today", true},
		{"xml", "non-existing", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s-%s", tt.format, tt.publisher), func(t *testing.T) {
			_, err := factory.GetParser(tt.format, tt.publisher)
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
