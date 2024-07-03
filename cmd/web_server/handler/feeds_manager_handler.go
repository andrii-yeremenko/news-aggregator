package handler

import (
	"encoding/json"
	"net/http"
	"news-aggregator/aggregator/model/resource"
)

// FeedsManagerHandler handles requests for managing news sources.
type FeedsManagerHandler struct {
	manager ResourceManager
}

// NewFeedsManagerHandler creates a new FeedsManagerHandler instance.
func NewFeedsManagerHandler(manager ResourceManager) *FeedsManagerHandler {
	return &FeedsManagerHandler{
		manager: manager,
	}
}

// Handle routes the request based on the HTTP method.
func (ch *FeedsManagerHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ch.GetSources(w)
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
func (ch *FeedsManagerHandler) GetSources(w http.ResponseWriter) {
	sources := ch.manager.AvailableSources()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(sources)
	if err != nil {
		http.Error(w, "Failed to encode sources", http.StatusInternalServerError)
		return
	}
}

// AddSource handles POST /sources to add a new source.
func (ch *FeedsManagerHandler) AddSource(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ch.manager.RegisterSource(resource.Source(source.Name), source.URL, format)
	if err != nil {
		http.Error(w, "Failed to add source", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateSource handles PUT /sources to update an existing source.
func (ch *FeedsManagerHandler) UpdateSource(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ch.manager.UpdateSource(resource.Source(source.Name), source.URL, format)

	if err != nil {
		http.Error(w, "Failed to update source", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteSource handles DELETE /sources to delete a source.
func (ch *FeedsManagerHandler) DeleteSource(w http.ResponseWriter, r *http.Request) {
	var source struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := ch.manager.DeleteSource(resource.Source(source.Name))
	if err != nil {
		http.Error(w, "Failed to delete source", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
