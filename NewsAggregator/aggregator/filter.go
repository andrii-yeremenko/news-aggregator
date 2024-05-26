package aggregator

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/parser"
	"github.com/reiver/go-porterstemmer"
	"strings"
	"time"
)

// Filter is a tool that produces from a given set of data a specific subset of data that meets certain conditions.
type Filter struct {
	keywords  []string
	startDate *time.Time
	endDate   *time.Time
	sources   map[string]struct{}
}

type FilterBuilder struct {
	filter Filter
}

func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{filter: Filter{}}
}

// WithKeywords sets the keywords that the filter will use to filter the data.
func (builder *FilterBuilder) WithKeywords(keywords []string) *FilterBuilder {
	builder.filter.keywords = keywords
	return builder
}

// WithStartDate sets the start date that the filter will use to filter the data.
func (builder *FilterBuilder) WithStartDate(startDateStr string) *FilterBuilder {
	dateParser := parser.NewDateParser()
	startDate, err := dateParser.ParseDefaultDateFormat(startDateStr)
	if err != nil {
		panic("invalid start date")
	}
	builder.filter.startDate = &startDate
	return builder
}

// WithEndDate sets the end date that the filter will use to filter the data.
func (builder *FilterBuilder) WithEndDate(endDateStr string) *FilterBuilder {
	dateParser := parser.NewDateParser()
	endDate, err := dateParser.ParseDefaultDateFormat(endDateStr)
	if err != nil {
		panic("invalid end date")
	}
	builder.filter.endDate = &endDate
	return builder
}

// WithSources sets the sources that the filter will use to filter the data.
func (builder *FilterBuilder) WithSources(sources []string) *FilterBuilder {
	sourceSet := make(map[string]struct{})
	for _, source := range sources {
		sourceSet[source] = struct{}{}
	}
	builder.filter.sources = sourceSet
	return builder
}

// Build creates a new Filter instance based on the provided configuration.
func (builder *FilterBuilder) Build() Filter {
	return builder.filter
}

// Apply filters the provided articles based on the filter configuration.
func (f *Filter) Apply(articles []article.Article) []article.Article {
	var filteredArticles []article.Article

	for _, art := range articles {
		if f.matchKeywords(art) && f.matchDateRange(art) && f.matchSources(art) {
			filteredArticles = append(filteredArticles, art)
		}
	}

	return filteredArticles
}

func (f *Filter) matchKeywords(art article.Article) bool {
	if len(f.keywords) == 0 {
		return true
	}
	title := strings.ToLower(string(art.Title()))
	description := strings.ToLower(string(art.Description()))
	for _, keyword := range f.keywords {
		keyword = strings.ToLower(keyword)
		stemmedKeyword := porterstemmer.StemString(keyword)
		if strings.Contains(title, stemmedKeyword) || strings.Contains(description, stemmedKeyword) {
			return true
		}
	}
	return false
}

func (f *Filter) matchDateRange(art article.Article) bool {
	if f.startDate == nil && f.endDate == nil {
		return true
	}
	articleDate := time.Time(art.Date())
	if f.startDate != nil && articleDate.Before(*f.startDate) {
		return false
	}
	if f.endDate != nil && articleDate.After(*f.endDate) {
		return false
	}
	return true
}

func (f *Filter) matchSources(art article.Article) bool {
	if len(f.sources) == 0 {
		return true
	}
	_, exists := f.sources[string(art.Source())]
	return exists
}
