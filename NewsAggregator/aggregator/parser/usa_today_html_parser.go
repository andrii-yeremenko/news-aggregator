package parser

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

// USATodayHTMLParser is a parser for HTML pages from USA Today.
type USATodayHTMLParser struct{}

func (p *USATodayHTMLParser) Parse(resource resource.Resource) ([]article.Article, error) {
	var articles []article.Article

	content := string(resource.Content())
	reader := strings.NewReader(content)

	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return nil, err
	}

	doc.Find("main.gnt_cw div.gnt_m_flm a.gnt_m_flm_a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		description, _ := s.Attr("data-c-br")
		dateAttr, _ := s.Find("div.gnt_m_flm_sbt").Attr("data-c-dt")
		creationDate, err := NewDateParser().Parse(dateAttr)
		if err != nil {
			return
		}

		builder := article.NewArticleBuilder().
			SetTitle(article.Title(strings.TrimSpace(title))).
			SetDescription(article.Description(strings.TrimSpace(description))).
			SetDate(article.CreationDate(creationDate)).
			SetSource(resource.Source())

		newArticle, err := builder.Build()
		if err != nil {
			return
		}

		articles = append(articles, *newArticle)
	})

	if len(articles) == 0 {
		return nil, errors.New("no articles found")
	}

	return articles, nil
}
