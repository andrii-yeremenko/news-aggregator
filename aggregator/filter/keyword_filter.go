package filter

import (
	"github.com/reiver/go-porterstemmer"
	"news-aggregator/aggregator/model/article"
	"strings"
)

// KeywordFilter is an aggregator.Filter that creates a subset from a given set of article.Article's
// corresponding to a given keywords set.
type KeywordFilter struct {
	keywords []string
}

// NewKeywordFilter creates a new KeywordFilter instance with the given keywords.
func NewKeywordFilter(keywords []string) *KeywordFilter {
	return &KeywordFilter{keywords: keywords}
}

// Apply filters the article.Article's and returns a subset of article.Article's containing the specified keywords.
func (f *KeywordFilter) Apply(articles []article.Article) []article.Article {

	var filteredArticles []article.Article

	for _, selectedArticle := range articles {
		if f.matchKeywords(selectedArticle) {
			filteredArticles = append(filteredArticles, selectedArticle)
		}
	}

	return filteredArticles
}

// matchKeywords checks if the given article contains any of the keywords.
func (f *KeywordFilter) matchKeywords(a article.Article) bool {

	if len(f.keywords) == 0 {
		return true
	}

	title := strings.ToLower(string(a.Title()))
	description := strings.ToLower(string(a.Description()))
	for _, keyword := range f.keywords {
		keyword = strings.ToLower(keyword)
		stemmedKeyword := porterstemmer.StemString(keyword)
		if strings.Contains(title, stemmedKeyword) || strings.Contains(description, stemmedKeyword) {
			return true
		}
	}
	return false
}
