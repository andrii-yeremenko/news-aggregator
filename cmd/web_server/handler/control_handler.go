package handler

import (
	"encoding/json"
	"net/http"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/resource_manager"
)

// ControlHandler handles requests for managing news sources.
type ControlHandler struct {
	manager *resource_manager.ResourceManager
}

// NewControlHandler creates a new ControlHandler instance.
func NewControlHandler(manager *resource_manager.ResourceManager) *ControlHandler {
	return &ControlHandler{
		manager: manager,
	}
}

// Handle routes the request based on the HTTP method.
func (ch *ControlHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ch.GetSources(w, r)
	case http.MethodPost:
		ch.AddSource(w, r)
	case http.MethodPut:
		ch.UpdateSource(w, r)
	case http.MethodDelete:
		ch.DeleteSource(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetSources handles GET /sources to retrieve all sources.
func (ch *ControlHandler) GetSources(w http.ResponseWriter, r *http.Request) {
	sources := ch.manager.AvailableSources()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(sources)
	if err != nil {
		http.Error(w, "Failed to encode sources", http.StatusInternalServerError)
		return
	}
}

// AddSource handles POST /sources to add a new source.
func (ch *ControlHandler) AddSource(w http.ResponseWriter, r *http.Request) {
	var source struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		Format string `json:"format"`
	}
	if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	format, err := resource.ParseFormat(source.Format)

	if err != nil {
		ch.respondWithError(w, err, http.StatusBadRequest)
		return
	}

	err = ch.manager.RegisterSource(resource.Source(source.Name), source.URL, format)
	if err != nil {
		ch.respondWithError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateSource handles PUT /sources to update an existing source.
func (ch *ControlHandler) UpdateSource(w http.ResponseWriter, r *http.Request) {
	var source struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		Format string `json:"format"`
	}

	if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	format, err := resource.ParseFormat(source.Format)

	if err != nil {
		ch.respondWithError(w, err, http.StatusBadRequest)
		return
	}

	err = ch.manager.UpdateSource(resource.Source(source.Name), source.URL, format)

	if err != nil {
		ch.respondWithError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteSource handles DELETE /sources to delete a source.
func (ch *ControlHandler) DeleteSource(w http.ResponseWriter, r *http.Request) {
	var source struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	ch.manager.DeleteSource(source.Name)
	w.WriteHeader(http.StatusOK)
}

// respondWithError responds with an error.
func (ch *ControlHandler) respondWithError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := NewErrorResponse(err)
	_, _ = w.Write(errorResponse.getJSON())
}
