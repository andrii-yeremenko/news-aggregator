package v1

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var feedlog = log.Log.WithName("feed-resource")

// Client for querying Kubernetes API
var k8sClient client.Client

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *Feed) SetupWebhookWithManager(mgr ctrl.Manager) error {
	k8sClient = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-news-aggregator-com-teamdev-v1-feed,mutating=true,failurePolicy=fail,sideEffects=None,groups=news-aggregator.com.teamdev,resources=feeds,verbs=create;update,versions=v1,name=mfeed.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Feed{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Feed) Default() {
	feedlog.Info("default", "name", r.Name)
}

// +kubebuilder:webhook:path=/validate-news-aggregator-com-teamdev-v1-feed,mutating=false,failurePolicy=fail,sideEffects=None,groups=news-aggregator.com.teamdev,resources=feeds,verbs=create;update,versions=v1,name=vfeed.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Feed{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Feed) ValidateCreate() (admission.Warnings, error) {
	feedlog.Info("validate create", "name", r.Name)

	if err := validateFeed(r); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Feed) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	feedlog.Info("validate update", "name", r.Name)

	if err := validateFeed(r); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Feed) ValidateDelete() (admission.Warnings, error) {
	feedlog.Info("validate delete", "name", r.Name)

	return nil, nil
}

// validateFeed performs the validation checks on the Feed object.
func validateFeed(feed *Feed) error {
	if feed.Spec.Name == "" {
		return fmt.Errorf("name field cannot be empty")
	}

	if len(feed.Spec.Name) > 20 {
		return fmt.Errorf("name field cannot be more than 20 characters")
	}

	if !isValidName(feed.Spec.Name) {
		return fmt.Errorf("name field contains invalid characters")
	}

	if feed.Spec.Link == "" {
		return fmt.Errorf("link field cannot be empty")
	}

	if err := validateURL(feed.Spec.Link); err != nil {
		return err
	}

	if err := checkNameUniqueness(feed); err != nil {
		return err
	}

	return nil
}

// isValidName checks that the name contains only alphanumeric characters, dashes, and underscores.
func isValidName(name string) bool {
	for _, ch := range name {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_') {
			return false
		}
	}
	return true
}

// validateURL checks if the URL is valid.
func validateURL(link string) error {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("link field is not a valid URL")
	}
	return nil
}

// checkNameUniqueness queries the Kubernetes API to ensure that no other Feed with the same name exists in the same namespace.
func checkNameUniqueness(feed *Feed) error {
	feedList := &FeedList{}
	listOpts := client.ListOptions{Namespace: feed.Namespace}
	err := k8sClient.List(context.Background(), feedList, &listOpts)
	if err != nil {
		return fmt.Errorf("failed to list feeds: %v", err)
	}

	for _, existingFeed := range feedList.Items {
		if existingFeed.Spec.Name == feed.Spec.Name && existingFeed.Namespace == feed.Namespace && existingFeed.UID != feed.UID {
			return fmt.Errorf("a Feed with name '%s' already exists in namespace '%s'", feed.Spec.Name, feed.Namespace)
		}
	}

	return nil
}
