package aggregator

import "NewsAggregator/aggregator/model/article"

// Filter is a construct that processes a collection of article.Article's and returns
// a new collection containing only those article.Article's that meet specific criteria.
type Filter interface {

	// Apply filters the data and returns a subset of articles.
	Apply(articles []article.Article) []article.Article
}
