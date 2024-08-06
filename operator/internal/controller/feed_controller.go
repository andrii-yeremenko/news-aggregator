package controller

import (
	"bytes"
	newsaggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
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

	var feed newsaggregatorv1.Feed
	if err := r.Get(ctx, req.NamespacedName, &feed); err != nil {
		err := r.handleDelete(req.Name)
		return ctrl.Result{}, err
	}

	err := r.handleAdd(ctx, &feed)

	return ctrl.Result{}, err
}

func (r *FeedReconcile) handleAdd(ctx context.Context, feed *newsaggregatorv1.Feed) error {
	data := map[string]string{
		"name":   feed.Spec.Name,
		"url":    feed.Spec.Link,
		"format": "RSS",
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	var httpErr error

	switch feed.Status.CurrentState {
	case "":
		httpErr = r.addSource(dataBytes)
		feed.Status.CurrentState = "Added"
	case "Failed":
		httpErr = r.addSource(dataBytes)
		feed.Status.CurrentState = "Added"
	case "Added":
		return nil
	}

	if httpErr != nil {
		feed.Status.CurrentState = "Failed"
		feed.Status.Message = httpErr.Error()
		return fmt.Errorf("failed to add source: %w", httpErr)
	} else {
		feed.Status.Message = "Source " + feed.Spec.Name + " added successfully"
	}

	feed.Status.LastUpdated = metav1.Now().String()

	if err := r.Client.Status().Update(ctx, feed); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return err
}

func (r *FeedReconcile) handleDelete(feedName string) error {

	log.Log.Info("Trying deleting source", "name", feedName)

	data := map[string]string{
		"name": feedName,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := r.deleteSource(dataBytes); err != nil {
		return fmt.Errorf("failed to delete source: %w", err)
	}

	return nil
}

// addSource sends a POST request to add a source.
func (r *FeedReconcile) addSource(data []byte) error {

	log.Log.Info("Trying adding source", "data", string(data))

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

// deleteSource sends a DELETE request to remove a source.
func (r *FeedReconcile) deleteSource(data []byte) error {
	u, err := url.JoinPath(serviceURL, sourcesEndpoint)
	if err != nil {
		return fmt.Errorf("failed to join URL path: %w", err)
	}

	httpClient := newInsecureHTTPClient()

	req, err := http.NewRequest(http.MethodDelete, u, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send DELETE request: %w", err)
	}
	defer closeResponseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete source: %s", resp.Status)
	}

	log.Log.Info("Source deleted successfully")

	return nil
}

// newInsecureHTTPClient creates a custom HTTP client with insecure transport.
func newInsecureHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
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
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				log.Log.Info("Create event", "feed", e.Object)
				return true
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				log.Log.Info("Update event", "feed", e.ObjectNew)
				return true
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				log.Log.Info("Delete event", "feed", e.Object)
				return true
			},
		}).
		Complete(r)
}
