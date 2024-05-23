package parser

import (
	"NewsAggregator/aggregator/model/resource"
	"fmt"
)

// parserProperties represents a composite key for the parser map.
type parserProperties struct {
	format    resource.Format
	publisher resource.Publisher
}

// Factory is a parser's selector according to the resource format and publisher.
type Factory struct {
	parsers map[parserProperties]Parser
}

// NewParserFactory creates a new factory.
func NewParserFactory() *Factory {
	return &Factory{
		parsers: make(map[parserProperties]Parser),
	}
}

// RegisterParser registers a parser with a specific format and publisher.
func (factory *Factory) RegisterParser(format resource.Format, publisher resource.Publisher, parser Parser) {
	key := parserProperties{format: format, publisher: publisher}
	factory.parsers[key] = parser
}

// GetParser returns a parser for the given resource.
func (factory *Factory) GetParser(format resource.Format, publisher resource.Publisher) (Parser, error) {
	key := parserProperties{format: format, publisher: publisher}
	parser, exists := factory.parsers[key]
	if !exists {
		return nil, fmt.Errorf("no parser found for format: %s and publisher: %s", format, publisher)
	}
	return parser, nil
}
