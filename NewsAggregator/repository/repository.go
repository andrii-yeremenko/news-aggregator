package repository

import (
	"NewsAggregator/aggregator/model/resource"
	"bufio"
	"fmt"
	"os"
)

// NewsRepository is a simple in-memory repository for resources.
type NewsRepository struct {
	resources []resource.Resource
}

// NewRepository creates a new NewsRepository.
func NewRepository() *NewsRepository {
	return &NewsRepository{}
}

// ReadFile reads a file and returns a resource.
func (r *NewsRepository) ReadFile(publisher resource.Source, format resource.Format, filename string) (resource.Resource, error) {
	file, err := os.Open(filename)
	if err != nil {
		return resource.Resource{}, fmt.Errorf("error opening file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file: %v", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	var content string
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return resource.Resource{}, fmt.Errorf("error scanning file: %v", err)
	}

	res, err := resource.NewResource(publisher, format, resource.Content(content))
	if err != nil {
		return resource.Resource{}, fmt.Errorf("error creating resource: %v", err)
	}

	return *res, nil
}
