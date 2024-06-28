package resource_manager

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/storage"
	"os"
)

// ResourceDetails is a struct that contains the format and link of a resource.
type ResourceDetails struct {
	Format resource.Format
	Link   string
}

// ResourceManager is a manager that responsible for retrieval of resources from the storage,
// forming them into structures.
type ResourceManager struct {
	storage        *storage.Storage
	resources      map[resource.Source]ResourceDetails
	resourcesSetup string
}

// New creates a new ResourceManager.
func New(storagePath string, resourcesPath string) (*ResourceManager, error) {
	s := storage.New(storagePath)

	rFormats, err := loadResources(resourcesPath)
	if err != nil {
		return nil, fmt.Errorf("error loading resources: %v", err)
	}

	return &ResourceManager{
		storage:        s,
		resources:      rFormats,
		resourcesSetup: resourcesPath,
	}, nil
}

// RegisterSource registers a new source.
func (rm *ResourceManager) RegisterSource(name resource.Source, url string, format resource.Format) error {

	rm.resources[name] = ResourceDetails{
		Format: format,
		Link:   url,
	}

	return saveResources(rm.resourcesSetup, rm.resources)
}

// UpdateSource updates the source.
func (rm *ResourceManager) UpdateSource(name resource.Source, url string, format resource.Format) error {

	rm.resources[name] = ResourceDetails{
		Format: format,
		Link:   url,
	}

	return saveResources(rm.resourcesSetup, rm.resources)
}

// DeleteSource deletes the source.
func (rm *ResourceManager) DeleteSource(name string) error {

	source := resource.Source(name)
	delete(rm.resources, source)

	return saveResources(rm.resourcesSetup, rm.resources)
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
	case resource.RSS:
		return rm.updateRSSResource(source, details)
	case resource.HTML:
		return rm.updateHTMLResource(source, details)
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

func (rm *ResourceManager) updateHTMLResource(source resource.Source, details ResourceDetails) error {
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

	err = rm.storage.UpdateHTMLSource(source, body)
	if err != nil {
		return err
	}

	return nil
}

func loadResources(resourcesPath string) (map[resource.Source]ResourceDetails, error) {
	file, err := os.Open(resourcesPath)
	if err != nil {
		return nil, fmt.Errorf("error opening resources file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing resources file: %v\n", err)
		}
	}(file)

	var resourceList []struct {
		Source string `json:"source"`
		Format string `json:"format"`
		Link   string `json:"link"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&resourceList); err != nil {
		return nil, fmt.Errorf("error decoding resources file: %v", err)
	}

	rFormats := make(map[resource.Source]ResourceDetails)
	for _, res := range resourceList {
		format, err := resource.ParseFormat(res.Format)
		if err != nil {
			return nil, err
		}

		rFormats[resource.Source(res.Source)] = ResourceDetails{
			Format: format,
			Link:   res.Link,
		}
	}

	return rFormats, nil
}

func saveResources(resourcesPath string, resources map[resource.Source]ResourceDetails) error {
	resourceList := make([]struct {
		Source string `json:"source"`
		Format string `json:"format"`
		Link   string `json:"link"`
	}, 0, len(resources))

	for source, details := range resources {
		resourceList = append(resourceList, struct {
			Source string `json:"source"`
			Format string `json:"format"`
			Link   string `json:"link"`
		}{
			Source: string(source),
			Format: resource.FormatToString(details.Format),
			Link:   details.Link,
		})
	}

	file, err := os.Create(resourcesPath)
	if err != nil {
		return fmt.Errorf("error creating resources file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing resources file: %v\n", err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&resourceList); err != nil {
		return fmt.Errorf("error encoding resources file: %v", err)
	}

	return nil
}
