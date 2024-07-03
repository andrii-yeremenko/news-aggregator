package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"news-aggregator/aggregator/model/resource"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockResourceManager is a mock implementation of the ResourceManager interface.
type MockResourceManager struct {
	mock.Mock
}

func (m *MockResourceManager) AvailableSources() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockResourceManager) RegisterSource(name resource.Source, url string, format resource.Format) error {
	args := m.Called(name, url, int(format))
	return args.Error(0)
}

func (m *MockResourceManager) UpdateSource(name resource.Source, url string, format resource.Format) error {
	args := m.Called(name, url, int(format))
	return args.Error(0)
}

func (m *MockResourceManager) UpdateResource(source resource.Source) error {
	args := m.Called(source)
	return args.Error(0)
}

func (m *MockResourceManager) DeleteSource(name resource.Source) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockResourceManager) IsSourceSupported(source resource.Source) bool {
	args := m.Called(source)
	return args.Bool(0)
}

func (m *MockResourceManager) GetAllResources() ([]resource.Resource, error) {
	args := m.Called()
	return args.Get(0).([]resource.Resource), args.Error(1)
}

func (m *MockResourceManager) GetSelectedResources(sources []string) ([]resource.Resource, error) {
	args := m.Called(sources)
	return args.Get(0).([]resource.Resource), args.Error(1)
}

func TestControlHandler_GetSources(t *testing.T) {
	mockManager := new(MockResourceManager)
	mockManager.On("AvailableSources").Return("source1,source2")

	handler := NewFeedsManagerHandler(mockManager)

	req := httptest.NewRequest(http.MethodGet, "/sources", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var sources string
	err := json.NewDecoder(resp.Body).Decode(&sources)
	assert.NoError(t, err)
	assert.Equal(t, "source1,source2", sources)
}

func TestControlHandler_AddSource(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		mockManager.On("RegisterSource", resource.Source("source1"), "http://example.com", int(resource.JSON)).Return(nil)

		handler := NewFeedsManagerHandler(mockManager)

		sourceData := map[string]string{
			"name":   "source1",
			"url":    "http://example.com",
			"format": "json",
		}
		body, _ := json.Marshal(sourceData)

		req := httptest.NewRequest(http.MethodPost, "/sources", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		handler := NewFeedsManagerHandler(mockManager)

		req := httptest.NewRequest(http.MethodPost, "/sources", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("unknown format", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		handler := NewFeedsManagerHandler(mockManager)

		sourceData := map[string]string{
			"name":   "source1",
			"url":    "http://example.com",
			"format": "unknown",
		}
		body, _ := json.Marshal(sourceData)

		req := httptest.NewRequest(http.MethodPost, "/sources", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("registration error", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		mockManager.On("RegisterSource", resource.Source("source1"), "http://example.com", int(resource.JSON)).Return(fmt.Errorf("registration error"))

		handler := NewFeedsManagerHandler(mockManager)

		sourceData := map[string]string{
			"name":   "source1",
			"url":    "http://example.com",
			"format": "json",
		}
		body, _ := json.Marshal(sourceData)

		req := httptest.NewRequest(http.MethodPost, "/sources", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestControlHandler_UpdateSource(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		mockManager.On("UpdateSource", resource.Source("source1"), "http://example.com", int(resource.JSON)).Return(nil)

		handler := NewFeedsManagerHandler(mockManager)

		sourceData := map[string]string{
			"name":   "source1",
			"url":    "http://example.com",
			"format": "json",
		}
		body, _ := json.Marshal(sourceData)

		req := httptest.NewRequest(http.MethodPut, "/sources", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		handler := NewFeedsManagerHandler(mockManager)

		req := httptest.NewRequest(http.MethodPut, "/sources", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("unknown format", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		handler := NewFeedsManagerHandler(mockManager)

		sourceData := map[string]string{
			"name":   "source1",
			"url":    "http://example.com",
			"format": "unknown",
		}
		body, _ := json.Marshal(sourceData)

		req := httptest.NewRequest(http.MethodPut, "/sources", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("update error", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		mockManager.On("UpdateSource", resource.Source("source1"), "http://example.com", int(resource.JSON)).Return(fmt.Errorf("update error"))

		handler := NewFeedsManagerHandler(mockManager)

		sourceData := map[string]string{
			"name":   "source1",
			"url":    "http://example.com",
			"format": "json",
		}
		body, _ := json.Marshal(sourceData)

		req := httptest.NewRequest(http.MethodPut, "/sources", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestControlHandler_DeleteSource(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		mockManager.On("DeleteSource", resource.Source("source1")).Return(nil)

		handler := NewFeedsManagerHandler(mockManager)

		sourceData := map[string]string{
			"name": "source1",
		}
		body, _ := json.Marshal(sourceData)

		req := httptest.NewRequest(http.MethodDelete, "/sources", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		handler := NewFeedsManagerHandler(mockManager)

		req := httptest.NewRequest(http.MethodDelete, "/sources", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("deletion error", func(t *testing.T) {
		mockManager := new(MockResourceManager)
		mockManager.On("DeleteSource", resource.Source("source1")).Return(fmt.Errorf("deletion error"))

		handler := NewFeedsManagerHandler(mockManager)

		sourceData := map[string]string{
			"name": "source1",
		}
		body, _ := json.Marshal(sourceData)

		req := httptest.NewRequest(http.MethodDelete, "/sources", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
