package storage

import (
	"bufio"
	"fmt"
	"news-aggregator/aggregator/model/resource"
	"os"
	"path/filepath"
)

// Storage is a component enabling the retrieval and manipulation of known files from a file system.
type Storage struct {
	basePath      string
	resourcesPath map[resource.Source]string
}

// New creates a new Storage.
func New(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
		resourcesPath: map[resource.Source]string{
			"nbc-news":         "resources/nbc-news.json",
			"abc-news":         "resources/abc-news.xml",
			"washington-times": "resources/washington-times.xml",
			"bbc-world":        "resources/bbc-world.xml",
			"usa-today":        "resources/usa-today-world-news.html",
		},
	}
}

// FileExists checks if a file exists in the storage.
func (s *Storage) fileExists(filename string) bool {
	absPath := filepath.Join(s.basePath, filename)
	_, err := os.Stat(absPath)

	if err != nil {
		return false
	}

	return true
}

// ReadSource reads the content of a file in the storage.
func (s *Storage) ReadSource(source resource.Source) (string, error) {

	if s.resourcesPath[source] == "" {
		return "", fmt.Errorf("source %s is unknown", source)
	}

	absPath := filepath.Join(s.basePath, s.resourcesPath[source])
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

	return s.scanFile(file)
}

// AvailableSources returns all the available registered is storage sources.
func (s *Storage) AvailableSources() []resource.Source {
	sources := make([]resource.Source, 0, len(s.resourcesPath))
	for source := range s.resourcesPath {
		if s.fileExists(s.resourcesPath[source]) {
			sources = append(sources, source)
		}
	}
	return sources
}

func (s *Storage) scanFile(file *os.File) (string, error) {
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
