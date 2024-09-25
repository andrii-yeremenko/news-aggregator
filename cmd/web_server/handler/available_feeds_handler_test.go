package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"news-aggregator/cmd/web_server/handler/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type errorResponseWriter struct {
	http.ResponseWriter
}

func (e *errorResponseWriter) Write(b []byte) (int, error) {
	return 0, errors.New("simulated write error")
}

func TestAvailableFeedsHandler_GetSources_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mocks.NewMockResourceManager(ctrl)
	feedData := `feed1,feed2`
	mockManager.EXPECT().AvailableFeeds().Return(feedData)

	handler := NewAvailableFeedsHandler(mockManager)
	req := httptest.NewRequest(http.MethodGet, "/feeds", nil)
	rr := httptest.NewRecorder()

	handler.Handle(rr, req)

	res := rr.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var responseBody string
	err := json.NewDecoder(res.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, feedData, "feed1,feed2")
}

func TestAvailableFeedsHandler_Handle_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mocks.NewMockResourceManager(ctrl)
	handler := NewAvailableFeedsHandler(mockManager)
	req := httptest.NewRequest(http.MethodPost, "/feeds", nil)
	rr := httptest.NewRecorder()

	handler.Handle(rr, req)

	res := rr.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestAvailableFeedsHandler_GetSources_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockManager := mocks.NewMockResourceManager(ctrl)
	feedData := "feed1,feed2"
	mockManager.EXPECT().AvailableFeeds().Return(feedData)

	handler := NewAvailableFeedsHandler(mockManager)
	req := httptest.NewRequest(http.MethodGet, "/feeds", nil)

	rr := httptest.NewRecorder()
	errorWriter := &errorResponseWriter{rr}

	handler.Handle(errorWriter, req)

	res := rr.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
