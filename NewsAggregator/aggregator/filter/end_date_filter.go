package filter

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/parser"
	"time"
)

// EndDateFilter is a Filter that creates a subset from a given set of article.Article's
// corresponding to a given end date.
type EndDateFilter struct {
	endDate *time.Time
}

// NewEndDateFilter creates a new EndDateFilter instance.
func NewEndDateFilter(endDateStr string) *EndDateFilter {

	dateParser := parser.NewDateParser()

	endDate, err := dateParser.ParseDefaultDateFormat(endDateStr)
	if err != nil {
		panic("invalid end date")
	}

	return &EndDateFilter{endDate: &endDate}
}

// Apply filters the data and returns a subset of articles.
func (filter *EndDateFilter) Apply(articles []article.Article) []article.Article {

	var filteredArticles []article.Article

	for _, selectedArticle := range articles {
		if filter.matchEndDate(selectedArticle) {
			filteredArticles = append(filteredArticles, selectedArticle)
		}
	}

	return filteredArticles
}

func (filter *EndDateFilter) matchEndDate(selectedArticle article.Article) bool {

	articleDate := time.Time(selectedArticle.Date())

	if filter.endDate != nil && articleDate.After(*filter.endDate) {
		return false
	}

	return true
}
