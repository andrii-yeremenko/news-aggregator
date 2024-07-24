package handler

import "news-aggregator/aggregator/model/resource"

//go:generate mockgen -source=resource_manager.go -destination=mocks/mock_resource_manager.go -package=mocks
type ResourceManager interface {
	AvailableSources() string
	RegisterSource(name resource.Source, url string, format resource.Format) error
	UpdateSource(name resource.Source, url string, format resource.Format) error
	UpdateResource(source resource.Source) error
	DeleteSource(name resource.Source) error
	IsSourceSupported(source resource.Source) bool
	GetSelectedResources(sourceNames []string) ([]resource.Resource, error)
	GetAllResources() ([]resource.Resource, error)
}
