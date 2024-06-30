package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"news-aggregator/aggregator/model/resource"
)

func createTestFile(t *testing.T, dir, name, content string) {
	err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
	if err != nil {
		t.Fatalf("error creating test file: %v", err)
	}
}

func TestFileExists(t *testing.T) {
	dir := "testdata"
	s := New(dir)
	testFile := "testfile.txt"
	createTestFile(t, dir+"/resources", testFile, "content")

	defer os.Remove(filepath.Join(dir, "resources", testFile))

	exists := s.fileExists(testFile)
	if !exists {
		t.Errorf("expected file %s to exist", testFile)
	}

	exists = s.fileExists("nonexistentfile.txt")
	if exists {
		t.Errorf("expected file %s to not exist", "nonexistentfile.txt")
	}
}

func TestAvailableSources(t *testing.T) {
	dir := "testdata"
	storage := New(dir)

	createTestFile(t, dir+"/resources", "source1_20210101.xml", "content")
	createTestFile(t, dir+"/resources", "source2_20210101.json", "content")
	createTestFile(t, dir+"/resources", "source1_20210102.html", "content")

	defer os.Remove(filepath.Join(dir, "resources", "source1_20210101.xml"))
	defer os.Remove(filepath.Join(dir, "resources", "source2_20210101.json"))
	defer os.Remove(filepath.Join(dir, "resources", "source1_20210102.html"))

	sources, err := storage.AvailableSources()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedSources := map[resource.Source]bool{
		"source1": true,
		"source2": true,
	}

	if len(sources) != len(expectedSources) {
		t.Fatalf("expected %d sources, got %d", len(expectedSources), len(sources))
	}

	for _, source := range sources {
		if !expectedSources[source] {
			t.Errorf("unexpected source: %s", source)
		}
	}
}

func TestReadSource(t *testing.T) {
	dir := "testdata"
	storage := New(dir)

	createTestFile(t, dir+"/resources", "source1_20210101.xml", "content1")
	createTestFile(t, dir+"/resources", "source1_20210102.html", "content2")

	defer os.Remove(filepath.Join(dir, "resources", "source1_20210101.xml"))
	defer os.Remove(filepath.Join(dir, "resources", "source1_20210102.html"))

	contents, err := storage.ReadSource("source1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedContents := []string{"content1\n", "content2\n"}
	if len(contents) != len(expectedContents) {
		t.Fatalf("expected %d contents, got %d", len(expectedContents), len(contents))
	}

	for i, content := range contents {
		if content != expectedContents[i] {
			t.Errorf("expected content %q, got %q", expectedContents[i], content)
		}
	}

	_, err = storage.ReadSource("unknown_source")
	if err == nil {
		t.Fatalf("expected error for unknown source")
	}
}

func TestUpdateXMLSource(t *testing.T) {
	dir := "testdata"
	storage := New(dir)
	source := resource.Source("source1")
	content := []byte("<xml>content</xml>")

	err := storage.UpdateXMLSource(source, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := ioutil.ReadDir(filepath.Join(dir, "resources"))
	if err != nil {
		t.Fatalf("error reading directory: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	expectedFileName := string(source) + "_" + time.Now().Format("20060102") + ".xml"
	if files[0].Name() != expectedFileName {
		t.Errorf("expected file name %q, got %q", expectedFileName, files[0].Name())
	}

	fileContent, err := ioutil.ReadFile(filepath.Join(dir, "resources", files[0].Name()))
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}

	defer os.Remove(filepath.Join(dir, "resources", files[0].Name()))

	if string(fileContent) != string(content) {
		t.Errorf("expected file content %q, got %q", content, fileContent)
	}
}

func TestUpdateJSONSource(t *testing.T) {
	dir := "testdata"
	storage := New(dir)
	source := resource.Source("source2")
	content := []byte(`{"key": "value"}`)

	err := storage.UpdateJSONSource(source, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := ioutil.ReadDir(filepath.Join(dir, "resources"))
	if err != nil {
		t.Fatalf("error reading directory: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	expectedFileName := string(source) + "_" + time.Now().Format("20060102") + ".json"
	if files[0].Name() != expectedFileName {
		t.Errorf("expected file name %q, got %q", expectedFileName, files[0].Name())
	}

	fileContent, err := ioutil.ReadFile(filepath.Join(dir, "resources", files[0].Name()))
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}

	defer os.Remove(filepath.Join(dir, "resources", files[0].Name()))

	if string(fileContent) != string(content) {
		t.Errorf("expected file content %q, got %q", content, fileContent)
	}
}

func TestUpdateHTMLSource(t *testing.T) {
	dir := "testdata"
	storage := New(dir)
	source := resource.Source("source3")
	content := []byte("<html>content</html>")

	err := storage.UpdateHTMLSource(source, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := ioutil.ReadDir(filepath.Join(dir, "resources"))
	if err != nil {
		t.Fatalf("error reading directory: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	expectedFileName := string(source) + "_" + time.Now().Format("20060102") + ".html"
	if files[0].Name() != expectedFileName {
		t.Errorf("expected file name %q, got %q", expectedFileName, files[0].Name())
	}

	fileContent, err := ioutil.ReadFile(filepath.Join(dir, "resources", files[0].Name()))
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}

	defer os.Remove(filepath.Join(dir, "resources", files[0].Name()))

	if string(fileContent) != string(content) {
		t.Errorf("expected file content %q, got %q", content, fileContent)
	}
}
