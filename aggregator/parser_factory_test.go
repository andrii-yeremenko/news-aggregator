package aggregator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"news-aggregator/aggregator/model/article"
	"news-aggregator/aggregator/model/resource"
	"testing"
)

type mockParser struct{}

//goland:noinspection ALL
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
		{resource.JSON, "nbc-news", true},
		{resource.HTML, "usa-today", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d-%s", tt.format, tt.publisher), func(t *testing.T) {
			key := parserProperties{format: tt.format, publisher: tt.publisher}
			_, exists := factory.parsers[key]
			assert.Equal(t, tt.expected, exists)
		})
	}
}

func TestAddNewParser(t *testing.T) {
	factory := NewParserFactory()
	mockParser := &mockParser{}
	factory.AddNewParser(resource.RSS, "new-publisher", mockParser)

	key := parserProperties{format: resource.RSS, publisher: "new-publisher"}
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
		{resource.JSON, "nbc-news", true},
		{resource.RSS, "abc-news", true},
		{resource.RSS, "washington-times", true},
		{resource.RSS, "bbc-world", true},
		{resource.HTML, "usa-today", true},
		{resource.HTML, "non-existing", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d-%s", tt.format, tt.publisher), func(t *testing.T) {
			_, err := factory.GetParser(tt.format, tt.publisher)
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
