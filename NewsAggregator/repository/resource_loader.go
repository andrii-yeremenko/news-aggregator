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

// NewResourceLoader creates a new ResourceLoader.
func NewResourceLoader(newsAggregator *aggregator.Aggregator) *ResourceLoader {
	return &ResourceLoader{
		newsAggregator: newsAggregator,
		newsRepository: NewRepository(),
		resourceDetails: map[resource.Source]resourceDetail{
			"nbc-news":         {format: "json", path: "repository/news-resources/nbc-news.json"},
			"abc-news":         {format: "rss", path: "repository/news-resources/abc-news.xml"},
			"washington-times": {format: "rss", path: "repository/news-resources/washington-times.xml"},
			"bbc-world":        {format: "rss", path: "repository/news-resources/bbc-world.xml"},
			"usa-today":        {format: "html", path: "repository/news-resources/usa-today-world-news.html"},
		},
	}
}

// GetAvailableSources returns the available sources.
func (loader *ResourceLoader) GetAvailableSources() string {
	var sources string
	for source := range loader.resourceDetails {
		sources += string(source) + ", "
	}
	return sources
}

// RegisterResource adds a resource to be loaded.
func (loader *ResourceLoader) RegisterResource(source resource.Source, format resource.Format, path string) {
	loader.resourceDetails[source] = resourceDetail{format: format, path: path}
}

// LoadAllResources loads all resources into the aggregator.
func (loader *ResourceLoader) LoadAllResources() {
	for source, detail := range loader.resourceDetails {
		loader.loadResource(source, detail)
	}
}

// LoadSelectedResources loads the specified and registered resources into the aggregator.
func (loader *ResourceLoader) LoadSelectedResources(sourceNames []string) {
	for _, sourceName := range sourceNames {
		source := resource.Source(sourceName)
		if detail, exists := loader.resourceDetails[source]; exists {
			loader.loadResource(source, detail)
		}
	}
}

// loadResource is a helper function to load a single resource into the aggregator.
func (loader *ResourceLoader) loadResource(source resource.Source, detail resourceDetail) {
	res, err := loader.newsRepository.ReadFile(source, detail.format, detail.path)
	if err != nil {
		panic(err)
	}
	err = loader.newsAggregator.LoadResource(res)
	if err != nil {
		panic(err)
	}
}
