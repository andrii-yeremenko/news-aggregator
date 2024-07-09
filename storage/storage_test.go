package storage

import (
	"github.com/stretchr/testify/assert"
	"news-aggregator/aggregator/model/resource"
	"os"
	"path/filepath"
	"testing"
)

func TestStorage_FileExists(t *testing.T) {
	testDir := "testdata"
	storage := New(testDir)

	// Test existing file
	exists := storage.fileExists("test.txt")
	assert.True(t, exists, "expected file to exist")

	// Test non-existent file
	exists = storage.fileExists("nonexistent.txt")
	assert.False(t, exists, "expected file to not exist")
}

func TestStorage_ReadSource(t *testing.T) {
	testDir := "testdata"
	storage := New(testDir)

	testFilePath := filepath.Join(testDir, "resources/abc-news.xml")
	err := os.MkdirAll(filepath.Dir(testFilePath), 0755)
	assert.NoError(t, err)
	file, err := os.Create(testFilePath)
	assert.NoError(t, err)
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("error removing test file: %v", err)
		}
	}(filepath.Dir(testFilePath))
	_, _ = file.WriteString("test content\nanother line\n")
	err = file.Close()
	if err != nil {
		t.Fatalf("error closing test file: %v", err)
	}

	expectedContent := "test content\nanother line\n"
	readContent, err := storage.ReadSource("abc-news")
	assert.NoError(t, err, "expected no error reading file")
	assert.Equal(t, expectedContent, readContent, "expected file content to match")

	_, err = storage.ReadSource("unknown-source")
	assert.Error(t, err, "expected error reading unknown source")
	assert.Equal(t, "source unknown-source is unknown", err.Error(), "expected specific error message")
}

func TestStorage_AvailableSources(t *testing.T) {
	testDir := "testdata"
	storage := New(testDir)

	for _, path := range storage.resourcesPath {
		testFilePath := filepath.Join(testDir, path)
		err := os.MkdirAll(filepath.Dir(testFilePath), 0755)
		assert.NoError(t, err)
		file, err := os.Create(testFilePath)
		assert.NoError(t, err)
		err = file.Close()
		if err != nil {
			t.Fatalf("error closing test file: %v", err)
		}
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				t.Fatalf("error removing test file: %v", err)
			}
		}(filepath.Dir(testFilePath))
	}

	expectedSources := []resource.Source{"nbc-news", "abc-news", "washington-times", "bbc-world", "usa-today"}
	availableSources := storage.AvailableSources()
	assert.ElementsMatch(t, expectedSources, availableSources, "expected available sources to match")
}
