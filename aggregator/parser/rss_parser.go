package parser

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"encoding/xml"
	"strings"
)

// RSSParser is the aggregator.Parser for parsing RSS 2.0 feeds.
type RSSParser struct{}

type mediaThumbnail struct {
	URL    string `xml:"url,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type rssItem struct {
	Title       string           `xml:"title"`
	Link        string           `xml:"link"`
	PubDate     string           `xml:"pubDate"`
	Description string           `xml:"description"`
	MediaThumbs []mediaThumbnail `xml:"thumbnail"`
	Keywords    string           `xml:"keywords"`
	LinkedVideo string           `xml:"LinkedVideo"`
	Category    string           `xml:"category"`
	Creator     string           `xml:"creator" xml:"http://purl.org/dc/elements/1.1/creator"`
}

type rssChannel struct {
	Items []rssItem `xml:"channel>item"`
}

// Parse parses the RSS feed from the provided content and returns a list of articles.
func (p *RSSParser) Parse(resource resource.Resource) ([]article.Article, error) {
	byteContent := []byte(string(resource.Content()))

	rssChannel, err := p.unmarshalRSS(byteContent)
	if err != nil {
		return nil, err
	}

	articles, err := p.extractArticles(rssChannel, resource)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (p *RSSParser) unmarshalRSS(content []byte) (*rssChannel, error) {
	var channel rssChannel
	err := xml.Unmarshal(content, &channel)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (p *RSSParser) extractArticles(channel *rssChannel, resource resource.Resource) ([]article.Article, error) {
	var articles []article.Article

	for _, item := range channel.Items {
		art, err := p.parseArticle(item, resource)
		if err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}

	return articles, nil
}

func (p *RSSParser) parseArticle(item rssItem, resource resource.Resource) (article.Article, error) {
	creationDate, err := NewDateParser().Parse(item.PubDate)
	if err != nil {
		return article.Article{}, err
	}

	builder := article.NewArticleBuilder().
		SetTitle(article.Title(strings.TrimSpace(item.Title))).
		SetDescription(article.Description(strings.TrimSpace(item.Description))).
		SetDate(article.CreationDate(creationDate)).
		SetSource(resource.Source()).
		SetAuthor(article.Author(strings.TrimSpace(item.Creator))).
		SetLink(article.Link(strings.TrimSpace(item.Link)))

	newArticle, err := builder.Build()
	if err != nil {
		return article.Article{}, err
	}

	return *newArticle, nil
}
