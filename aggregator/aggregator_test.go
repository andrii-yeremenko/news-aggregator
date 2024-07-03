package aggregator_test

import (
	"news-aggregator/aggregator"
	"news-aggregator/aggregator/model/article"
	"news-aggregator/aggregator/model/resource"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockFilter struct{}

func (m *MockFilter) Apply(articles []article.Article) []article.Article {
	var filtered []article.Article
	for _, a := range articles {
		if a.Source() == "source2" {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

type MockFactory struct{}

func (m *MockFactory) GetParser(format resource.Format, source resource.Source) (aggregator.Parser, error) {
	if format == 0 || source == "invalid" {
		return nil, assert.AnError
	}
	return &MockParser{}, nil
}

//goland:noinspection GoUnusedParameter,GoUnusedParameter,GoUnusedParameter
func (m *MockFactory) AddNewParser(format resource.Format, source resource.Source, parser aggregator.Parser) {

}

type MockParser struct{}

func (m *MockParser) Parse(res resource.Resource) ([]article.Article, error) {
	var articles []article.Article

	art, _ := article.NewArticleBuilder().SetTitle("Mock Title 1").
		SetDescription("Mock Description 1").
		SetDate(article.CreationDate(time.Now())).
		SetSource(res.Source()).
		SetAuthor("Mock Author 1").
		SetLink("Mock Link 1").
		Build()

	articles = append(articles, *art)

	return articles, nil
}

func TestAggregator(t *testing.T) {
	factory := &MockFactory{}
	agg, _ := aggregator.New(factory)

	t.Run("New Aggregator with nil factory", func(t *testing.T) {
		errorAgg, err := aggregator.New(nil)
		assert.Error(t, err)
		assert.Nil(t, errorAgg)
	})

	t.Run("Aggregate single resource without applied filters", func(t *testing.T) {
		res, err := resource.New("source1", resource.JSON, "content1")
		assert.NoError(t, err)
		articles, err := agg.Aggregate(*res)
		assert.Equal(t, 1, len(articles))
		assert.Equal(t, "Mock Title 1", string(articles[0].Title()))
		assert.Equal(t, "Mock Description 1", string(articles[0].Description()))
	})

	t.Run("Aggregate multiple resources without applied filters", func(t *testing.T) {
		res1, err := resource.New("source1", resource.JSON, "content1")
		assert.NoError(t, err)
		res2, err := resource.New("source2", resource.JSON, "content1")
		assert.NoError(t, err)
		articles, err := agg.AggregateMultiple([]resource.Resource{*res1, *res2})
		assert.Equal(t, 2, len(articles))
		assert.Equal(t, "source1", string(articles[0].Source()))
		assert.Equal(t, "source2", string(articles[1].Source()))
	})

	t.Run("Aggregate multiple resources with applied filters", func(t *testing.T) {
		res1, err := resource.New("source1", resource.JSON, "content1")
		assert.NoError(t, err)
		res2, err := resource.New("source2", resource.JSON, "content1")
		assert.NoError(t, err)
		agg.AddFilter(&MockFilter{})
		articles, err := agg.AggregateMultiple([]resource.Resource{*res1, *res2})
		assert.Equal(t, 1, len(articles))
		assert.Equal(t, "source2", string(articles[0].Source()))
	})

	t.Run("Aggregate incorrect resource", func(t *testing.T) {
		res, err := resource.New("invalid", resource.JSON, "invalid")
		assert.NoError(t, err)
		_, err = agg.Aggregate(*res)
		assert.Error(t, err)
	})
}
