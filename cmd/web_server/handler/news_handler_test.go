package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"news-aggregator/aggregator"
	"news-aggregator/manager"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewsAggregatorHandler_Handle(t *testing.T) {
	t.Run("valid request with articles", func(t *testing.T) {
		managerConfigPath := "../../../config/feeds_dictionary.json"
		storagePath := "../../../resources"

		manager, err := manager.New(storagePath, managerConfigPath)
		parserFactory := aggregator.NewParserFactory()

		handler := &NewsAggregatorHandler{
			resourceManager: manager,
			parserPool:      parserFactory,
		}

		req := httptest.NewRequest(http.MethodGet, "/news?sort=asc&sources=usa-today&keywords=Ukraine", nil)
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Error(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var articlesJSON []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&articlesJSON)
		assert.NoError(t, err)
		assert.NotEmpty(t, articlesJSON)

	})

	t.Run("invalid source", func(t *testing.T) {
		managerConfigPath := "../../../config/feeds_dictionary.json"
		storagePath := "../../../resources"

		manager, err := manager.New(storagePath, managerConfigPath)

		assert.NoError(t, err)

		parserFactory := aggregator.NewParserFactory()

		handler := &NewsAggregatorHandler{
			resourceManager: manager,
			parserPool:      parserFactory,
		}

		// Create request and response recorder
		req := httptest.NewRequest(http.MethodGet, "/?sources=invalidSource", nil)
		w := httptest.NewRecorder()

		// Call the handler
		handler.Handle(w, req)

		// Check the response
		resp := w.Result()
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Error(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
