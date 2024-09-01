package feed

import (
	"errors"
)

// Feed a structured data supply containing supplier information
type Feed struct {
	source Source
	format Format
	link   Link
}

// New is a constructor function for creating a new Feed.
func New(source Source, format Format, link Link) (*Feed, error) {

	if source == "" {
		return nil, errors.New("source cannot be empty")
	}

	if format == 0 {
		return nil, errors.New("format cannot be unknown")
	}

	if link == "" {
		return nil, errors.New("link cannot be empty")
	}

	return &Feed{
		source: source,
		format: format,
		link:   link,
	}, nil
}

// Source returns the source of the feed.
func (r *Feed) Source() Source {
	return r.source
}

// Format returns the format of the feed.
func (r *Feed) Format() Format {
	return r.format
}

// Link returns the link to the feed origin.
func (r *Feed) Link() Link {
	return r.link
}
