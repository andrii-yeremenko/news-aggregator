package storage_test

import (
	"NewsAggregator/aggregator/model/resource"
	storage2 "NewsAggregator/storage"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	testdataDir := "testdata"
	err := os.MkdirAll(testdataDir, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create testdata directory: %v", err)
	}

	filename := filepath.Join(testdataDir, "sample.txt")
	content := "This is a test file.\nWith multiple lines of text.\n"
	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create sample file: %v", err)
	}

	publisher := resource.Source("Test Publisher")
	format := resource.Format("Test Format")

	storage := storage2.New()

	res, err := storage.ReadFile(publisher, format, filename)
	if err != nil {
		t.Fatalf("ReadFile returned an error: %v", err)
	}

	expectedContent := content
	if string(res.Content()) != expectedContent {
		t.Errorf("unexpected content: got %v, want %v", res.Content(), expectedContent)
	}

	_, err = storage.ReadFile(publisher, format, filepath.Join(testdataDir, "nonexistent.txt"))
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}

	_, err = storage.ReadFile("", "", filename)
	if err == nil {
		t.Fatal("expected error during resource creation, got nil")
	}

	errorContentFile := filepath.Join(testdataDir, "error.txt")
	err = os.WriteFile(errorContentFile, []byte("error"), 0644)
	if err != nil {
		t.Fatalf("failed to create error content file: %v", err)
	}

	err = os.Remove(filename)
	err = os.Remove(errorContentFile)
	err = os.Remove(testdataDir)
}
