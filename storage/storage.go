package storage

import (
	"bufio"
	"fmt"
	"news-aggregator/aggregator/model/resource"
	"os"
	"path/filepath"
)

// ResourceDetail holds the details of a resource to be loaded.
type ResourceDetail struct {
	format resource.Format
	path   string
}

// Storage is a tool for loading structured resources into the aggregator.Aggregator.
type Storage struct {
	resourceDetails map[resource.Source]ResourceDetail
	basePath        string
}

// New creates a new Storage.
func New(basePath string) *Storage {
	return &Storage{
		resourceDetails: map[resource.Source]ResourceDetail{
			"nbc-news":         {format: resource.JSON, path: "resources/nbc-news.json"},
			"abc-news":         {format: resource.RSS, path: "resources/abc-news.xml"},
			"washington-times": {format: resource.RSS, path: "resources/washington-times.xml"},
			"bbc-world":        {format: resource.RSS, path: "resources/bbc-world.xml"},
			"usa-today":        {format: resource.HTML, path: "resources/usa-today-world-news.html"},
		},
		basePath: basePath,
	}
}

// GetAvailableSources returns the available sources.
func (l *Storage) GetAvailableSources() string {
	var sources string
	for source := range l.resourceDetails {
		sources += string(source) + ", "
	}
	return sources
}

// GetAllResources returns all known resource.Resource's from a file system.
func (l *Storage) GetAllResources() ([]resource.Resource, error) {

	fetchedResources := make([]resource.Resource, 0)

	for s, detail := range l.resourceDetails {
		res, err := l.getResource(s, detail)
		if err != nil {
			return fetchedResources, fmt.Errorf("error getting resource : %v", err)
		}

		fetchedResources = append(fetchedResources, res)
	}

	return fetchedResources, nil

}

// GetSelectedResources returns the specified and known resource.Resource's from a file system.
func (l *Storage) GetSelectedResources(sourceNames []string) ([]resource.Resource, error) {

	fetchedResources := make([]resource.Resource, 0)

	for _, name := range sourceNames {
		s := resource.Source(name)
		if detail, exists := l.resourceDetails[s]; exists {
			res, err := l.getResource(s, detail)
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

func (l *Storage) getResource(source resource.Source, detail ResourceDetail) (resource.Resource, error) {
	res, err := l.readFile(source, detail.format, detail.path)
	if err != nil {
		return resource.Resource{}, fmt.Errorf("error reading file: %v", err)
	}

	return res, nil
}

func (l *Storage) readFile(publisher resource.Source, format resource.Format, filename string) (resource.Resource, error) {
	content, err := l.readFileContent(filename)
	if err != nil {
		return resource.Resource{}, fmt.Errorf("error reading file content: %v", err)
	}

	res, err := resource.New(publisher, format, resource.Content(content))
	if err != nil {
		return resource.Resource{}, fmt.Errorf("error creating resource: %v", err)
	}

	return *res, nil
}

func (l *Storage) readFileContent(filename string) (string, error) {
	absPath := filepath.Join(l.basePath, filename)
	file, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}(file)

	return l.scanFile(file)
}

func (l *Storage) scanFile(file *os.File) (string, error) {
	scanner := bufio.NewScanner(file)
	var content string
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning file: %v", err)
	}

	return content, nil
}
