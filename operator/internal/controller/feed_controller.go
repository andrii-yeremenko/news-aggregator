package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	newsaggregatorv1 "com.teamdev/news-aggregator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	serviceURL      = "https://news-aggregator.news-aggregator-namespace.svc.cluster.local:443"
	sourcesEndpoint = "/sources"
	finalizer       = "feed.finalizer.news-aggregator.teamdev.com"
)

//go:generate mockgen -destination=mocks/mock_http_client.go -package=mocks com.teamdev/news-aggregator/internal/controller HTTPClient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

type FeedReconcile struct {
	client.Client
	Scheme     *runtime.Scheme
	HTTPClient HTTPClient
}

// Reconcile performs the reconciliation logic for Feed objects.
func (r *FeedReconcile) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var feed newsaggregatorv1.Feed

	if err := r.Client.Get(ctx, req.NamespacedName, &feed); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.processFeed(&feed); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// processFeed handles the main reconciliation logic based on the Feed's state.
func (r *FeedReconcile) processFeed(feed *newsaggregatorv1.Feed) error {

	if feed.ObjectMeta.DeletionTimestamp.IsZero() {
		err := r.ensureFinalizer(feed)
		if err != nil {
			return err
		}
	}

	if len(feed.Status.Conditions) == 0 {
		return r.handleFeedCreation(feed)
	}

	lastCondition := feed.Status.Conditions[len(feed.Status.Conditions)-1]
	if lastCondition.Type == newsaggregatorv1.ConditionDeleted {
		return nil
	}

	if !feed.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.handleFeedDeletion(feed)
	}

	if !feed.CreationTimestamp.IsZero() {
		return r.handleFeedUpdation(feed)
	}

	return fmt.Errorf("unexpected state for Feed %s", feed.Spec.Name)
}

// ensureFinalizer ensures that the finalizer is added to the Feed if not already present.
func (r *FeedReconcile) ensureFinalizer(feed *newsaggregatorv1.Feed) error {
	if !containsFinalizer(feed.Finalizers, finalizer) {
		feed.Finalizers = append(feed.Finalizers, finalizer)
		if err := r.Client.Update(context.Background(), feed); err != nil {
			return err
		}
		log.Log.Info("Added finalizer", "name", feed.Spec.Name)
	}
	return nil
}

// handleFeedCreation handles the logic for when a Feed is created.
func (r *FeedReconcile) handleFeedCreation(feed *newsaggregatorv1.Feed) error {
	err := r.addSource(feed)
	if err != nil {
		return r.updateStatus(feed, newsaggregatorv1.ConditionFailed, false, err.Error())
	}
	return r.updateStatus(feed, newsaggregatorv1.ConditionAdded, true, "Feed added successfully")
}

// handleFeedUpdation handles the logic for when a Feed is updated.
func (r *FeedReconcile) handleFeedUpdation(feed *newsaggregatorv1.Feed) error {
	if err := r.updateSource(feed); err != nil {
		return r.updateStatus(feed, newsaggregatorv1.ConditionFailed, false, err.Error())
	}
	return r.updateStatus(feed, newsaggregatorv1.ConditionUpdated, true, "Feed updated successfully")
}

// handleFeedDeletion handles the logic for when a Feed is deleted.
func (r *FeedReconcile) handleFeedDeletion(feed *newsaggregatorv1.Feed) error {
	if containsFinalizer(feed.Finalizers, finalizer) {
		if err := r.deleteSource(feed.Spec.Name); err != nil {
			return err
		}

		if err := r.updateStatus(feed, newsaggregatorv1.ConditionDeleted, true,
			"Feed deleted successfully"); err != nil {
			return err
		}

		feed.Finalizers = removeFinalizer(feed.Finalizers, finalizer)
		if err := r.Client.Update(context.Background(), feed); err != nil {
			return err
		}
		log.Log.Info("Removed finalizer", "name", feed.Spec.Name)
	}
	return nil
}

// updateStatus updates the status of the Feed object.
func (r *FeedReconcile) updateStatus(feed *newsaggregatorv1.Feed, conditionType newsaggregatorv1.ConditionType,
	status bool, message string) error {
	condition := newsaggregatorv1.Condition{
		Type:           conditionType,
		Status:         status,
		Message:        message,
		LastUpdateTime: metav1.Now(),
	}

	feed.Status.Conditions = append(feed.Status.Conditions, condition)
	return r.Client.Status().Update(context.Background(), feed)
}

// addSource sends a POST request to add a source.
func (r *FeedReconcile) addSource(feed *newsaggregatorv1.Feed) error {
	data, err := r.prepareSourceData(feed)
	if err != nil {
		return err
	}

	return r.postSource(data)
}

// updateSource sends a PUT request to update a source.
func (r *FeedReconcile) updateSource(feed *newsaggregatorv1.Feed) error {
	data, err := r.prepareSourceData(feed)
	if err != nil {
		return err
	}

	return r.putSource(data)
}

// prepareSourceData prepares the source data to be sent in the POST request.
func (r *FeedReconcile) prepareSourceData(feed *newsaggregatorv1.Feed) ([]byte, error) {
	data := map[string]string{
		"name":   feed.Spec.Name,
		"url":    feed.Spec.Link,
		"format": "RSS",
	}

	return json.Marshal(data)
}

// postSource sends the POST request to add the source.
func (r *FeedReconcile) postSource(data []byte) error {
	u, err := url.JoinPath(serviceURL, sourcesEndpoint)
	if err != nil {
		return fmt.Errorf("failed to join URL path: %w", err)
	}

	resp, err := r.HTTPClient.Post(u, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer closeResponseBody(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add source: %s", resp.Status)
	}

	return nil
}

// putSource sends a PUT request to update a source.
func (r *FeedReconcile) putSource(data []byte) error {
	u, err := url.JoinPath(serviceURL, sourcesEndpoint)
	if err != nil {
		return fmt.Errorf("failed to join URL path: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, u, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send PUT request: %w", err)
	}
	defer closeResponseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update source: %s", resp.Status)
	}

	log.Log.Info("Source updated successfully")
	return nil
}

// deleteSource sends a DELETE request to remove a source.
func (r *FeedReconcile) deleteSource(feedName string) error {
	data, err := r.prepareDeleteData(feedName)
	if err != nil {
		return err
	}

	u, err := url.JoinPath(serviceURL, sourcesEndpoint)
	if err != nil {
		return fmt.Errorf("failed to join URL path: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, u, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.HTTPClient.Do(req)
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

// prepareDeleteData prepares the data for the DELETE request.
func (r *FeedReconcile) prepareDeleteData(feedName string) ([]byte, error) {
	data := map[string]string{"name": feedName}
	return json.Marshal(data)
}

// closeResponseBody ensures that the response body is closed and logs any errors.
func closeResponseBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Log.Error(err, "Failed to close response body")
	}
}

// containsFinalizer checks if a finalizer is present in the list.
func containsFinalizer(finalizers []string, finalizer string) bool {
	for _, f := range finalizers {
		if f == finalizer {
			return true
		}
	}
	return false
}

// removeFinalizer removes a finalizer from the list.
func removeFinalizer(finalizers []string, finalizer string) []string {
	for i, f := range finalizers {
		if f == finalizer {
			return append(finalizers[:i], finalizers[i+1:]...)
		}
	}
	return finalizers
}

// SetupWithManager sets up the controller with the Manager.
func (r *FeedReconcile) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&newsaggregatorv1.Feed{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				return true
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				return !e.DeleteStateUnknown
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				return e.ObjectNew.GetGeneration() != e.ObjectOld.GetGeneration()
			},
		}).
		Complete(r)
}
