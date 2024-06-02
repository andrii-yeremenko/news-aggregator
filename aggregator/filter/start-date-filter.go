package filter

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/parser"
	"time"
)

// StartDateFilter is an aggregator.Filter that creates a subset from a given set of article.Article's
// corresponding to a given start date.
type StartDateFilter struct {
	startDate *time.Time
}

// NewStartDateFilter creates a new StartDateFilter instance with the given start date.
func NewStartDateFilter(startDateStr string) *StartDateFilter {

	dateParser := parser.NewDateParser()

	startDate, err := dateParser.ParseDefaultDateFormat(startDateStr)
	if err != nil {
		panic("invalid start date")
	}

	return &StartDateFilter{startDate: &startDate}
}

// Apply filters the article.Article's and returns a subset that meets predefined start date.
func (f *StartDateFilter) Apply(articles []article.Article) []article.Article {

	var filteredArticles []article.Article

	for _, selectedArticle := range articles {
		if f.matchStartDate(selectedArticle) {
			filteredArticles = append(filteredArticles, selectedArticle)
		}
	}

	return filteredArticles
}

func (f *StartDateFilter) matchStartDate(a article.Article) bool {

	articleDate := time.Time(a.Date())

	if f.startDate != nil && articleDate.Before(*f.startDate) {
		return false
	}

	return true
}
