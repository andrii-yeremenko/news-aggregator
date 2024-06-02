package aggregator

import (
	"NewsAggregator/aggregator/model/resource"
	"NewsAggregator/aggregator/parser"
	"fmt"
)

// parserProperties represents a composite key for the parser map.
type parserProperties struct {
	format    resource.Format
	publisher resource.Source
}

// Factory is a parser's selector, according to the resource format and publisher.
type Factory struct {
	parsers map[parserProperties]Parser
}

// NewParserFactory creates a new factory with predefined default parsers.
func NewParserFactory() *Factory {
	return &Factory{
		parsers: map[parserProperties]Parser{
			{format: "json", publisher: "nbc-news"}:        &parser.JSONParser{},
			{format: "rss", publisher: "abc-news"}:         &parser.RSSParser{},
			{format: "rss", publisher: "washington-times"}: &parser.RSSParser{},
			{format: "rss", publisher: "bbc-world"}:        &parser.RSSParser{},
			{format: "html", publisher: "usa-today"}:       &parser.USATodayHTMLParser{},
		},
	}
}

// RegisterParser registers a parser with a specific format and publisher.
func (f *Factory) RegisterParser(format resource.Format, publisher resource.Source, parser Parser) {
	key := parserProperties{format: format, publisher: publisher}
	f.parsers[key] = parser
}

// GetParser returns a parser for the given resource.
func (f *Factory) GetParser(format resource.Format, publisher resource.Source) (Parser, error) {
	key := parserProperties{format: format, publisher: publisher}
	p, exists := f.parsers[key]
	if !exists {
		return nil, fmt.Errorf("no parser found for format: %s and publisher: %s", format, publisher)
	}
	return p, nil
}
