package manager

import (
	"fmt"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/storage"
)

// ResourceManager is a manager that responsible for retrieval of resources from the storage,
// forming them into structures.
type ResourceManager struct {
	storage   *storage.Storage
	resources map[resource.Source]resource.Format
}

// New creates a new ResourceManager.
func New(storagePath string) *ResourceManager {
	s := storage.New(storagePath)
	rFormats := map[resource.Source]resource.Format{
		"nbc-news":         resource.JSON,
		"abc-news":         resource.RSS,
		"washington-times": resource.RSS,
		"bbc-world":        resource.RSS,
		"usa-today":        resource.HTML,
	}

	return &ResourceManager{
		storage:   s,
		resources: rFormats,
	}
}

// AvailableSources returns the available sources.
func (rm *ResourceManager) AvailableSources() string {
	sources := rm.storage.AvailableSources()

	if len(sources) == 0 {
		return "no available sources"
	}

	sourcesStr := ""

	for _, source := range sources {
		sourcesStr += string(source) + ","
	}

	return sourcesStr
}

// AllResources returns all known resource.Resource's from a file system.
func (rm *ResourceManager) AllResources() ([]resource.Resource, error) {

	fetchedResources := make([]resource.Resource, 0)

	for s := range rm.resources {
		res, err := rm.getResource(s)
		if err != nil {
			return fetchedResources, fmt.Errorf("error getting resource : %v", err)
		}

		fetchedResources = append(fetchedResources, res)
	}

	return fetchedResources, nil

}

// GetSelectedResources returns the specified and known resource.Resource's from a file system.
func (rm *ResourceManager) GetSelectedResources(sourceNames []string) ([]resource.Resource, error) {

	fetchedResources := make([]resource.Resource, 0)

	for _, name := range sourceNames {
		s := resource.Source(name)
		if _, exists := rm.resources[s]; exists {
			res, err := rm.getResource(s)
			if err != nil {
				return nil, fmt.Errorf("error getting resource from source \"%s\" : %v", name, err)
			}

			fetchedResources = append(fetchedResources, res)
		} else {
			return nil, fmt.Errorf("source \"%s\" is not available", name)
		}
	}

	return fetchedResources, nil
}

func (rm *ResourceManager) getResource(source resource.Source) (resource.Resource, error) {
	resContent, err := rm.storage.ReadSource(source)

	if err != nil {
		return resource.Resource{}, fmt.Errorf("error reading file: %v", err)
	}

	res, err := resource.New(source, rm.resources[source], resource.Content(resContent))

	if err != nil {
		return resource.Resource{}, fmt.Errorf("error creating resource: %v", err)
	}

	return *res, nil
}
