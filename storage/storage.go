package storage

import (
	"NewsAggregator/aggregator/model/resource"
	"bufio"
	"fmt"
	"os"
)

// NewsStorage is a simple in-memory repository for resources.
type NewsStorage struct {
	resources []resource.Resource
}

// New creates a new NewsStorage.
func New() *NewsStorage {
	return &NewsStorage{}
}

// ReadFile reads a file and returns a resource.
func (r *NewsStorage) ReadFile(publisher resource.Source, format resource.Format, filename string) (resource.Resource, error) {
	content, err := r.readFileContent(filename)
	if err != nil {
		return resource.Resource{}, fmt.Errorf("error reading file content: %v", err)
	}

	res, err := resource.New(publisher, format, resource.Content(content))
	if err != nil {
		return resource.Resource{}, fmt.Errorf("error creating resource: %v", err)
	}

	return *res, nil
}

// readFileContent reads the content of the given file.
func (r *NewsStorage) readFileContent(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}(file)

	return r.scanFile(file)
}

// scanFile scans the content of the provided file.
func (r *NewsStorage) scanFile(file *os.File) (string, error) {
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
