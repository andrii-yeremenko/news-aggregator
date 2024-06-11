package storage_test

import (
	"NewsAggregator/storage"
	"testing"
)

func TestStorage_GetAvailableSources(t *testing.T) {
	store := storage.New("")
	expectedSources := "nbc-news, abc-news, washington-times, bbc-world, usa-today, "
	sources := store.GetAvailableSources()
	if sources != expectedSources {
		t.Errorf("expected %s but got %s", expectedSources, sources)
	}
}

func TestStorage_GetAllResources(t *testing.T) {
	store := storage.New("")

	resources, err := store.GetAllResources()
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if len(resources) != 5 {
		t.Fatalf("expected 5 resources but got %d", len(resources))
	}
}

func TestStorage_GetSelectedResources(t *testing.T) {
	store := storage.New("")

	selectedSources := []string{"nbc-news", "abc-news"}
	resources, err := store.GetSelectedResources(selectedSources)
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
