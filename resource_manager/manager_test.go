package resource_manager_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/resource_manager"
)

const (
	testStoragePath    = "./testdata"
	testFeedDictionary = "./testdata/feeds.json"
)

func TestNewResourceManager(t *testing.T) {
	rm, err := resource_manager.New(testStoragePath, testFeedDictionary)
	assert.NoError(t, err)
	assert.NotNil(t, rm)
}

func TestRegisterSource(t *testing.T) {
	rm, err := resource_manager.New(testStoragePath, testFeedDictionary)
	assert.NoError(t, err)

	err = rm.RegisterSource("source", "http://source.com/source", resource.RSS)
	assert.NoError(t, err)

	data, err := ioutil.ReadFile(testFeedDictionary)
	assert.NoError(t, err)

	var feeds []struct {
		Source string `json:"source"`
		Format string `json:"format"`
		Link   string `json:"link"`
	}
	err = json.Unmarshal(data, &feeds)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(feeds))
}

func TestUpdateSource(t *testing.T) {
	rm, err := resource_manager.New(testStoragePath, testFeedDictionary)
	assert.NoError(t, err)

	err = rm.UpdateSource("source", "http://source.com/updated", resource.HTML)
	assert.NoError(t, err)

	data, err := os.ReadFile(testFeedDictionary)
	assert.NoError(t, err)

	var feeds []struct {
		Source string `json:"source"`
		Format string `json:"format"`
		Link   string `json:"link"`
	}
	err = json.Unmarshal(data, &feeds)
	assert.NoError(t, err)

	for _, feed := range feeds {
		if feed.Source == "source" {
			assert.Equal(t, "http://source.com/updated", feed.Link)
			assert.Equal(t, "HTML", feed.Format)
		}
	}
}

func TestDeleteSource(t *testing.T) {
	rm, err := resource_manager.New(testStoragePath, testFeedDictionary)
	assert.NoError(t, err)

	err = rm.DeleteSource("source")
	assert.NoError(t, err)

	data, err := os.ReadFile(testFeedDictionary)
	assert.NoError(t, err)

	var feeds []struct {
		Source string `json:"source"`
		Format string `json:"format"`
		Link   string `json:"link"`
	}
	err = json.Unmarshal(data, &feeds)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(feeds))
	assert.NotEqual(t, "source", feeds[0].Source)
}

func TestSourceIsSupported(t *testing.T) {
	rm, err := resource_manager.New(testStoragePath, testFeedDictionary)
	assert.NoError(t, err)

	assert.True(t, rm.SourceIsSupported("supported_source"))
}

func TestAvailableSources(t *testing.T) {
	rm, err := resource_manager.New(testStoragePath, testFeedDictionary)
	assert.NoError(t, err)

	expectedSources := "supported,test,"
	assert.Equal(t, expectedSources, rm.AvailableSources())
}

func TestGetAllResources(t *testing.T) {
	rm, err := resource_manager.New(testStoragePath, testFeedDictionary)
	assert.NoError(t, err)

	resources, err := rm.GetAllResources()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(resources))

	for _, r := range resources {
		assert.NotEmpty(t, r.Content())
		assert.NotEmpty(t, r.Source())
		assert.NotEmpty(t, r.Format())
	}
}

func TestGetSelectedResources(t *testing.T) {
	rm, err := resource_manager.New(testStoragePath, testFeedDictionary)
	assert.NoError(t, err)

	resources, err := rm.GetSelectedResources([]string{"supported_source"})
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resources))

	for _, r := range resources {
		assert.Equal(t, "supported_source", string(r.Source()))
		assert.Equal(t, "Test text\n", string(r.Content()))
		assert.Equal(t, resource.HTML, int(r.Format()))
	}
}
