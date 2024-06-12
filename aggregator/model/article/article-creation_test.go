package article_test

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"testing"
	"time"
)

func TestArticleCreation(t *testing.T) {
	builder := article.NewArticleBuilder()

	title := article.Title("Test Title")
	description := article.Description("Test Description")
	date := article.CreationDate(time.Now())
	source := resource.Source("Test Source")
	author := article.Author("Test Author")
	link := article.Link("http://testlink.com")

	art, err := builder.
		SetTitle(title).
		SetDescription(description).
		SetDate(date).
		SetSource(source).
		SetAuthor(author).
		SetLink(link).
		Build()

	if err != nil {
		t.Errorf("Error occurred while creating art: %v", err)
	}

	if art.Title() != title {
		t.Errorf("Expected title: %s, Got: %s", title, art.Title())
	}
	if art.Description() != description {
		t.Errorf("Expected description: %s, Got: %s", description, art.Description())
	}
	if art.Date() != date {
		t.Errorf("Expected date: %v, Got: %v", date, art.Date())
	}
	if art.Source() != source {
		t.Errorf("Expected source: %s, Got: %s", source, art.Source())
	}
	if art.Author() != author {
		t.Errorf("Expected author: %s, Got: %s", author, art.Author())
	}
	if art.Link() != link {
		t.Errorf("Expected link: %s, Got: %s", link, art.Link())
	}
	if art.TitleStr() != string(title) {
		t.Errorf("Expected title string: %s, Got: %s", string(title), art.TitleStr())
	}
	if art.DescriptionStr() != string(description) {
		t.Errorf("Expected description string: %s, Got: %s", string(description), art.DescriptionStr())
	}
}

func TestEmptyTitle(t *testing.T) {
	builder := article.NewArticleBuilder()

	_, err := builder.
		SetDescription("Test Description").
		SetDate(article.CreationDate(time.Now())).
		SetSource("Test Source").
		Build()
	if err == nil || err.Error() != "title cannot be empty" {
		t.Errorf("Expected error: 'title cannot be empty', Got: %v", err)
	}
}

func TestEmptyDescription(t *testing.T) {
	builder := article.NewArticleBuilder()

	_, err := builder.
		SetTitle("Test Title").
		SetDate(article.CreationDate(time.Now())).
		SetSource("Test Source").
		Build()
	if err == nil || err.Error() != "description cannot be empty" {
		t.Errorf("Expected error: 'description cannot be empty', Got: %v", err)
	}
}

func TestEmptyCreationDate(t *testing.T) {
	builder := article.NewArticleBuilder()

	_, err := builder.
		SetTitle("Test Title").
		SetDescription("Test Description").
		SetSource("Test Source").
		Build()
	if err == nil || err.Error() != "creationDate cannot be empty" {
		t.Errorf("Expected error: 'creationDate cannot be empty', Got: %v", err)
	}
}

func TestEmptySource(t *testing.T) {
	builder := article.NewArticleBuilder()

	_, err := builder.
		SetTitle("Test Title").
		SetDescription("Test Description").
		SetDate(article.CreationDate(time.Now())).
		Build()
	if err == nil || err.Error() != "source cannot be empty" {
		t.Errorf("Expected error: 'source cannot be empty', Got: %v", err)
	}
}
