package resource_manager

import (
	"fmt"
	"io"
	"net/http"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/storage"
)

// ResourceDetails is a struct that contains the format and link of a resource.
type ResourceDetails struct {
	Format resource.Format
	Link   string
}

// ResourceManager is a manager that responsible for retrieval of resources from the storage,
// forming them into structures.
type ResourceManager struct {
	storage   *storage.Storage
	resources map[resource.Source]ResourceDetails
}

// New creates a new ResourceManager.
func New(storagePath string) *ResourceManager {
	s := storage.New(storagePath)
	rFormats := map[resource.Source]ResourceDetails{
		"nbc-news":         {Format: resource.JSON, Link: "https://www.nbcnews.com/rss.xml"},
		"abc-news":         {Format: resource.RSS, Link: "https://feeds.abcnews.com/abcnews/internationalheadlines"},
		"washington-times": {Format: resource.RSS, Link: "https://www.washingtontimes.com/rss/headlines/news/world/"},
		"bbc-world":        {Format: resource.RSS, Link: "https://feeds.bbci.co.uk/news/rss.xml"},
		"usa-today":        {Format: resource.HTML, Link: "https://www.nbcnews.com/rss.xml"},
	}

	return &ResourceManager{
		storage:   s,
		resources: rFormats,
	}
}

// SourceIsSupported checks if the source is supported.
func (rm *ResourceManager) SourceIsSupported(source resource.Source) bool {
	_, exists := rm.resources[source]
	return exists
}

// AvailableSources returns the available sources.
func (rm *ResourceManager) AvailableSources() string {
	sources, err := rm.storage.AvailableSources()

	if err != nil {
		return fmt.Sprintf("error getting available sources: %v", err)
	}

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

		fetchedResources = append(fetchedResources, res...)
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

			fetchedResources = append(fetchedResources, res...)
		} else {
			return nil, fmt.Errorf("source \"%s\" is not available", name)
		}
	}

	return fetchedResources, nil
}

// UpdateResource updates the source in the storage.
func (rm *ResourceManager) UpdateResource(source resource.Source) error {
	details, exists := rm.resources[source]
	if !exists {
		return fmt.Errorf("source \"%s\" is not supported", source)
	}

	switch details.Format {
	case resource.JSON:
		return fmt.Errorf("format JSON is not supported")
	case resource.RSS:
		return rm.updateRSSResource(source, details)
	case resource.HTML:
		return fmt.Errorf("format HTML is not supported")
	default:
		return fmt.Errorf("unknown format")
	}
}

func (rm *ResourceManager) getResource(source resource.Source) ([]resource.Resource, error) {
	resContent, err := rm.storage.ReadSource(source)

	if err != nil {
		return []resource.Resource{}, fmt.Errorf("error reading file: %v", err)
	}

	resources := make([]resource.Resource, 0)

	for _, content := range resContent {
		res, err := resource.New(source, rm.resources[source].Format, resource.Content(content))
		if err != nil {
			return resources, fmt.Errorf("error creating resource: %v", err)
		}
		resources = append(resources, *res)
	}

	return resources, nil
}

func (rm *ResourceManager) updateRSSResource(source resource.Source, details ResourceDetails) error {
	resp, err := http.Get(details.Link)
	if err != nil {
		return fmt.Errorf("error fetching resource from link: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching resource from link: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading resource content: %v", err)
	}

	err = rm.storage.UpdateXMLSource(source, body)
	if err != nil {
		return err
	}

	return nil
}
