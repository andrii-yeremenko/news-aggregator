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
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	newsaggregatorv1 "com.teamdev/news-aggregator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	sourcesEndpoint = "/sources"
)

// +kubebuilder:rbac:groups=news-aggregator.com.teamdev,resources=feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=news-aggregator.com.teamdev,resources=feeds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=news-aggregator.com.teamdev,resources=feeds/finalizers,verbs=update

//go:generate mockgen -destination=mocks/mock_http_client.go -package=mocks com.teamdev/news-aggregator/internal/controller HTTPClient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

type FeedReconcile struct {
	client.Client
	Scheme     *runtime.Scheme
	HTTPClient HTTPClient
	ServiceURL string
	Finalizer  string
}

// Reconcile performs the reconciliation logic for Feed objects.
func (r *FeedReconcile) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := ctrl.LoggerFrom(ctx)

	var feed newsaggregatorv1.Feed

	if err := r.Client.Get(ctx, req.NamespacedName, &feed); err != nil {
		if client.IgnoreNotFound(err) == nil {
			l.Info("Feed resource not found. It might have been deleted.", "name", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if err := r.processFeed(&feed, ctx); err != nil {
		return ctrl.Result{}, err
	}

	l.Info("Successfully reconciled Feed", "name", req.NamespacedName)

	return ctrl.Result{}, nil
}

// processFeed handles the logic for the Feed object based on its state.
func (r *FeedReconcile) processFeed(feed *newsaggregatorv1.Feed, ctx context.Context) error {

	if feed.ObjectMeta.DeletionTimestamp.IsZero() {
		err := r.ensureFinalizer(feed, ctx)
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
func (r *FeedReconcile) ensureFinalizer(feed *newsaggregatorv1.Feed, ctx context.Context) error {
	if !containsFinalizer(feed.Finalizers, r.Finalizer) {
		feed.Finalizers = append(feed.Finalizers, r.Finalizer)
		if err := r.Client.Update(ctx, feed); err != nil {
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
		statusErr := r.updateStatus(feed, newsaggregatorv1.ConditionFailed, false, err.Error())
		if statusErr != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}
		return err
	}
	return r.updateStatus(feed, newsaggregatorv1.ConditionAdded, true, "Feed added successfully")
}

// handleFeedUpdation handles the logic for when a Feed is updated.
func (r *FeedReconcile) handleFeedUpdation(feed *newsaggregatorv1.Feed) error {
	if err := r.updateSource(feed); err != nil {
		statusErr := r.updateStatus(feed, newsaggregatorv1.ConditionFailed, false, err.Error())
		if statusErr != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}
		return err
	}
	return r.updateStatus(feed, newsaggregatorv1.ConditionUpdated, true, "Feed updated successfully")
}

// handleFeedDeletion handles the logic for when a Feed is deleted.
func (r *FeedReconcile) handleFeedDeletion(feed *newsaggregatorv1.Feed) error {
	if containsFinalizer(feed.Finalizers, r.Finalizer) {
		if err := r.deleteSource(feed.Spec.Name); err != nil {
			return err
		}

		if err := r.updateStatus(feed, newsaggregatorv1.ConditionDeleted, true,
			"Feed deleted successfully"); err != nil {
			return err
		}

		feed.Finalizers = removeFinalizer(feed.Finalizers, r.Finalizer)

		contextWithTimeout, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
		defer cancel()

		if err := r.Client.Update(contextWithTimeout, feed); err != nil {
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

	contextWithTimeout, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	return r.Client.Status().Update(contextWithTimeout, feed)
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
		"format": "RSS", // By default, the format is set to RSS
	}

	return json.Marshal(data)
}

// postSource sends the POST request to add the source.
func (r *FeedReconcile) postSource(data []byte) error {
	u, err := url.JoinPath(r.ServiceURL, sourcesEndpoint)
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
	u, err := url.JoinPath(r.ServiceURL, sourcesEndpoint)
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

	u, err := url.JoinPath(r.ServiceURL, sourcesEndpoint)
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
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
