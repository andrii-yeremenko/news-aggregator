package parser

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/aggregator/model/resource"
	"fmt"
)

// parserProperties represents a composite key for the parser map.
type parserProperties struct {
	format    resource.Format
	publisher resource.Source
}

// Factory is a parser's selector according to the resource format and publisher.
type Factory struct {
	parsers map[parserProperties]aggregator.Parser
}

// NewParserFactory creates a new factory with predefined default parsers.
func NewParserFactory() *Factory {
	return &Factory{
		parsers: map[parserProperties]aggregator.Parser{
			{format: "json", publisher: "nbc-news"}:        &JSONParser{},
			{format: "rss", publisher: "abc-news"}:         &RSSParser{},
			{format: "rss", publisher: "washington-times"}: &RSSParser{},
			{format: "rss", publisher: "bbc-world"}:        &RSSParser{},
			{format: "html", publisher: "usa-today"}:       &USATodayHTMLParser{},
		},
	}
}

// RegisterParser registers a parser with a specific format and publisher.
func (f *Factory) RegisterParser(format resource.Format, publisher resource.Source, parser aggregator.Parser) {
	key := parserProperties{format: format, publisher: publisher}
	f.parsers[key] = parser
}

// GetParser returns a parser for the given resource.
func (f *Factory) GetParser(format resource.Format, publisher resource.Source) (aggregator.Parser, error) {
	key := parserProperties{format: format, publisher: publisher}
	parser, exists := f.parsers[key]
	if !exists {
		return nil, fmt.Errorf("no parser found for format: %s and publisher: %s", format, publisher)
	}
	return parser, nil
}
