package aggregator

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
)

// Aggregator is responsible for aggregating news articles.
type Aggregator struct {
	articles      []article.Article
	parserFactory *Factory
	filters       []Filter
}

// New creates a new Aggregator instance.
func New(factory *Factory) *Aggregator {
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

// AddFilter adds a filter to the aggregator.
func (agr *Aggregator) AddFilter(filter Filter) {
	agr.filters = append(agr.filters, filter)
}

// GetFilteredArticles applies all filters to the articles and returns this filtered articles.
func (agr *Aggregator) GetFilteredArticles() []article.Article {

	if agr.filters == nil {
		return agr.articles
	}

	filteredArticles := agr.articles

	for _, selectedFilter := range agr.filters {
		filteredArticles = selectedFilter.Apply(filteredArticles)
	}

	return filteredArticles
}
