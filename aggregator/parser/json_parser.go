package parser

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"encoding/json"
	"strings"
)

// JSONParser is an aggregator.Parser that parses JSON data.
type JSONParser struct{}

type jsonArticle struct {
	Source struct {
		Name string `json:"name"`
	} `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PublishedAt string `json:"publishedAt"`
}

type jsonResponse struct {
	Articles []jsonArticle `json:"articles"`
}

// Parse parses the JSON content into a list of articles.
func (p *JSONParser) Parse(resource resource.Resource) ([]article.Article, error) {
	byteContent := []byte(resource.Content())

	response, err := p.unmarshalJSON(byteContent)
	if err != nil {
		return nil, err
	}

	articles, err := p.extractArticles(response, resource)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (p *JSONParser) unmarshalJSON(content []byte) (*jsonResponse, error) {
	var response jsonResponse
	err := json.Unmarshal(content, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (p *JSONParser) extractArticles(response *jsonResponse, resource resource.Resource) ([]article.Article, error) {
	var articles []article.Article

	for _, a := range response.Articles {
		art, err := p.parseArticle(a, resource)
		if err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}

	return articles, nil
}

func (p *JSONParser) parseArticle(a jsonArticle, resource resource.Resource) (article.Article, error) {
	publishedAt, err := NewDateParser().Parse(a.PublishedAt)
	if err != nil {
		return article.Article{}, err
	}

	builder := article.NewArticleBuilder().
		SetTitle(article.Title(strings.TrimSpace(a.Title))).
		SetDescription(article.Description(strings.TrimSpace(a.Description))).
		SetDate(article.CreationDate(publishedAt)).
		SetSource(resource.Source()).
		SetAuthor(article.Author(strings.TrimSpace(a.Author)))

	newArticle, err := builder.Build()
	if err != nil {
		return article.Article{}, err
	}

	return *newArticle, nil
}
