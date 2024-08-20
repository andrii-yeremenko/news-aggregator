package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
	"updater/updater/model/feed"
)

func TestNew(t *testing.T) {
	t.Run("basePath provided", func(t *testing.T) {
		basePath := "testdata"
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Println("Error of removing directory")
			}
		}(basePath)

		storage, err := New(basePath)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if storage.basePath != basePath {
			t.Errorf("expected basePath %s, got %s", basePath, storage.basePath)
		}

		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			t.Errorf("expected directory %s to be created", basePath)
		}
	})

	t.Run("basePath not provided", func(t *testing.T) {
		defaultPath := "/resources"
		storage, _ := New("")

		if storage.basePath != defaultPath {
			t.Errorf("expected basePath %s, got %s", defaultPath, storage.basePath)
		}
	})

	t.Run("cannot create directory", func(t *testing.T) {
		restrictedPath := "restrictedDir"
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Println("Error of removing directory")
			}
		}(restrictedPath)

		err := os.Mkdir(restrictedPath, 0555)
		if err != nil {
			t.Fatalf("expected to create directory %s, got error %v", restrictedPath, err)
		}

		_ = os.Chdir(restrictedPath)

		_, err = New("")
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		_ = os.Chdir("..")

		_ = os.Chmod(restrictedPath, 0755)
		err = os.Remove(restrictedPath)
		if err != nil {
			t.Fatalf("expected to remove directory %s, got error %v", restrictedPath, err)
		}
	})
}

func TestStorage_UpdateRSSFeed(t *testing.T) {
	basePath := "testdata"
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Println("Error of removing directory")
		}
	}(basePath)
	storage, _ := New(basePath)

	source := feed.Source("test-source")
	content := []byte("<rss>test content</rss>")

	err := storage.UpdateRSSFeed(source, content)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	timestamp := time.Now().Format("20060102")
	expectedFilePath := filepath.Join(basePath, fmt.Sprintf("%s_%s.xml", source, timestamp))

	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Errorf("expected file %s to be created", expectedFilePath)
	}

	fileContent, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("expected to read file %s, got error %v", expectedFilePath, err)
	}

	if string(fileContent) != string(content) {
		t.Errorf("expected file content %s, got %s", content, fileContent)
	}
}

func TestStorage_UpdateHTMLFeed(t *testing.T) {
	basePath := "testdata"
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Println("Error of removing directory")
		}
	}(basePath)
	storage, _ := New(basePath)

	source := feed.Source("test-source")
	content := []byte("<html>test content</html>")

	err := storage.UpdateHTMLFeed(source, content)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	timestamp := time.Now().Format("20060102")
	expectedFilePath := filepath.Join(basePath, fmt.Sprintf("%s_%s.html", source, timestamp))

	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Errorf("expected file %s to be created", expectedFilePath)
	}

	fileContent, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("expected to read file %s, got error %v", expectedFilePath, err)
	}

	if string(fileContent) != string(content) {
		t.Errorf("expected file content %s, got %s", content, fileContent)
	}
}
