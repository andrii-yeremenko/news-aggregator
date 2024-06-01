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
}

type rssChannel struct {
	Items []rssItem `xml:"channel>item"`
}

// Parse parses the RSS feed from the provided content and returns a list of articles.
func (p *RSSParser) Parse(resource resource.Resource) ([]article.Article, error) {
	var rssChannel rssChannel

	byteContent := []byte(string(resource.Content()))

	err := xml.Unmarshal(byteContent, &rssChannel)
	if err != nil {
		return nil, err
	}

	articles := make([]article.Article, 0, len(rssChannel.Items))

	for _, item := range rssChannel.Items {

		creationDate, err := NewDateParser().Parse(item.PubDate)

		if err != nil {
			return nil, err
		}

		builder := article.NewArticleBuilder().
			SetTitle(article.Title(strings.TrimSpace(item.Title))).
			SetDescription(article.Description(strings.TrimSpace(item.Description))).
			SetDate(article.CreationDate(creationDate)).
			SetSource(resource.Source())

		art, err := builder.Build()

		if err != nil {
			return nil, err
		}

		articles = append(articles, *art)
	}

	return articles, nil
}
