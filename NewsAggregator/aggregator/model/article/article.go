package article

import (
	"errors"
	"time"
)

// Article is a piece of writing about a particular subject.
type Article struct {
	title        Title
	description  Description
	creationDate CreationDate
	source       Source
	author       Author
}

func (a *Article) Title() Title {
	return a.title
}

func (a *Article) Description() Description {
	return a.description
}

func (a *Article) Date() CreationDate {
	return a.creationDate
}

func (a *Article) Source() Source {
	return a.source
}

func (a *Article) Author() Author {
	return a.author
}

type Builder struct {
	article *Article
}

func NewArticleBuilder() *Builder {
	return &Builder{article: &Article{}}
}

func (builder *Builder) SetTitle(title Title) *Builder {
	builder.article.title = title
	return builder
}

func (builder *Builder) SetDescription(description Description) *Builder {
	builder.article.description = description
	return builder
}

func (builder *Builder) SetDate(creationDate CreationDate) *Builder {
	builder.article.creationDate = creationDate
	return builder
}

func (builder *Builder) SetSource(source Source) *Builder {
	builder.article.source = source
	return builder
}

func (builder *Builder) SetAuthor(author Author) *Builder {
	builder.article.author = author
	return builder
}

// Build validates the article and returns the final Article instance.
// Checks all required fields are set. If not, returns an error.
// If all fields are set, returns the Article instance.
func (builder *Builder) Build() (*Article, error) {

	if builder.article.title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if builder.article.description == "" {
		return nil, errors.New("description cannot be empty")
	}
	if time.Time(builder.article.creationDate).IsZero() {
		return nil, errors.New("creationDate cannot be empty")
	}
	if builder.article.source == "" {
		return nil, errors.New("source cannot be empty")
	}

	return builder.article, nil
}
