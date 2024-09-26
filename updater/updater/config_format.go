package updater

// FeedJSON is a struct that represents how feed is stored in the JSON config file.
type FeedJSON struct {
	Source string `json:"source"`
	Format string `json:"format"`
	Link   string `json:"link"`
}
