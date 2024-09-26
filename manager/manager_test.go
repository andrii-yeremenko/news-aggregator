package manager_test

import (
	"encoding/json"
	"news-aggregator/manager"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"news-aggregator/aggregator/model/resource"
)

const (
	testStoragePath    = "./testdata"
	testFeedDictionary = "./testdata/feeds.json"
)

func TestNewResourceManager(t *testing.T) {
	rm, err := manager.New(testStoragePath, testFeedDictionary)
	assert.NoError(t, err)
	assert.NotNil(t, rm)
}

func TestRegisterSource(t *testing.T) {
	rm, err := manager.New(testStoragePath, testFeedDictionary)
	assert.NoError(t, err)

	err = rm.RegisterSource("source", "http://source.com/source", resource.RSS)
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
	assert.Equal(t, 3, len(feeds))
}

func TestUpdateSource(t *testing.T) {
	rm, err := manager.New(testStoragePath, testFeedDictionary)
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
	rm, err := manager.New(testStoragePath, testFeedDictionary)
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
	rm, err := manager.New(testStoragePath, testFeedDictionary)
	assert.NoError(t, err)

	assert.True(t, rm.IsSourceSupported("supported_source"))
}

func TestAvailableSources(t *testing.T) {
	rm, err := manager.New(testStoragePath, testFeedDictionary)
	assert.NoError(t, err)

	expectedSources := "supported,test,"
	assert.Equal(t, expectedSources, rm.AvailableSources())
}

func TestGetAllResources(t *testing.T) {
	rm, err := manager.New(testStoragePath, testFeedDictionary)
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
	rm, err := manager.New(testStoragePath, testFeedDictionary)
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

func TestGetAvailableFeeds_SingleFeed(t *testing.T) {

	file := "./testdata/empty.json"
	_ = os.Remove(file)
	_, _ = os.Create(file)

	rm, err := manager.New(testStoragePath, "./testdata/empty.json")
	assert.NoError(t, err)

	_ = rm.RegisterSource("test", "http://test.com/test", resource.HTML)

	feeds := rm.AvailableFeeds()

	assert.Equal(t, "test", feeds)
}

func TestGetAvailableFeeds_MultipleFeeds(t *testing.T) {

	file := "./testdata/empty.json"
	_ = os.Remove(file)
	_, _ = os.Create(file)

	rm, err := manager.New(testStoragePath, "./testdata/empty.json")
	assert.NoError(t, err)

	_ = rm.RegisterSource("test", "http://test.com/test", resource.HTML)
	_ = rm.RegisterSource("test2", "http://test.com/test2", resource.HTML)

	feeds := rm.AvailableFeeds()

	ok := feeds == "test,test2" || feeds == "test2,test"

	assert.True(t, ok)
}

func TestGetAvailableFeeds_NoFeeds(t *testing.T) {

	file := "./testdata/empty.json"
	_ = os.Remove(file)
	_, _ = os.Create(file)

	rm, err := manager.New(testStoragePath, "./testdata/empty.json")

	_ = rm.DeleteSource("supported_source")
	_ = rm.DeleteSource("test")

	assert.NoError(t, err)

	assert.Equal(t, "no available feeds", rm.AvailableFeeds())
}
