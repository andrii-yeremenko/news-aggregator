package resource

import (
	"errors"
)

// Resource a data supply containing a news report.
type Resource struct {
	publisher Source
	format    Format
	content   Content
}

// NewResource is a constructor function for creating a new Resource.
func NewResource(source Source, format Format, content Content) (*Resource, error) {

	if source == "" {
		return nil, errors.New("source cannot be empty")
	}

	if format == "" {
		return nil, errors.New("format cannot be empty")
	}

	if content == "" {
		return nil, errors.New("content cannot be empty")
	}

	return &Resource{
		publisher: source,
		format:    format,
		content:   content,
	}, nil
}

// Source returns the publisher of the resource.
func (r *Resource) Source() Source {
	return r.publisher
}

// Format returns the format of the resource file.
func (r *Resource) Format() Format {
	return r.format
}

// Content returns the content of the resource.
func (r *Resource) Content() Content {
	return r.content
}
