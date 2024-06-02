package resource

import (
	"errors"
)

// Resource a structured data supply containing a news report content with information about its parameters required
// for future processing.
type Resource struct {
	source  Source
	format  Format
	content Content
}

// New is a constructor function for creating a new Resource.
func New(source Source, format Format, content Content) (*Resource, error) {

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
		source:  source,
		format:  format,
		content: content,
	}, nil
}

// Source returns the source of the resource.
func (r *Resource) Source() Source {
	return r.source
}

// Format returns the format of the resource file.
func (r *Resource) Format() Format {
	return r.format
}

// Content returns the content of the resource.
func (r *Resource) Content() Content {
	return r.content
}
