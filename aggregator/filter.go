package aggregator

import "NewsAggregator/aggregator/model/article"

// Filter is a tool that produces from a given set of data a specific subset of data that meets certain condition.
type Filter interface {

	// Apply filters the data and returns a subset of articles.
	Apply(articles []article.Article) []article.Article
}
