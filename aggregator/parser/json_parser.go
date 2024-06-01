package parser

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"encoding/json"
	"strings"
)

// JSONParser is a aggregator.Parser that parses JSON data.
type JSONParser struct{}

// Parse parses the JSON content into a list of articles.
func (p *JSONParser) Parse(resource resource.Resource) ([]article.Article, error) {
	var jsonResponse struct {
		Articles []struct {
			Source struct {
				Name string `json:"name"`
			} `json:"source"`
			Author      string `json:"author"`
			Title       string `json:"title"`
			Description string `json:"description"`
			PublishedAt string `json:"publishedAt"`
		} `json:"articles"`
	}

	byteContent := []byte(resource.Content())

	if err := json.Unmarshal(byteContent, &jsonResponse); err != nil {
		return nil, err
	}

	var articles []article.Article
	for _, a := range jsonResponse.Articles {
		publishedAt, err := NewDateParser().Parse(a.PublishedAt)
		if err != nil {
			return nil, err
		}

		art, err := article.NewArticleBuilder().
			SetTitle(article.Title(strings.TrimSpace(a.Title))).
			SetDescription(article.Description(strings.TrimSpace(a.Description))).
			SetDate(article.CreationDate(publishedAt)).
			SetSource(resource.Source()).
			SetAuthor(article.Author(strings.TrimSpace(a.Author))).
			Build()

		if err != nil {
			return nil, err
		}

		articles = append(articles, *art)
	}

	return articles, nil
}
