package aggregator

import (
	"news-aggregator/aggregator/model/article"
	"sort"
	"time"
)

// SortArticlesByDateAsc sorts the given array of articles by date in ascending order
func SortArticlesByDateAsc(articles []article.Article) []article.Article {
	sortedArticles := make([]article.Article, len(articles))
	copy(sortedArticles, articles)

	sort.Slice(sortedArticles, func(i, j int) bool {
		return time.Time(sortedArticles[i].Date()).Before(time.Time(sortedArticles[j].Date()))
	})
	return sortedArticles
}

// SortArticlesByDateDesc sorts the given array of articles by date in descending order
func SortArticlesByDateDesc(articles []article.Article) []article.Article {
	sortedArticles := make([]article.Article, len(articles))
	copy(sortedArticles, articles)

	sort.Slice(sortedArticles, func(i, j int) bool {
		return time.Time(sortedArticles[i].Date()).After(time.Time(sortedArticles[j].Date()))
	})
	return sortedArticles
}
