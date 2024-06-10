package aggregator_test

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockFactory struct{}

func (m *MockFactory) GetParser(format resource.Format, source resource.Source) (aggregator.Parser, error) {
	return &MockParser{}, nil
}

func (m *MockFactory) RegisterParser(format resource.Format, source resource.Source, parser aggregator.Parser) {

}

type MockParser struct{}

func (m *MockParser) Parse(res resource.Resource) ([]article.Article, error) {
	var articles []article.Article

	art, _ := article.NewArticleBuilder().SetTitle("Mock Title 1").
		SetDescription("Mock Description 1").
		SetDate(article.CreationDate(time.Now())).
		SetSource("Mock Source 1").
		SetAuthor("Mock Author 1").
		SetLink("Mock Link 1").
		Build()

	articles = append(articles, *art)

	return articles, nil
}

type MockFilter struct{}

func (m *MockFilter) Apply(articles []article.Article) []article.Article {
	var filtered []article.Article
	for _, a := range articles {
		if a.Title() == "Mock Title 1" {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

func TestAggregator(t *testing.T) {
	factory := &MockFactory{}
	agg, _ := aggregator.New(factory)

	t.Run("New Aggregator with nil factory", func(t *testing.T) {
		errorAgg, err := aggregator.New(nil)
		assert.Error(t, err)
		assert.Nil(t, errorAgg)
	})

	t.Run("New Aggregator", func(t *testing.T) {
		assert.NotNil(t, agg)
		assert.Equal(t, 0, len(agg.GetAllArticles()))
	})

	t.Run("Load Resource", func(t *testing.T) {
		res, err := resource.New("source1", "format1", "content1")
		assert.NoError(t, err)
		err = agg.LoadResource(*res)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(agg.GetAllArticles()))
	})

	t.Run("Get All Articles", func(t *testing.T) {
		articles := agg.GetAllArticles()
		assert.Equal(t, 1, len(articles))
		assert.Equal(t, "Mock Title 1", string(articles[0].Title()))
	})

	t.Run("Get Filtered Articles without given filter", func(t *testing.T) {
		filteredArticles := agg.GetFilteredArticles()
		assert.Equal(t, 1, len(filteredArticles))
		assert.Equal(t, "Mock Title 1", string(filteredArticles[0].Title()))
	})

	t.Run("Get Filtered Articles", func(t *testing.T) {
		agg.AddFilter(&MockFilter{})
		filteredArticles := agg.GetFilteredArticles()
		assert.Equal(t, 1, len(filteredArticles))
		assert.Equal(t, "Mock Title 1", string(filteredArticles[0].Title()))
	})

	t.Run("Load Resource with Error", func(t *testing.T) {
		factoryErr := &MockFactoryWithError{}
		aggErr, _ := aggregator.New(factoryErr)
		res, err := resource.New("source1", "format1", "content1")
		assert.NoError(t, err)
		err = aggErr.LoadResource(*res)
		assert.Error(t, err)
	})
}

type MockFactoryWithError struct{}

func (m *MockFactoryWithError) GetParser(format resource.Format, source resource.Source) (aggregator.Parser, error) {
	return nil, errors.New("parser not found")
}

func (m *MockFactoryWithError) RegisterParser(format resource.Format, source resource.Source, parser aggregator.Parser) {

}
