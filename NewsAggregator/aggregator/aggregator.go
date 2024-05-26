package aggregator

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"NewsAggregator/aggregator/parser"
)

// Aggregator is responsible for aggregating news articles.
type Aggregator struct {
	articles      []article.Article
	parserFactory *parser.Factory
}

// NewAggregator creates a new Aggregator instance.
func NewAggregator(factory *parser.Factory) *Aggregator {
	if factory == nil {
		panic("factory is nil")
	}

	return &Aggregator{
		articles:      []article.Article{},
		parserFactory: factory,
	}
}

// LoadResource loads articles from a resource and aggregates them.
func (agr *Aggregator) LoadResource(resource resource.Resource) error {
	newArticles, err := agr.Aggregate(resource)
	if err != nil {
		return err
	}
	agr.articles = append(agr.articles, newArticles...)
	return nil
}

// GetAllArticles returns all articles from the aggregator.
func (agr *Aggregator) GetAllArticles() []article.Article {
	return agr.articles
}

// Aggregate fetches articles from a resource and parses them.
func (agr *Aggregator) Aggregate(resource resource.Resource) ([]article.Article, error) {

	articlesParser, err := agr.parserFactory.GetParser(resource.Format(), resource.Source())
	if err != nil {
		return nil, err
	}

	return articlesParser.Parse(resource)
}

// ApplyFilter filters the articles based on the provided filter.
func (agr *Aggregator) ApplyFilter(filter Filter) []article.Article {
	return filter.Apply(agr.articles)
}
