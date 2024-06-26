package handler

import (
	"net/http"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/resource_manager"
)

// UpdateHandler is a handler for updating news sources from the internet.
type UpdateHandler struct {
	ResourceManager *resource_manager.ResourceManager
}

// NewUpdateHandler creates a new UpdateHandler instance.
func NewUpdateHandler(resourceManager *resource_manager.ResourceManager) *UpdateHandler {
	return &UpdateHandler{
		ResourceManager: resourceManager,
	}
}

// Handle handles the HTTP request and response for updating news sources.
func (h *UpdateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	sourceType := r.URL.Query().Get("source")

	if sourceType == "" {
		http.Error(w, "Source not specified", http.StatusBadRequest)
		return
	}

	if !h.ResourceManager.SourceIsSupported(resource.Source(sourceType)) {
		http.Error(w, "Source not supported", http.StatusBadRequest)
		return
	}

	err := h.ResourceManager.UpdateResource(resource.Source(sourceType))

	if err != nil {
		http.Error(w, "Failed to update source", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
