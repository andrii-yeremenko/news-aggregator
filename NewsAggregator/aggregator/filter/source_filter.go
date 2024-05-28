package filter

import "NewsAggregator/aggregator/model/article"

// SourceFilter is a Filter that creates a subset from a given set of article.Article's
// corresponding to a given sources set.
type SourceFilter struct {
	sources map[string]struct{}
}

// NewSourceFilter creates a new SourceFilter instance.
func NewSourceFilter(sources []string) *SourceFilter {

	sourceSet := make(map[string]struct{})

	for _, source := range sources {
		sourceSet[source] = struct{}{}
	}

	return &SourceFilter{sources: sourceSet}
}

// Apply filters the data and returns a subset of articles.
func (filter *SourceFilter) Apply(articles []article.Article) []article.Article {

	var filteredArticles []article.Article

	for _, selectedArticle := range articles {
		if filter.matchSources(selectedArticle) {
			filteredArticles = append(filteredArticles, selectedArticle)
		}
	}

	return filteredArticles
}

func (filter *SourceFilter) matchSources(art article.Article) bool {
	if len(filter.sources) == 0 {
		return true
	}
	_, exists := filter.sources[string(art.Source())]
	return exists
}
