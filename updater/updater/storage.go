package updater

import (
	"updater/model/feed"
)

// StorageInterface is an interface for storing feed data.
// MockStorage is a mock implementation of the StorageInterface
//
//go:generate mockgen -destination=mocks/mock_storage.go -package=mocks . StorageInterface
type StorageInterface interface {
	UpdateRSSFeed(source feed.Source, content []byte) error
	UpdateHTMLFeed(source feed.Source, content []byte) error
}
