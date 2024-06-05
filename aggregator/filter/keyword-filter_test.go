package filter_test

import (
	"NewsAggregator/aggregator/filter"
	"NewsAggregator/aggregator/model/article"
	"testing"
	"time"
)

func TestKeywordFilter_Apply(t *testing.T) {
	articles := []article.Article{
		createArticleWithKeywords("Title Ukraine", "Description Wonderful"),
		createArticleWithKeywords("Title Southland", "Description Wonderful"),
		createArticleWithKeywords("Title Kharkiv", "Description best city in Ukraine"),
	}

	tests := []struct {
		name     string
		keywords []string
		expected int
	}{
		{"Filter by Single Keyword",
			[]string{"russia"}, 0}, // No articles should match
		{"Filter by Title Keyword",
			[]string{"ukraine"}, 2}, // Two articles should match
		{"Filter by Description Keyword",
			[]string{"wonderful"}, 2}, // Two articles should match
		{"Filter by Description Keyword",
			[]string{"description"}, 3}, // All articles should match
		{"Filter by Non-existent Keyword",
			[]string{"best"}, 1}, // One articles should match
		{"Filter by non provided keyword",
			[]string{}, 3}, // All articles should match
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			keywordFilter := filter.NewKeywordFilter(test.keywords)
			filteredArticles := keywordFilter.Apply(articles)
			if len(filteredArticles) != test.expected {
				t.Errorf("Expected %d articles, got %d", test.expected, len(filteredArticles))
			}
		})
	}
}

func createArticleWithKeywords(title, description string) article.Article {
	builder := article.NewArticleBuilder().
		SetTitle(article.Title(title)).
		SetDescription(article.Description(description)).
		SetSource("Source").
		SetDate(article.CreationDate(time.Now()))

	art, err := builder.Build()

	if err != nil {
		panic(err)
	}

	return *art
}
