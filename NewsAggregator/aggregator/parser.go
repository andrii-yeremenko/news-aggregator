package aggregator

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
)

// Parser is a component that takes input resource.Resource and converts it into a structured and unified
// article.Article format.
type Parser interface {
	Parse(content resource.Resource) ([]article.Article, error)
}
