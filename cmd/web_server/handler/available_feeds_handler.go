package handler

import (
	"encoding/json"
	"net/http"
)

// AvailableFeedsHandler is a handler that returns available feeds.
type AvailableFeedsHandler struct {
	manager ResourceManager
}

// NewAvailableFeedsHandler creates a new AvailableFeedsHandler.
func NewAvailableFeedsHandler(manager ResourceManager) *AvailableFeedsHandler {
	return &AvailableFeedsHandler{
		manager: manager,
	}
}

// Handle handles the request.
func (ch *AvailableFeedsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ch.GetSources(w)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetSources returns available feeds.
func (ch *AvailableFeedsHandler) GetSources(w http.ResponseWriter) {
	feeds := ch.manager.AvailableFeeds()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(feeds)
	if err != nil {
		http.Error(w, "Failed to encode feeds", http.StatusInternalServerError)
		return
	}
}
