package article

import (
	"NewsAggregator/aggregator/model/resource"
	"errors"
	"time"
)

// The Article is a structured piece of writing about a particular subject.
// It consists of a title, description, creation date, source, and optionally author and link.
type Article struct {
	title        Title
	description  Description
	creationDate CreationDate
	source       resource.Source
	author       Author
	link         Link
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

func (a *Article) Source() resource.Source {
	return a.source
}

func (a *Article) Author() Author {
	return a.author
}

func (a *Article) Link() Link {
	return a.link
}

type Builder struct {
	article *Article
}

func NewArticleBuilder() *Builder {
	return &Builder{article: &Article{}}
}

func (b *Builder) SetTitle(title Title) *Builder {
	b.article.title = title
	return b
}

func (b *Builder) SetDescription(description Description) *Builder {
	b.article.description = description
	return b
}

func (b *Builder) SetDate(creationDate CreationDate) *Builder {
	b.article.creationDate = creationDate
	return b
}

func (b *Builder) SetSource(source resource.Source) *Builder {
	b.article.source = source
	return b
}

func (b *Builder) SetAuthor(author Author) *Builder {
	b.article.author = author
	return b
}

func (b *Builder) SetLink(link Link) *Builder {
	b.article.link = link
	return b
}

// Build validates the article and returns the final Article instance.
// Checks all required fields are set. If not, returns an error.
// If all fields are set, return the Article instance.
func (b *Builder) Build() (*Article, error) {

	if b.article.title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if b.article.description == "" {
		return nil, errors.New("description cannot be empty")
	}
	if time.Time(b.article.creationDate).IsZero() {
		return nil, errors.New("creationDate cannot be empty")
	}
	if b.article.source == "" {
		return nil, errors.New("source cannot be empty")
	}

	return b.article, nil
}
