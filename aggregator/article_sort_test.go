package aggregator

import (
	"news-aggregator/aggregator/model/article"
	"reflect"
	"testing"
	"time"
)

func TestSortArticlesByDateAsc(t *testing.T) {
	articleA, _ := article.NewArticleBuilder().
		SetDate(article.CreationDate(time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC))).
		SetTitle("Article A").
		SetDescription("A").
		SetSource("nbc").
		Build()

	articleB, _ := article.NewArticleBuilder().
		SetDate(article.CreationDate(time.Date(2023, 6, 2, 0, 0, 0, 0, time.UTC))).
		SetTitle("Article A").
		SetDescription("A").
		SetSource("nbc").
		Build()

	articleC, _ := article.NewArticleBuilder().
		SetDate(article.CreationDate(time.Date(2023, 6, 3, 0, 0, 0, 0, time.UTC))).
		SetTitle("Article A").
		SetDescription("A").
		SetSource("nbc").
		Build()

	articles := []article.Article{*articleA, *articleC, *articleB}

	expectedAsc := []article.Article{*articleA, *articleB, *articleC}

	sortedAsc := SortArticlesByDateAsc(articles)

	if !reflect.DeepEqual(sortedAsc, expectedAsc) {
		t.Errorf("SortArticlesByDateAsc returned %v, expected %v", sortedAsc, expectedAsc)
	}
}

func TestSortArticlesByDateDesc(t *testing.T) {
	articleA, _ := article.NewArticleBuilder().
		SetDate(article.CreationDate(time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC))).
		SetTitle("Article A").
		SetDescription("A").
		SetSource("nbc").
		Build()

	articleB, _ := article.NewArticleBuilder().
		SetDate(article.CreationDate(time.Date(2023, 6, 2, 0, 0, 0, 0, time.UTC))).
		SetTitle("Article A").
		SetDescription("A").
		SetSource("nbc").
		Build()

	articleC, _ := article.NewArticleBuilder().
		SetDate(article.CreationDate(time.Date(2023, 6, 3, 0, 0, 0, 0, time.UTC))).
		SetTitle("Article A").
		SetDescription("A").
		SetSource("nbc").
		Build()

	articles := []article.Article{*articleA, *articleC, *articleB}

	expectedDesc := []article.Article{*articleC, *articleB, *articleA}

	sortedDesc := SortArticlesByDateDesc(articles)

	if !reflect.DeepEqual(sortedDesc, expectedDesc) {
		t.Errorf("SortArticlesByDateDesc returned %v, expected %v", sortedDesc, expectedDesc)
	}
}
