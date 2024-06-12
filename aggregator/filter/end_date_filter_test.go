package filter_test

import (
	"NewsAggregator/aggregator/filter"
	"NewsAggregator/aggregator/model/article"
	"testing"
	"time"
)

func TestEndDateFilter_Apply(t *testing.T) {
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
		{"Filter articles before end date", "2024-26-06", 3},
		{"Filter articles before end date", "2024-21-06", 2},
		{"Filter articles that are after end date", "2024-11-06", 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dateFilter, _ := filter.NewEndDateFilter(test.endDate)
			filteredArticles := dateFilter.Apply(articles)
			if len(filteredArticles) != test.expected {
				t.Errorf("Expected %d articles, got %d", test.expected, len(filteredArticles))
			}
		})
	}
}

func TestEndDateFilter_Error(t *testing.T) {
	_, err := filter.NewEndDateFilter("invalid date!")

	if err == nil {
		t.Errorf("End date filter should return an error for an invalid date")
	}
}

func createArticleWithDate(dateStr string) article.Article {
	date, _ := time.Parse("2006-02-01", dateStr)
	builder := article.NewArticleBuilder().
		SetTitle("Title").
		SetDescription("Description").
		SetDate(article.CreationDate(date)).
		SetSource("Sample Source")

	art, err := builder.Build()

	if err != nil {
		panic(err)
	}

	return *art
}
