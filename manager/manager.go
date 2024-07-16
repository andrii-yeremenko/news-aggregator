package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/storage"
	"os"
	"path/filepath"
	"sort"
)

// ResourceDetails is a struct that contains the format and link of a resource.
type ResourceDetails struct {
	Format resource.Format
	Link   string
}

// ResourceManager is a manager that responsible for retrieval of feeds from the storage,
// forming them into structures.
type ResourceManager struct {
	storage            *storage.Storage
	feeds              map[resource.Source]ResourceDetails
	feedDictionaryPath string
}

// New creates a new ResourceManager.
func New(storagePath string, feedDictionaryPath string) (*ResourceManager, error) {

	dir := filepath.Dir(feedDictionaryPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("error creating directory: %v", err)
		}
	}

	if fileInfo, err := os.Stat(feedDictionaryPath); os.IsNotExist(err) {
		err := createEmptyJSONFile(feedDictionaryPath)
		if err != nil {
			return nil, fmt.Errorf("error creating feed dictionary file: %v", err)
		}
	} else if err == nil && fileInfo.Size() == 0 {
		err := createEmptyJSONFile(feedDictionaryPath)
		if err != nil {
			return nil, fmt.Errorf("error initializing feed dictionary file: %v", err)
		}
	}

	feeds, err := loadResources(feedDictionaryPath)
	if err != nil {
		return nil, fmt.Errorf("error loading feeds: %v", err)
	}

	return &ResourceManager{
		storage:            storage.New(storagePath),
		feeds:              feeds,
		feedDictionaryPath: feedDictionaryPath,
	}, nil
}

// RegisterSource registers a new source.
func (rm *ResourceManager) RegisterSource(name resource.Source, url string, format resource.Format) error {

	rm.feeds[name] = ResourceDetails{
		Format: format,
		Link:   url,
	}

	return rm.saveFeeds()
}

// UpdateSource updates the source.
func (rm *ResourceManager) UpdateSource(name resource.Source, url string, format resource.Format) error {

	rm.feeds[name] = ResourceDetails{
		Format: format,
		Link:   url,
	}

	return rm.saveFeeds()
}

// DeleteSource deletes the source.
func (rm *ResourceManager) DeleteSource(name resource.Source) error {

	delete(rm.feeds, name)

	return rm.saveFeeds()
}

// IsSourceSupported checks if the source is supported.
func (rm *ResourceManager) IsSourceSupported(source resource.Source) bool {
	_, exists := rm.feeds[source]
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

// GetAllResources returns all known resource.Resource's from a file system.
func (rm *ResourceManager) GetAllResources() ([]resource.Resource, error) {

	fetchedResources := make([]resource.Resource, 0)

	for s := range rm.feeds {
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
		if _, exists := rm.feeds[s]; exists {
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

// UpdateAllSources updates all sources in the storage.
func (rm *ResourceManager) UpdateAllSources() error {

	for source := range rm.feeds {
		err := rm.UpdateResource(source)
		if err != nil {
			return fmt.Errorf("error updating source \"%s\": %v", source, err)
		}
	}

	return nil
}

// UpdateResource updates the source in the storage.
func (rm *ResourceManager) UpdateResource(source resource.Source) error {
	details, exists := rm.feeds[source]
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
		res, err := resource.New(source, rm.feeds[source].Format, resource.Content(content))
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

func (rm *ResourceManager) saveFeeds() error {
	resourceList := make([]struct {
		Source string `json:"source"`
		Format string `json:"format"`
		Link   string `json:"link"`
	}, 0, len(rm.feeds))

	for source, details := range rm.feeds {
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

	file, err := os.Create(rm.feedDictionaryPath)

	if err != nil {
		return fmt.Errorf("error creating feeds file: %v", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing feeds file: %v\n", err)
		}
	}(file)

	sort.Slice(resourceList, func(i, j int) bool {
		return resourceList[i].Source < resourceList[j].Source
	})

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&resourceList); err != nil {
		return fmt.Errorf("error encoding feeds file: %v", err)
	}

	return nil
}

func loadResources(path string) (map[resource.Source]ResourceDetails, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening feeds file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing feeds file: %v\n", err)
		}
	}(file)

	var resourceList []struct {
		Source string `json:"source"`
		Format string `json:"format"`
		Link   string `json:"link"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&resourceList); err != nil {
		return nil, fmt.Errorf("error decoding feeds file: %v", err)
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

// createEmptyJSONFile creates a file with an empty JSON array.
func createEmptyJSONFile(path string) error {
	emptyFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(emptyFile *os.File) {
		err := emptyFile.Close()
		if err != nil {
			fmt.Printf("error closing empty JSON file: %v\n", err)
		}
	}(emptyFile)

	_, err = emptyFile.Write([]byte("[]"))
	return err
}
