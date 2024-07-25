package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/cmd/web_server/handler/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestControlHandler_GetSources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mocks.NewMockResourceManager(ctrl)
	mockManager.EXPECT().AvailableSources().Return("source1,source2")

	handler := NewFeedsManagerHandler(mockManager)

	req := httptest.NewRequest(http.MethodGet, "/sources", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var sources string
	err := json.NewDecoder(resp.Body).Decode(&sources)
	assert.NoError(t, err)
	assert.Equal(t, "source1,source2", sources)
}

func TestControlHandler_AddSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
		mockManager.EXPECT().
			RegisterSource(resource.Source("source1"), "http://example.com", resource.Format(3)).
			Return(nil)

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
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
		handler := NewFeedsManagerHandler(mockManager)

		req := httptest.NewRequest(http.MethodPost, "/sources", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("unknown format", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
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
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("registration error", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
		mockManager.EXPECT().
			RegisterSource(resource.Source("source1"), "http://example.com", resource.Format(3)).
			Return(fmt.Errorf("registration error"))

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
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestControlHandler_UpdateSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
		mockManager.EXPECT().
			UpdateSource(resource.Source("source1"), "http://example.com", resource.Format(3)).
			Return(nil)

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
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
		handler := NewFeedsManagerHandler(mockManager)

		req := httptest.NewRequest(http.MethodPut, "/sources", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("unknown format", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
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
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("update error", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
		mockManager.EXPECT().
			UpdateSource(resource.Source("source1"), "http://example.com", resource.Format(3)).
			Return(fmt.Errorf("update error"))

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
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestControlHandler_DeleteSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
		mockManager.EXPECT().
			DeleteSource(resource.Source("source1")).
			Return(nil)

		handler := NewFeedsManagerHandler(mockManager)

		sourceData := map[string]string{
			"name": "source1",
		}
		body, _ := json.Marshal(sourceData)

		req := httptest.NewRequest(http.MethodDelete, "/sources", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
		handler := NewFeedsManagerHandler(mockManager)

		req := httptest.NewRequest(http.MethodDelete, "/sources", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()
		handler.Handle(w, req)

		resp := w.Result()
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("deletion error", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
		mockManager.EXPECT().
			DeleteSource(resource.Source("source1")).
			Return(fmt.Errorf("deletion error"))

		handler := NewFeedsManagerHandler(mockManager)

		sourceData := map[string]string{
			"name": "source1",
		}
		body, _ := json.Marshal(sourceData)

		req := httptest.NewRequest(http.MethodDelete, "/sources", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Handle(w, req)

		resp := w.Result()
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}(resp.Body)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
