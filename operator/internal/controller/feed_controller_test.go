package controller

import (
	"bytes"
	newsaggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

const (
	finalizer  = "feed.finalizer.news-aggregator.teamdev.com"
	serviceURL = "https://news-aggregator.news-aggregator-namespace.svc.cluster.local:443"
)

func TestFeedReconcile(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = newsaggregatorv1.AddToScheme(scheme)

	initialFeed := &newsaggregatorv1.Feed{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-feed",
			Namespace: "default",
		},
		Spec: newsaggregatorv1.FeedSpec{
			Name: "Test Feed",
			Link: "https://example.com/rss",
		},
		Status: newsaggregatorv1.FeedStatus{
			Conditions: []newsaggregatorv1.Condition{},
		},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initialFeed).Build()

	feed := &newsaggregatorv1.Feed{}
	err := client.Get(context.Background(), types.NamespacedName{
		Name:      "test-feed",
		Namespace: "default",
	}, feed)
	assert.NoError(t, err, "initial Feed object should be found")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHTTPClient := mocks.NewMockHTTPClient(ctrl)

	mockHTTPClient.EXPECT().
		Post("https://news-aggregator.news-aggregator-namespace.svc.cluster.local:443/sources", "application/json", gomock.Any()).
		Return(&http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}, nil)

	r := &FeedReconcile{
		Client:     client,
		Scheme:     scheme,
		HTTPClient: mockHTTPClient,
		ServiceURL: serviceURL,
		Finalizer:  finalizer,
	}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-feed",
			Namespace: "default",
		},
	}

	res, err := r.Reconcile(context.Background(), req)
	assert.False(t, res.Requeue)

	// Retrieve the Feed object after reconciliation
	err = client.Get(context.Background(), req.NamespacedName, feed)
	assert.NoError(t, err, "Feed object should be found after reconciliation")

	// Verify that the finalizer was added
	if !assert.Contains(t, feed.Finalizers, finalizer, "Finalizer should be added") {
		t.Logf("Finalizers found: %v", feed.Finalizers)
	}
}
