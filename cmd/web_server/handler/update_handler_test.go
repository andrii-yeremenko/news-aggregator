package handler

import (
	"net/http"
	"net/http/httptest"
	"news-aggregator/aggregator/model/resource"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateHandler_Handle(t *testing.T) {
	t.Run("valid source", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		mockManager.On("IsSourceSupported", resource.Source("validSource")).Return(true)
		mockManager.On("UpdateResource", resource.Source("validSource")).Return(nil)

		handler := NewUpdateHandler(mockManager)

		req := httptest.NewRequest(http.MethodGet, "/update?source=validSource", nil)
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		mockManager.AssertExpectations(t)
	})

	t.Run("source not specified", func(t *testing.T) {
		mockManager := new(MockResourceManager)

		handler := NewUpdateHandler(mockManager)

		req := httptest.NewRequest(http.MethodGet, "/update", nil)
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockManager.AssertExpectations(t)
	})

	t.Run("unsupported source", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		mockManager.On("IsSourceSupported", resource.Source("unsupportedSource")).Return(false)

		handler := NewUpdateHandler(mockManager)

		req := httptest.NewRequest(http.MethodGet, "/update?source=unsupportedSource", nil)
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, readResponseBody(resp), "Source not supported")

		mockManager.AssertExpectations(t)
	})

	t.Run("update error", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		mockManager.On("IsSourceSupported", resource.Source("validSource")).Return(true)
		mockManager.On("UpdateResource", resource.Source("validSource")).Return(assert.AnError)

		handler := NewUpdateHandler(mockManager)

		req := httptest.NewRequest(http.MethodGet, "/update?source=validSource", nil)
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Contains(t, readResponseBody(resp), "Failed to update source")

		mockManager.AssertExpectations(t)
	})
}

// readResponseBody reads and returns the response body as a string.
func readResponseBody(resp *http.Response) string {
	body := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		body = append(body, buf[:n]...)
		if err != nil {
			break
		}
	}
	return string(body)
}
