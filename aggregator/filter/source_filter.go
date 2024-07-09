package filter

import "news-aggregator/aggregator/model/article"

// SourceFilter is an aggregator.Filter that creates a subset from a given set of article.Article's
// corresponding to a given sources set.
type SourceFilter struct {
	sources map[string]struct{}
}

// NewSourceFilter creates a new SourceFilter instance with the given sources.
func NewSourceFilter(sources []string) *SourceFilter {

	sourceSet := make(map[string]struct{})

	for _, source := range sources {
		sourceSet[source] = struct{}{}
	}

	return &SourceFilter{sources: sourceSet}
}

// Apply filters the article.Article's and returns a subset of article.Article's containing the specified sources.
func (f *SourceFilter) Apply(articles []article.Article) []article.Article {

	var filteredArticles []article.Article

	for _, selectedArticle := range articles {
		if f.matchSources(selectedArticle) {
			filteredArticles = append(filteredArticles, selectedArticle)
		}
	}

	return filteredArticles
}

func (f *SourceFilter) matchSources(a article.Article) bool {

	if len(f.sources) == 0 {
		return true
	}
	_, exists := f.sources[string(a.Source())]
	return exists
}
