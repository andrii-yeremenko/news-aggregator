package filter

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/parser"
	"time"
)

// EndDateFilter is an aggregator.Filter that creates a subset from a given set of article.Article's
// corresponding to a given end date string.
type EndDateFilter struct {
	endDate *time.Time
}

// NewEndDateFilter creates a new EndDateFilter instance with the given end date.
func NewEndDateFilter(endDateStr string) (*EndDateFilter, error) {

	dateParser := parser.NewDateParser()

	endDate, err := dateParser.ParseDefaultDateFormat(endDateStr)
	if err != nil {
		return nil, err
	}

	return &EndDateFilter{endDate: &endDate}, nil
}

// Apply filters the article.Article's and returns a subset that meets predefined end date.
func (f *EndDateFilter) Apply(articles []article.Article) []article.Article {

	var filteredArticles []article.Article

	for _, selectedArticle := range articles {
		if f.matchEndDate(selectedArticle) {
			filteredArticles = append(filteredArticles, selectedArticle)
		}
	}

	return filteredArticles
}

func (f *EndDateFilter) matchEndDate(a article.Article) bool {

	articleDate := time.Time(a.Date())

	if f.endDate != nil && articleDate.After(*f.endDate) {
		return false
	}

	return true
}
