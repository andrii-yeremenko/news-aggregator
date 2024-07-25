package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/cmd/web_server/handler/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("valid source", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
		mockManager.EXPECT().IsSourceSupported(resource.Source("validSource")).Return(true)
		mockManager.EXPECT().UpdateResource(resource.Source("validSource")).Return(nil)

		handler := NewUpdateHandler(mockManager)

		req := httptest.NewRequest(http.MethodGet, "/update?source=validSource", nil)
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

	t.Run("source not specified", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)

		handler := NewUpdateHandler(mockManager)

		req := httptest.NewRequest(http.MethodGet, "/update", nil)
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

	t.Run("unsupported source", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
		mockManager.EXPECT().IsSourceSupported(resource.Source("unsupportedSource")).Return(false)

		handler := NewUpdateHandler(mockManager)

		req := httptest.NewRequest(http.MethodGet, "/update?source=unsupportedSource", nil)
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
		assert.Contains(t, readResponseBody(resp), "Source not supported")
	})

	t.Run("update error", func(t *testing.T) {
		mockManager := mocks.NewMockResourceManager(ctrl)
		mockManager.EXPECT().IsSourceSupported(resource.Source("validSource")).Return(true)
		mockManager.EXPECT().UpdateResource(resource.Source("validSource")).Return(assert.AnError)

		handler := NewUpdateHandler(mockManager)

		req := httptest.NewRequest(http.MethodGet, "/update?source=validSource", nil)
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
		assert.Contains(t, readResponseBody(resp), "Failed to update source")
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
