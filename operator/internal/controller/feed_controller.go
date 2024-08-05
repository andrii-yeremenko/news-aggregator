package controller

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	newsaggregatorv1 "com.teamdev/news-aggregator/api/v1"
)

// FeedReconciler reconciles a Feed object
type FeedReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *FeedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

	serviceURL := "https://news-aggregator.news-aggregator-namespace.svc.cluster.local:443"

	if err := r.addSource(serviceURL, dataBytes); err != nil {
		logger.Error(err, "Failed to add source")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully added source")

	return ctrl.Result{}, nil
}

// Create a custom HTTP client with insecure transport
func newInsecureHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

// Update your addSource method to use this custom client
func (r *FeedReconciler) addSource(serviceURL string, data []byte) error {
	url := fmt.Sprintf("%s/source", serviceURL)
	client := newInsecureHTTPClient()
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Log.Error(err, "Failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add source: %s", resp.Status)
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FeedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&newsaggregatorv1.Feed{}).
		Complete(r)
}
