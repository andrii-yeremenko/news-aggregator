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

// NewEndDateFilter creates a new EndDateFilter instance with the given end date.
func NewEndDateFilter(endDateStr string) *EndDateFilter {

	dateParser := parser.NewDateParser()

	endDate, err := dateParser.ParseDefaultDateFormat(endDateStr)
	if err != nil {
		panic("invalid end date")
	}

	return &EndDateFilter{endDate: &endDate}
}

// Apply filters the data and returns a subset of articles.
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
