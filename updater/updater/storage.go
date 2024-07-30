package updater

import (
	"updater/model/feed"
)

// StorageInterface is an interface for storing feed data.
type StorageInterface interface {
	UpdateRSSFeed(source feed.Source, content []byte) error
	UpdateHTMLFeed(source feed.Source, content []byte) error
}
