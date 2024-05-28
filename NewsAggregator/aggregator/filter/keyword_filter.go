package filter

import (
	"NewsAggregator/aggregator/model/article"
	"github.com/reiver/go-porterstemmer"
	"strings"
)

// KeywordFilter is a Filter that creates a subset from a given set of article.Article's
// corresponding to a given keywords set.
type KeywordFilter struct {
	keywords []string
}

// NewKeywordFilter creates a new KeywordFilter instance.
func NewKeywordFilter(keywords []string) *KeywordFilter {
	return &KeywordFilter{keywords: keywords}
}

// Apply filters the data and returns a subset of articles.
func (filter *KeywordFilter) Apply(articles []article.Article) []article.Article {

	var filteredArticles []article.Article

	for _, selectedArticle := range articles {
		if filter.matchKeywords(selectedArticle) {
			filteredArticles = append(filteredArticles)
		}
	}

	return filteredArticles
}

func (filter *KeywordFilter) matchKeywords(art article.Article) bool {
	if len(filter.keywords) == 0 {
		return true
	}
	title := strings.ToLower(string(art.Title()))
	description := strings.ToLower(string(art.Description()))
	for _, keyword := range filter.keywords {
		keyword = strings.ToLower(keyword)
		stemmedKeyword := porterstemmer.StemString(keyword)
		if strings.Contains(title, stemmedKeyword) || strings.Contains(description, stemmedKeyword) {
			return true
		}
	}
	return false
}
