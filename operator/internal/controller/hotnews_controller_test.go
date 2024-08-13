package controller_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	newsaggregatorv1 "com.teamdev/news-aggregator/api/v1"
	corev1 "k8s.io/api/core/v1"

	"com.teamdev/news-aggregator/internal/controller"
)

// MockHTTPClient is a mock for HTTPClient interface.
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	panic("unexpected call to Post")
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestHotNewsReconciler_Reconcile_Success(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = newsaggregatorv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	hotNews := &newsaggregatorv1.HotNews{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hotnews",
			Namespace: "default",
		},
		Spec: newsaggregatorv1.HotNewsSpec{
			Feeds: []string{"feed1", "feed2"},
			SummaryConfig: newsaggregatorv1.SummaryConfig{
				TitlesCount: 2,
			},
		},
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "news-feeds",
			Namespace: "default",
		},
		Data: map[string]string{
			"group1": "feed1,feed2",
		},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(hotNews, configMap).Build()
	httpClient := new(MockHTTPClient)

	reconciler := &controller.HotNewsReconciler{
		Client:             client,
		HTTPClient:         httpClient,
		Scheme:             scheme,
		NewsAggregatorURL:  "http://news-aggregator",
		ConfigMapName:      "news-feeds",
		ConfigMapNamespace: "default",
	}

	httpResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`[{"title": "Title 1"}, {"title": "Title 2"}]`)),
	}
	httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(httpResponse, nil)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-hotnews",
			Namespace: "default",
		},
	}
	result, err := reconciler.Reconcile(context.Background(), req)

	assert.Equal(t, reconcile.Result{}, result)

	updatedHotNews := &newsaggregatorv1.HotNews{}
	err = client.Get(context.Background(), req.NamespacedName, updatedHotNews)
	assert.NoError(t, err)

	httpClient.AssertExpectations(t)
}

func TestHotNewsReconciler_Reconcile_FetchNewsError(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = newsaggregatorv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	hotNews := &newsaggregatorv1.HotNews{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hotnews",
			Namespace: "default",
		},
		Spec: newsaggregatorv1.HotNewsSpec{
			Feeds: []string{"feed1", "feed2"},
			SummaryConfig: newsaggregatorv1.SummaryConfig{
				TitlesCount: 2,
			},
		},
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "news-feeds",
			Namespace: "default",
		},
		Data: map[string]string{
			"group1": "feed1,feed2",
		},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(hotNews, configMap).Build()
	httpClient := new(MockHTTPClient)

	reconciler := &controller.HotNewsReconciler{
		Client:             client,
		HTTPClient:         httpClient,
		Scheme:             scheme,
		NewsAggregatorURL:  "http://news-aggregator",
		ConfigMapName:      "news-feeds",
		ConfigMapNamespace: "default",
	}

	httpResponse := &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader(`Internal Server Error`)),
	}
	httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(httpResponse, errors.New("failed to fetch news"))

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-hotnews",
			Namespace: "default",
		},
	}
	_, err := reconciler.Reconcile(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch news")

	httpClient.AssertExpectations(t)
}
