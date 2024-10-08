package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"updater/updater/model/feed"
)

type Updater struct {
	feedsConfigPath string
	storage         StorageInterface
	feeds           []*feed.Feed
}

func New(feedsConfigPath string, storage StorageInterface) (Updater, error) {

	if feedsConfigPath == "" {
		return Updater{}, fmt.Errorf("feeds config path not provided")
	}

	feeds, err := loadFeedsInfo(feedsConfigPath)

	return Updater{
		feedsConfigPath: feedsConfigPath,
		storage:         storage,
		feeds:           feeds,
	}, err
}

// UpdateAllFeeds updates all feeds. If some feed fails to update, it will continue with the next one.
func (u *Updater) UpdateAllFeeds() []error {

	var errors []error

	for _, f := range u.feeds {
		err := u.UpdateFeed(string(f.Source()))
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// UpdateFeed updates a specific feeds.
func (u *Updater) UpdateFeed(feedSource string) error {

	source := feed.Source(feedSource)
	if source == "" {
		return fmt.Errorf("feed source not provided")
	}

	var targetFeed *feed.Feed
	for _, f := range u.feeds {
		if f.Source() == source {
			targetFeed = f
			break
		}
	}

	if targetFeed == nil {
		return fmt.Errorf("feed source not found: %s", feedSource)
	}

	resp, err := http.Get(string(targetFeed.Link()))
	if err != nil {
		return fmt.Errorf("error fetching resource from link: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching resource from link: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading resource content: %v", err)
	}

	switch targetFeed.Format() {
	case feed.RSS:
		return u.storage.UpdateRSSFeed(targetFeed.Source(), body)
	case feed.HTML:
		return u.storage.UpdateHTMLFeed(targetFeed.Source(), body)
	default:
		return fmt.Errorf("unsupported format")
	}
}

// AvailableFeeds returns a list of available feeds.
func (u *Updater) AvailableFeeds() ([]string, error) {
	var sources []string
	for _, f := range u.feeds {
		sources = append(sources, string(f.Source()))
	}
	return sources, nil
}

func loadFeedsInfo(path string) ([]*feed.Feed, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error of closing file")
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("can't read file: %v", err)
	}

	var feedsJSON []FeedJSON
	err = json.Unmarshal(data, &feedsJSON)
	if err != nil {
		return nil, fmt.Errorf("can't parse JSON: %v", err)
	}

	var feeds []*feed.Feed
	for _, feedJSON := range feedsJSON {
		format, err := feed.ParseFormat(feedJSON.Format)
		if err != nil {
			return nil, fmt.Errorf("unknown format: %v", err)
		}

		newFeed, err := feed.New(feed.Source(feedJSON.Source), format, feed.Link(feedJSON.Link))
		if err != nil {
			return nil, fmt.Errorf("can't create feed: %v", err)
		}

		feeds = append(feeds, newFeed)
	}

	return feeds, nil
}
