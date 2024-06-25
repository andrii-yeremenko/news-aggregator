package handler

import (
	"fmt"
	"net/http"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/resource_manager"
)

// UpdateSourcesHandler is a handler for updating news sources from the internet.
type UpdateSourcesHandler struct {
	ResourceManager *resource_manager.ResourceManager
}

// NewUpdateSourcesHandler creates a new UpdateSourcesHandler instance.
func NewUpdateSourcesHandler(resourceManager *resource_manager.ResourceManager) *UpdateSourcesHandler {
	return &UpdateSourcesHandler{
		ResourceManager: resourceManager,
	}
}

// Handle handles the HTTP request and response for updating news sources.
func (h *UpdateSourcesHandler) Handle(w http.ResponseWriter, r *http.Request) {
	sourceType := r.URL.Query().Get("source")

	if sourceType == "" {
		h.respondWithError(w, fmt.Errorf("update source not specified"), http.StatusBadRequest)
		return
	}

	if !h.ResourceManager.SourceIsSupported(resource.Source(sourceType)) {
		h.respondWithError(w, fmt.Errorf("unsupported source: %s", sourceType), http.StatusBadRequest)
		return
	}

	err := h.ResourceManager.UpdateResource(resource.Source(sourceType))

	if err != nil {
		h.respondWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// respondWithError responds with an error.
func (h *UpdateSourcesHandler) respondWithError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := NewErrorResponse(err)
	_, _ = w.Write(errorResponse.getJSON())
}
