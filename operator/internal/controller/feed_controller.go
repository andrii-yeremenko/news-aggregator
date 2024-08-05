package controller

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	newsaggregatorv1 "com.teamdev/news-aggregator/api/v1"
)

const (
	serviceURL      = "https://news-aggregator.news-aggregator-namespace.svc.cluster.local:443"
	sourcesEndpoint = "/sources"
)

// FeedReconcile reconciles a Feed object.
type FeedReconcile struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main Kubernetes reconciliation loop, which aims to
// move the current state of the cluster closer to the desired state.
func (r *FeedReconcile) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var feed newsaggregatorv1.Feed
	if err := r.Get(ctx, req.NamespacedName, &feed); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	data := map[string]string{
		"name":   feed.Spec.Name,
		"url":    feed.Spec.Link,
		"format": "RSS",
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		logger.Error(err, "Failed to marshal data")
		return ctrl.Result{}, err
	}

	if err := r.addSource(dataBytes); err != nil {
		logger.Error(err, "Failed to add source")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully added source")
	return ctrl.Result{}, nil
}

// newInsecureHTTPClient creates a custom HTTP client with insecure transport.
func newInsecureHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

// addSource sends a POST request to add a source.
func (r *FeedReconcile) addSource(data []byte) error {
	u, err := url.JoinPath(serviceURL, sourcesEndpoint)

	if err != nil {
		return fmt.Errorf("failed to join URL path: %w", err)
	}

	httpClient := newInsecureHTTPClient()

	resp, err := httpClient.Post(u, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer closeResponseBody(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add source: %s", resp.Status)
	}

	return nil
}

// closeResponseBody ensures that the response body is closed and logs any errors.
func closeResponseBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Log.Error(err, "Failed to close response body")
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *FeedReconcile) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&newsaggregatorv1.Feed{}).
		Complete(r)
}
