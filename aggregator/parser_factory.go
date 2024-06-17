package aggregator

import (
	"fmt"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/aggregator/parser"
)

type Factory interface {
	AddNewParser(format resource.Format, publisher resource.Source, parser Parser)
	GetParser(format resource.Format, source resource.Source) (Parser, error)
}

// parserProperties represents a composite key for the parser map.
type parserProperties struct {
	format    resource.Format
	publisher resource.Source
}

// ParserFactory is a parser's selector, according to the resource format and publisher.
type ParserFactory struct {
	parsers map[parserProperties]Parser
}

// NewParserFactory creates a new factory with predefined default parsers.
func NewParserFactory() *ParserFactory {
	return &ParserFactory{
		parsers: map[parserProperties]Parser{
			{format: resource.JSON, publisher: "nbc-news"}:        &parser.JSONParser{},
			{format: resource.RSS, publisher: "abc-news"}:         &parser.RSSParser{},
			{format: resource.RSS, publisher: "washington-times"}: &parser.RSSParser{},
			{format: resource.RSS, publisher: "bbc-world"}:        &parser.RSSParser{},
			{format: resource.HTML, publisher: "usa-today"}:       &parser.USATodayHTMLParser{},
		},
	}
}

// AddNewParser registers a parser with a specific format and publisher.
func (f *ParserFactory) AddNewParser(format resource.Format, publisher resource.Source, parser Parser) {
	key := parserProperties{format: format, publisher: publisher}
	f.parsers[key] = parser
}

// GetParser returns a parser for the given resource.
func (f *ParserFactory) GetParser(format resource.Format, publisher resource.Source) (Parser, error) {
	key := parserProperties{format: format, publisher: publisher}
	p, exists := f.parsers[key]
	if !exists {
		return nil, fmt.Errorf("no parser found for format: %d and publisher: %s", format, publisher)
	}
	return p, nil
}
