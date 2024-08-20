package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"updater/updater/model/feed"
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

// UpdateRSSFeed creates a new XML file with the content of the source.
func (s *Storage) UpdateRSSFeed(source feed.Source, content []byte) error {
	return s.updateSource(source, content, "xml")
}

// UpdateHTMLFeed creates a new HTML file with the content of the source.
func (s *Storage) UpdateHTMLFeed(source feed.Source, content []byte) error {
	return s.updateSource(source, content, "html")
}

func (s *Storage) updateSource(source feed.Source, content []byte, ext string) error {
	timestamp := time.Now().Format("20060102")
	filePath := fmt.Sprintf("%s_%s.%s", source, timestamp, ext)
	filePath = filepath.Join(s.basePath, filePath)

	err := os.WriteFile(filePath, content, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing resource to file: %v", err)
	}

	return nil
}
