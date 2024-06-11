package filter_test

import (
	"NewsAggregator/aggregator/filter"
	"NewsAggregator/aggregator/model/article"
	"testing"
)

func TestStartDateFilter_Apply(t *testing.T) {
	articles := []article.Article{
		createArticleWithDate("2024-15-06"),
		createArticleWithDate("2024-20-06"),
		createArticleWithDate("2024-25-06"),
	}

	tests := []struct {
		name     string
		endDate  string
		expected int
	}{
		{"Filter articles after start date", "2024-10-06", 3},
		{"Filter a few articles after start date", "2024-16-06", 2},
		{"Filter articles that are before start date", "2024-27-06", 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dateFilter, _ := filter.NewStartDateFilter(test.endDate)
			filteredArticles := dateFilter.Apply(articles)
			if len(filteredArticles) != test.expected {
				t.Errorf("Expected %d articles, got %d", test.expected, len(filteredArticles))
			}
		})
	}
}

func TestStartDateFilter_Error(t *testing.T) {
	_, err := filter.NewStartDateFilter("invalid date!")

	if err == nil {
		t.Errorf("Start date filter should return an error for an invalid date")
	}
}
