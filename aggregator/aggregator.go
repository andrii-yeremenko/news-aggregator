package aggregator

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"errors"
	"fmt"
)

// Aggregator is a processor that collects specific information from resource.Resource
// and turns it into a collection of article.Article.
type Aggregator struct {
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
		parserFactory: factory,
	}, nil
}

// AddFilter adds a filter to the aggregator.
func (agr *Aggregator) AddFilter(filter Filter) {
	agr.filters = append(agr.filters, filter)
}

// Aggregate fetches articles from a resource and parses them.
func (agr *Aggregator) Aggregate(resource resource.Resource) ([]article.Article, error) {

	articlesParser, err := agr.parserFactory.GetParser(resource.Format(), resource.Source())
	if err != nil {
		return nil, err
	}

	articles, err := articlesParser.Parse(resource)

	if err != nil {
		return nil, fmt.Errorf("failed to parse articles: %w", err)
	}

	if agr.filters != nil {
		return agr.getFilteredArticles(articles), nil
	}

	return articles, nil
}

// AggregateMultiple fetches articles from a multiple resources and parses them.
func (agr *Aggregator) AggregateMultiple(resources []resource.Resource) ([]article.Article, error) {

	var articles []article.Article

	for _, res := range resources {
		art, err := agr.Aggregate(res)
		if err != nil {
			return nil, fmt.Errorf("failed to aggregate articles: %w", err)
		}
		articles = append(articles, art...)
	}

	return articles, nil
}

// GetFilteredArticles applies all filters to the articles and returns this filtered articles.
func (agr *Aggregator) getFilteredArticles(parsedArticles []article.Article) []article.Article {

	filteredArticles := parsedArticles

	for _, selectedFilter := range agr.filters {
		filteredArticles = selectedFilter.Apply(filteredArticles)
	}

	return filteredArticles
}
