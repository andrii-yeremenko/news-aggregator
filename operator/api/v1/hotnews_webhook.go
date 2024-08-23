package v1

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
)

// log is for logging in this package.
var hotnewsLog = logf.Log.WithName("hotnews-resource")
var configMapName string

// SetupWebhookWithManager sets up the webhook with the manager.
func (r *HotNews) SetupWebhookWithManager(mgr ctrl.Manager, mapName string) error {
	configMapName = mapName
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/validate-news-aggregator-com-teamdev-v1-hotnews,mutating=false,failurePolicy=fail,sideEffects=None,groups=news-aggregator.com.teamdev,resources=hotnews,verbs=create;update,versions=v1,name=vhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HotNews{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *HotNews) ValidateCreate() (admission.Warnings, error) {
	hotnewsLog.Info("validate create", "name", r.Name)

	if err := r.validateHotNews(); err != nil {
		return nil, err
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *HotNews) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	hotnewsLog.Info("validate update", "name", r.Name)

	if err := r.validateHotNews(); err != nil {
		return nil, err
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *HotNews) ValidateDelete() (admission.Warnings, error) {
	hotnewsLog.Info("validate delete", "name", r.Name)
	return nil, nil
}

// validateHotNews contains the core validation logic for HotNews.
func (r *HotNews) validateHotNews() error {
	var errors []string

	if len(r.Spec.Keywords) == 0 {
		errors = append(errors, "keywords must be provided")
	}

	if r.Spec.DateStart != nil && r.Spec.DateEnd != nil {
		if !r.Spec.DateEnd.After(r.Spec.DateStart.Time) {
			errors = append(errors, "dateEnd must be after dateStart")
		}
	} else if r.Spec.DateStart == nil || r.Spec.DateEnd == nil {
		errors = append(errors, "both dateStart and dateEnd must be provided")
	}

	if err := r.validateFeedGroups(); err != nil {
		errors = append(errors, err.Error())
	}

	if err := r.validateFeeds(); err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}

	return nil
}

// validateFeeds checks if the feeds specified in the HotNews resource exist.
func (r *HotNews) validateFeeds() error {

	if r.Spec.Feeds == nil {
		return nil
	}

	var notFoundFeeds []string

	for _, feed := range r.Spec.Feeds {
		var f Feed
		if err := k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: r.Namespace,
			Name:      feed,
		}, &f); err != nil {
			notFoundFeeds = append(notFoundFeeds, feed)
		}
	}

	if len(notFoundFeeds) > 0 {
		errStr := strings.Join(notFoundFeeds, ", ")
		return fmt.Errorf("feeds \"%s\" do not exist", errStr)
	}

	return nil
}

// validateFeedGroups checks if the feed groups specified in the HotNews resource exist in the ConfigMap.
func (r *HotNews) validateFeedGroups() error {

	var cm v1.ConfigMap
	if err := k8sClient.Get(context.TODO(), client.ObjectKey{
		Namespace: r.Namespace,
		Name:      configMapName,
	}, &cm); err != nil {
		return fmt.Errorf("failed to retrieve ConfigMap %s/%s: %v", r.Namespace, configMapName, err)
	}

	if r.Spec.FeedGroups == nil {
		return nil
	}

	for _, group := range r.Spec.FeedGroups {
		if _, exists := cm.Data[group]; !exists {
			return fmt.Errorf("feedGroup %s does not exist in ConfigMap %s/%s", group, r.Namespace, configMapName)
		}
	}

	return nil
}
