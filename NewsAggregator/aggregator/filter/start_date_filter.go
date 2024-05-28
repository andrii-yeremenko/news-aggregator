package filter

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/parser"
	"time"
)

// StartDateFilter is a Filter that creates a subset from a given set of article.Article's
// corresponding to a given start date.
type StartDateFilter struct {
	startDate *time.Time
}

// NewStartDateFilter creates a new StartDateFilter instance.
func NewStartDateFilter(startDateStr string) *StartDateFilter {

	dateParser := parser.NewDateParser()

	startDate, err := dateParser.ParseDefaultDateFormat(startDateStr)
	if err != nil {
		panic("invalid start date")
	}

	return &StartDateFilter{startDate: &startDate}
}

// Apply filters the data and returns a subset of articles.
func (filter *StartDateFilter) Apply(articles []article.Article) []article.Article {

	var filteredArticles []article.Article

	for _, selectedArticle := range articles {
		if filter.matchStartDate(selectedArticle) {
			filteredArticles = append(filteredArticles, selectedArticle)
		}
	}

	return filteredArticles
}

func (filter *StartDateFilter) matchStartDate(art article.Article) bool {

	articleDate := time.Time(art.Date())

	if filter.startDate != nil && articleDate.Before(*filter.startDate) {
		return false
	}

	return true
}
