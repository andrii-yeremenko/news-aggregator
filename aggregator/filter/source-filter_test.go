package filter_test

import (
	"NewsAggregator/aggregator/filter"
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"testing"
	"time"
)

func TestSourceFilter_Apply(t *testing.T) {
	articles := []article.Article{
		createArticleWithSource("Source 1"),
		createArticleWithSource("Source 2"),
		createArticleWithSource("Source 3"),
	}

	tests := []struct {
		name     string
		sources  []string
		expected int
	}{
		{"Filter by Single Source", []string{"Source 1"}, 1},                  // Only one article should match
		{"Filter by Multiple Sources", []string{"Source 1", "Source 2"}, 2},   // Two articles should match
		{"Filter by Non-existent Source", []string{"Non-existent Source"}, 0}, // No articles should match
		{"Filter by non provided source", []string{}, 3},                      // All articles should match
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sourceFilter := filter.NewSourceFilter(test.sources)
			filteredArticles := sourceFilter.Apply(articles)
			if len(filteredArticles) != test.expected {
				t.Errorf("Expected %d articles, got %d", test.expected, len(filteredArticles))
			}
		})
	}
}

func createArticleWithSource(source string) article.Article {
	builder := article.NewArticleBuilder().
		SetTitle("Title").
		SetDescription("Description").
		SetSource(resource.Source(source)).
		SetDate(article.CreationDate(time.Now()))

	art, err := builder.Build()

	if err != nil {
		panic(err)
	}

	return *art
}
