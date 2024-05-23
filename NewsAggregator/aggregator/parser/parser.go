package parser

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
)

// Parser is an interface for parsing news articles.
type Parser interface {
	Parse(content resource.Content) ([]article.Article, error)
}
