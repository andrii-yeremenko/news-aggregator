package handler

import "news-aggregator/aggregator/model/resource"

// ResourceManager is a manager that responsible for retrieval of feeds from the storage,
// forming them into structures.
//
//go:generate mockgen -source=resource_manager.go -destination=mocks/mock_resource_manager.go -package=mocks
type ResourceManager interface {
	// AvailableSources returns the available sources.
	AvailableSources() string
	// AvailableFeeds returns the available feeds registered in a system.
	AvailableFeeds() string
	// RegisterSource registers a new source with the given URL and format.
	RegisterSource(name resource.Source, url string, format resource.Format) error
	// UpdateSource updates the source.
	UpdateSource(name resource.Source, url string, format resource.Format) error
	// UpdateResource updates the source in the storage.
	UpdateResource(source resource.Source) error
	// DeleteSource deletes the source.
	DeleteSource(name resource.Source) error
	// IsSourceSupported checks if the source is supported.
	IsSourceSupported(source resource.Source) bool
	// GetSelectedResources returns the specified and known resource.Resource's from a file system.
	GetSelectedResources(sourceNames []string) ([]resource.Resource, error)
	// GetAllResources returns all known resource.Resource's from a file system.
	GetAllResources() ([]resource.Resource, error)
}
