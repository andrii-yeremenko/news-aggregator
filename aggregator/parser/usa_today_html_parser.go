package parser

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

// USATodayHTMLParser is an aggregator.Parser for HTML pages from USA Today.
type USATodayHTMLParser struct{}

// Parse extracts articles from the provided HTML resource.
func (p *USATodayHTMLParser) Parse(resource resource.Resource) ([]article.Article, error) {
	content := string(resource.Content())
	doc, err := p.createDocumentFromContent(content)
	if err != nil {
		return nil, err
	}

	articles, err := p.extractArticles(doc, resource)
	if err != nil {
		return nil, err
	}

	if len(articles) == 0 {
		return nil, errors.New("no articles found")
	}

	return articles, nil
}

func (p *USATodayHTMLParser) createDocumentFromContent(content string) (*goquery.Document, error) {
	reader := strings.NewReader(content)
	return goquery.NewDocumentFromReader(reader)
}

func (p *USATodayHTMLParser) extractArticles(doc *goquery.Document, resource resource.Resource) ([]article.Article, error) {
	var articles []article.Article

	doc.Find("main.gnt_cw div.gnt_m_flm a.gnt_m_flm_a").Each(func(i int, s *goquery.Selection) {
		art, err := p.parseArticle(s, resource)
		if err == nil {
			articles = append(articles, art)
		}
	})

	return articles, nil
}

func (p *USATodayHTMLParser) parseArticle(s *goquery.Selection, resource resource.Resource) (article.Article, error) {
	title := strings.TrimSpace(s.Text())
	description := p.getAttribute(s, "data-c-br")
	dateAttr := p.getChildAttribute(s, "div.gnt_m_flm_sbt", "data-c-dt")
	creationDate, err := NewDateParser().Parse(dateAttr)
	if err != nil {
		return article.Article{}, err
	}

	builder := article.NewArticleBuilder().
		SetTitle(article.Title(title)).
		SetDescription(article.Description(description)).
		SetDate(article.CreationDate(creationDate)).
		SetSource(resource.Source())

	newArticle, err := builder.Build()
	if err != nil {
		return article.Article{}, err
	}

	return *newArticle, nil
}

func (p *USATodayHTMLParser) getAttribute(s *goquery.Selection, attrName string) string {
	attrValue, _ := s.Attr(attrName)
	return strings.TrimSpace(attrValue)
}

func (p *USATodayHTMLParser) getChildAttribute(s *goquery.Selection, childSelector string, attrName string) string {
	child := s.Find(childSelector)
	return p.getAttribute(child, attrName)
}
