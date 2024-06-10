package aggregator

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"errors"
)

// Aggregator is a processor that collects specific information from resource.Resource
// and turns it into a collection of article.Article.
type Aggregator struct {
	articles      []article.Article
	parserFactory Factory
	filters       []Filter
}

// New creates a new Aggregator instance.
func New(factory Factory) (*Aggregator, error) {
	if factory == nil {
		err := errors.New("factory cannot be nil")
		return nil, err
	}

	return &Aggregator{
		articles:      []article.Article{},
		parserFactory: factory,
	}, nil
}

// AddFilter adds a filter to the aggregator.
func (agr *Aggregator) AddFilter(filter Filter) {
	agr.filters = append(agr.filters, filter)
}

// LoadResource loads articles from a resource and aggregates them.
func (agr *Aggregator) LoadResource(resource resource.Resource) error {
	newArticles, err := agr.aggregate(resource)
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

// aggregate fetches articles from a resource and parses them.
func (agr *Aggregator) aggregate(resource resource.Resource) ([]article.Article, error) {

	articlesParser, err := agr.parserFactory.GetParser(resource.Format(), resource.Source())
	if err != nil {
		return nil, err
	}

	return articlesParser.Parse(resource)
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
