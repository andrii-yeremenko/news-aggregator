package repository

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/aggregator/model/resource"
)

// resourceDetail holds the details of a resource to be loaded.
type resourceDetail struct {
	format resource.Format
	path   string
}

// ResourceLoader is responsible for loading resources into the aggregator.
type ResourceLoader struct {
	newsAggregator  *aggregator.Aggregator
	newsRepository  *NewsRepository
	resourceDetails map[resource.Source]resourceDetail
}

// NewLoader creates a new ResourceLoader.
func NewLoader(a *aggregator.Aggregator) *ResourceLoader {
	return &ResourceLoader{
		newsAggregator: a,
		newsRepository: New(),
		resourceDetails: map[resource.Source]resourceDetail{
			"nbc-news":         {format: "json", path: "repository/resources/nbc-news.json"},
			"abc-news":         {format: "rss", path: "repository/resources/abc-news.xml"},
			"washington-times": {format: "rss", path: "repository/resources/washington-times.xml"},
			"bbc-world":        {format: "rss", path: "repository/resources/bbc-world.xml"},
			"usa-today":        {format: "html", path: "repository/resources/usa-today-world-news.html"},
		},
	}
}

// GetAvailableSources returns the available sources.
func (l *ResourceLoader) GetAvailableSources() string {
	var sources string
	for source := range l.resourceDetails {
		sources += string(source) + ", "
	}
	return sources
}

// RegisterResource adds a resource to be loaded.
func (l *ResourceLoader) RegisterResource(source resource.Source, format resource.Format, path string) {
	l.resourceDetails[source] = resourceDetail{format: format, path: path}
}

// LoadAllResources loads all resources into the aggregator.
func (l *ResourceLoader) LoadAllResources() {
	for source, detail := range l.resourceDetails {
		l.loadResource(source, detail)
	}
}

// LoadSelectedResources loads the specified and registered resources into the aggregator.
func (l *ResourceLoader) LoadSelectedResources(sourceNames []string) {
	for _, sourceName := range sourceNames {
		source := resource.Source(sourceName)
		if detail, exists := l.resourceDetails[source]; exists {
			l.loadResource(source, detail)
		}
	}
}

// loadResource is a helper function to load a single resource into the aggregator.
func (l *ResourceLoader) loadResource(source resource.Source, detail resourceDetail) {
	res, err := l.newsRepository.ReadFile(source, detail.format, detail.path)
	if err != nil {
		panic(err)
	}
	err = l.newsAggregator.LoadResource(res)
	if err != nil {
		panic(err)
	}
}
