package web_server

// Manager is an interface that defines the methods for managing resources.
//
//go:generate mockgen -source=manager.go -destination=mocks/mock_manager.go -package=mocks
type Manager interface {
	UpdateAllSources() error
}
