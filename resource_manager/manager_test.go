package resource_manager

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestStorage_GetAvailableSources(t *testing.T) {
	basePath, _ := os.Getwd()
	manager, _ := New(path.Join(basePath, "../storage"), path.Join(basePath, "resource_dictionary.json"))
	expectedSources := []string{"bbc-world", "usa-today", "nbc-news", "abc-news", "washington-times"}
	expectedCount := len(expectedSources)
	sources := manager.AvailableSources()

	for _, source := range expectedSources {
		if !strings.Contains(sources, source) {
			t.Errorf("expected source %s is missing", source)
		}
	}

	actualCount := strings.Count(sources, ",")
	if actualCount != expectedCount {
		t.Errorf("expected %d sources but got %d", expectedCount, actualCount)
	}
}

func TestStorage_GetAllResources(t *testing.T) {
	basePath, _ := os.Getwd()
	manager, _ := New(path.Join(basePath, "../storage"), path.Join(basePath, "resource_dictionary.json"))

	resources, err := manager.AllResources()
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if len(resources) != 5 {
		t.Fatalf("expected 5 resources but got %d", len(resources))
	}
}

func TestStorage_GetSelectedResources(t *testing.T) {
	basePath, _ := os.Getwd()
	manager, _ := New(path.Join(basePath, "../storage"), path.Join(basePath, "resource_dictionary.json"))
	selectedSources := []string{"nbc-news", "abc-news"}
	resources, err := manager.GetSelectedResources(selectedSources)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
	if len(resources) != 2 {
		t.Fatalf("expected 2 resources but got %d", len(resources))
	}

	for _, res := range resources {
		if res.Source() != "nbc-news" && res.Source() != "abc-news" {
			t.Errorf("expected source to be nbc-news or abc-news but got %s", res.Source())
		}
	}
}

func TestStorage_GetSelectedResources_InvalidSource(t *testing.T) {
	basePath, _ := os.Getwd()
	manager, _ := New(path.Join(basePath, "../storage"), path.Join(basePath, "resource_dictionary.json"))
	selectedSources := []string{"nbc-news", "abc-news", "invalid-source"}
	_, err := manager.GetSelectedResources(selectedSources)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}
