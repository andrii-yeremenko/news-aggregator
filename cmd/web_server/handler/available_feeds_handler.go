package handler

import (
	"encoding/json"
	"net/http"
)

type AvailableFeedsHandler struct {
	manager ResourceManager
}

func NewAvailableFeedsHandler(manager ResourceManager) *AvailableFeedsHandler {
	return &AvailableFeedsHandler{
		manager: manager,
	}
}

func (ch *AvailableFeedsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ch.GetSources(w)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ch *AvailableFeedsHandler) GetSources(w http.ResponseWriter) {
	feeds := ch.manager.AvailableFeeds()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(feeds)
	if err != nil {
		http.Error(w, "Failed to encode feeds", http.StatusInternalServerError)
		return
	}
}
