package handler

import (
	"net/http"
	"news-aggregator/aggregator/model/resource"
)

// UpdateHandler is a handler for updating news sources from the internet.
type UpdateHandler struct {
	ResourceManager ResourceManager
}

// NewUpdateHandler creates a new UpdateHandler instance.
func NewUpdateHandler(resourceManager ResourceManager) *UpdateHandler {
	return &UpdateHandler{
		ResourceManager: resourceManager,
	}
}

// Handle handles the HTTP request and response for updating news sources.
func (h *UpdateHandler) Handle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sourceType := r.URL.Query().Get("source")

	if sourceType == "" {
		http.Error(w, "Source not specified", http.StatusBadRequest)
		return
	}

	source := resource.Source(sourceType)

	if !h.ResourceManager.IsSourceSupported(source) {
		http.Error(w, "Source not supported", http.StatusBadRequest)
		return
	}

	err := h.ResourceManager.UpdateResource(source)

	if err != nil {
		http.Error(w, "Failed to update source", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
