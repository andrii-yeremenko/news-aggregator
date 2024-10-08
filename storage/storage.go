package storage

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"news-aggregator/aggregator/model/resource"
)

// Storage is a component enabling the retrieval and manipulation of known files from a file system.
type Storage struct {
	basePath string
}

// New creates a new Storage.
func New(basePath string) *Storage {

	if basePath == "" {
		basePath = "/resources"
	}

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		err := os.MkdirAll(basePath, os.ModePerm)
		if err != nil {
			fmt.Printf("error creating directory: %v\n", err)
		}
	}

	return &Storage{
		basePath: basePath,
	}
}

// FileExists checks if a file exists in the storage.
func (s *Storage) fileExists(filename string) bool {
	absPath := filepath.Join(s.basePath, filename)
	_, err := os.Stat(absPath)
	return err == nil
}

// AvailableSources returns all the available registered in storage sources.
func (s *Storage) AvailableSources() ([]resource.Source, error) {
	files, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	sourceMap := make(map[string]bool)
	var sources []resource.Source

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		parts := strings.Split(fileName, "_")
		if len(parts) > 1 {
			sourceName := strings.Join(parts[:len(parts)-1], "_")
			if !sourceMap[sourceName] {
				sourceMap[sourceName] = true
				sources = append(sources, resource.Source(sourceName))
			}
		}
	}

	return sources, nil
}

// ReadSource reads the content of files starting with the source name.
func (s *Storage) ReadSource(source resource.Source) ([]string, error) {
	var contents []string
	prefix := string(source)
	files, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		if strings.HasPrefix(fileName, prefix) {
			absPath := filepath.Join(s.basePath, fileName)
			content, err := s.readFileContents(absPath)
			if err != nil {
				return nil, err
			}
			contents = append(contents, content)
		}
	}

	if len(contents) == 0 {
		return nil, fmt.Errorf("source %s is unknown", source)
	}

	return contents, nil
}

// UpdateXMLSource creates a new xml file with the content of the source.
func (s *Storage) UpdateXMLSource(source resource.Source, content []byte) error {
	return s.updateSource(source, content, "xml")
}

// UpdateJSONSource creates a new json file with the content of the source.
func (s *Storage) UpdateJSONSource(source resource.Source, content []byte) error {
	return s.updateSource(source, content, "json")
}

// UpdateHTMLSource creates a new html file with the content of the source.
func (s *Storage) UpdateHTMLSource(source resource.Source, content []byte) error {
	return s.updateSource(source, content, "html")
}

func (s *Storage) updateSource(source resource.Source, content []byte, ext string) error {
	timestamp := time.Now().Format("20060102")
	filePath := fmt.Sprintf("%s_%s.%s", source, timestamp, ext)
	filePath = filepath.Join(s.basePath, filePath)

	err := os.WriteFile(filePath, content, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing resource to file: %v", err)
	}

	return nil
}

func (s *Storage) readFileContents(absPath string) (string, error) {
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

	scanner := bufio.NewScanner(file)
	var content strings.Builder
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning file: %v", err)
	}

	return content.String(), nil
}
