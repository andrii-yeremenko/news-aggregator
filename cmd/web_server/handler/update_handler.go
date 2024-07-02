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
	sourceType := r.URL.Query().Get("source")

	if sourceType == "" {
		http.Error(w, "Source not specified", http.StatusBadRequest)
		return
	}

	//rename to IsSourceSupported
	if !h.ResourceManager.IsSourceSupported(resource.Source(sourceType)) {
		http.Error(w, "Source not supported", http.StatusBadRequest)
		return
	}
	//resource.Source(sourceType)
	err := h.ResourceManager.UpdateResource(resource.Source(sourceType))

	if err != nil {
		http.Error(w, "Failed to update source", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
