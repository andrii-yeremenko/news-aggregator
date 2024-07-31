package updater

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"updater/updater/mocks"

	"github.com/golang/mock/gomock"
	"updater/model/feed"
)

func TestNewUpdater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name            string
		feedsConfigPath string
		mockSetup       func(*mocks.MockStorageInterface)
		expectedError   error
	}{
		{
			name:            "valid path",
			feedsConfigPath: "testdata/feeds.json",
			mockSetup:       func(m *mocks.MockStorageInterface) {},
			expectedError:   nil,
		},
		{
			name:            "empty path",
			feedsConfigPath: "",
			mockSetup:       func(m *mocks.MockStorageInterface) {},
			expectedError:   fmt.Errorf("feeds config path not provided"),
		},
		{
			name:            "invalid JSON in path",
			feedsConfigPath: "testdata/invalid_feeds.json",
			mockSetup:       func(m *mocks.MockStorageInterface) {},
			expectedError:   fmt.Errorf("can't parse JSON: invalid character 'i' looking for beginning of value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := mocks.NewMockStorageInterface(ctrl)
			updater, err := New(tt.feedsConfigPath, storageMock)
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("expected error: %v, got nil", tt.expectedError)
			}
			if err == nil && updater.feedsConfigPath != tt.feedsConfigPath {
				t.Errorf("expected feedsConfigPath: %v, got: %v", tt.feedsConfigPath, updater.feedsConfigPath)
			}
		})
	}
}

func TestUpdateFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		feedSource    string
		mockSetup     func(*mocks.MockStorageInterface)
		expectedError error
	}{
		{
			name:          "unsupported format",
			feedSource:    "unsupportedFeedSource",
			mockSetup:     func(m *mocks.MockStorageInterface) {},
			expectedError: fmt.Errorf("feed source not found: unsupportedFeedSource"),
		},
		{
			name:          "feed source not found",
			feedSource:    "nonexistentFeedSource",
			mockSetup:     func(m *mocks.MockStorageInterface) {},
			expectedError: fmt.Errorf("feed source not found: nonexistentFeedSource"),
		},
	}

	testFeedABC, _ := feed.New("abc-news", feed.RSS, "https://feeds.abcnews.com/abcnews/internationalheadlines")
	testFeedWT, _ := feed.New("washington-times", feed.RSS, "https://www.washingtontimes.com/rss/headlines/news/world/")
	feeds := []*feed.Feed{
		testFeedABC,
		testFeedWT,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := mocks.NewMockStorageInterface(ctrl)
			updater := Updater{
				feeds:   feeds,
				storage: storageMock,
			}

			tt.mockSetup(storageMock)
			http.DefaultTransport = &mockTransport{
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("feed content")),
				},
				err: nil,
			}

			err := updater.UpdateFeed(tt.feedSource)
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("expected error: %v, got nil", tt.expectedError)
			}
		})
	}
}

func TestAvailableFeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testFeedABC, _ := feed.New("abc-news", feed.RSS, "https://feeds.abcnews.com/abcnews/internationalheadlines")
	testFeedWT, _ := feed.New("washington-times", feed.RSS, "https://www.washingtontimes.com/rss/headlines/news/world/")

	tests := []struct {
		name          string
		feeds         []*feed.Feed
		expectedFeeds []string
		expectedError error
	}{
		{
			name: "available feeds",
			feeds: []*feed.Feed{
				testFeedABC,
				testFeedWT,
			},
			expectedFeeds: []string{"abc-news", "washington-times"},
			expectedError: nil,
		},
		{
			name:          "no feeds available",
			feeds:         []*feed.Feed{},
			expectedFeeds: []string{},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := mocks.NewMockStorageInterface(ctrl)
			updater := Updater{
				feeds:   tt.feeds,
				storage: storageMock,
			}

			feeds, err := updater.AvailableFeeds()
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("expected error: %v, got nil", tt.expectedError)
			}
			if !equal(feeds, tt.expectedFeeds) {
				t.Errorf("expected feeds: %v, got: %v", tt.expectedFeeds, feeds)
			}
		})
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type mockTransport struct {
	resp *http.Response
	err  error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}
