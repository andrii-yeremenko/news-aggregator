package updater

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"updater/updater/mocks"
	feed2 "updater/updater/model/feed"

	"github.com/golang/mock/gomock"
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
		serverHandler http.HandlerFunc
		mockSetup     func(*mocks.MockStorageInterface)
		expectedError error
	}{
		{
			name:       "unsupported format",
			feedSource: "unsupportedFeedSource",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			mockSetup:     func(m *mocks.MockStorageInterface) {},
			expectedError: fmt.Errorf("feed source not found: unsupportedFeedSource"),
		},
		{
			name:       "feed source not found",
			feedSource: "nonexistentFeedSource",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			mockSetup:     func(m *mocks.MockStorageInterface) {},
			expectedError: fmt.Errorf("feed source not found: nonexistentFeedSource"),
		},
		{
			name:       "successful feed update",
			feedSource: "abc-news",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("feed content"))
			},
			mockSetup: func(m *mocks.MockStorageInterface) {
				m.EXPECT().UpdateRSSFeed(feed2.Source("abc-news"), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.serverHandler)
			defer server.Close()

			testFeedABC, _ := feed2.New("abc-news", feed2.RSS, feed2.Link(server.URL))
			testFeedWT, _ := feed2.New("washington-times", feed2.RSS, feed2.Link(server.URL))
			feeds := []*feed2.Feed{testFeedABC, testFeedWT}

			storageMock := mocks.NewMockStorageInterface(ctrl)
			updater := Updater{
				feeds:   feeds,
				storage: storageMock,
			}

			tt.mockSetup(storageMock)

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
