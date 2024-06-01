package storage

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
	newsRepository  *NewsStorage
	resourceDetails map[resource.Source]resourceDetail
}

// NewLoader creates a new ResourceLoader.
func NewLoader() *ResourceLoader {
	return &ResourceLoader{
		newsRepository: New(),
		resourceDetails: map[resource.Source]resourceDetail{
			"nbc-news":         {format: "json", path: "storage/resources/nbc-news.json"},
			"abc-news":         {format: "rss", path: "storage/resources/abc-news.xml"},
			"washington-times": {format: "rss", path: "storage/resources/washington-times.xml"},
			"bbc-world":        {format: "rss", path: "storage/resources/bbc-world.xml"},
			"usa-today":        {format: "html", path: "storage/resources/usa-today-world-news.html"},
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

// LoadAllResources loads all resource.Resource's into the given aggregator.Aggregator.
func (l *ResourceLoader) LoadAllResources(agr *aggregator.Aggregator) {
	for source, detail := range l.resourceDetails {
		l.LoadResource(source, detail, agr)
	}
}

// LoadSelectedResources loads the specified and registered resource.Resource's into the given aggregator.Aggregator.
func (l *ResourceLoader) LoadSelectedResources(sourceNames []string, agr *aggregator.Aggregator) {
	for _, sourceName := range sourceNames {
		source := resource.Source(sourceName)
		if detail, exists := l.resourceDetails[source]; exists {
			l.LoadResource(source, detail, agr)
		}
	}
}

// LoadResource is a helper function to load a single resource.Resource into the given aggregator.Aggregator.
func (l *ResourceLoader) LoadResource(source resource.Source, detail resourceDetail, agr *aggregator.Aggregator) {
	res, err := l.newsRepository.ReadFile(source, detail.format, detail.path)
	if err != nil {
		panic(err)
	}
	err = agr.LoadResource(res)
	if err != nil {
		panic(err)
	}
}
