package aggregator

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"NewsAggregator/aggregator/parser"
	"strings"
	"time"
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
	processor, err := agr.parserFactory.GetParser(resource.Format(), resource.Publisher())
	if err != nil {
		return nil, err
	}

	return processor.Parse(resource.Content())
}

// FilterByKeywords filters articles by keywords.
func (agr *Aggregator) FilterByKeywords(keywords []string) []article.Article {
	var filteredArticles []article.Article
	for _, art := range agr.articles {
		for _, keyword := range keywords {
			keyword = strings.ToLower(keyword)
			title := strings.ToLower(string(art.Title()))
			description := strings.ToLower(string(art.Description()))
			if strings.Contains(title, keyword) || strings.Contains(description, keyword) {
				filteredArticles = append(filteredArticles, art)
				break
			}
		}
	}
	return filteredArticles
}

// FilterByDateRange filters articles within a specified date range.
func (agr *Aggregator) FilterByDateRange(startDateStr, endDateStr string) []article.Article {
	var filteredArticles []article.Article

	startDate, err := time.Parse("2006-02-01", startDateStr)
	if err != nil {
		return filteredArticles
	}
	endDate, err := time.Parse("2006-02-01", endDateStr)
	if err != nil {
		return filteredArticles
	}

	for _, art := range agr.articles {
		articleDate := time.Time(art.Date())

		if articleDate.After(startDate) && articleDate.Before(endDate) {
			filteredArticles = append(filteredArticles, art)
		}
	}
	return filteredArticles
}

// FilterBySources filters articles by specified sources.
func (agr *Aggregator) FilterBySources(sources []string) []article.Article {
	sourceSet := make(map[string]struct{})
	for _, source := range sources {
		sourceSet[source] = struct{}{}
	}
	var filteredArticles []article.Article
	for _, art := range agr.articles {
		if _, exists := sourceSet[string(art.Source())]; exists {
			filteredArticles = append(filteredArticles, art)
		}
	}
	return filteredArticles
}
