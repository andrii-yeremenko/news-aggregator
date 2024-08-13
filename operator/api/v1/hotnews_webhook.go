package v1

import (
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var hotnewslog = logf.Log.WithName("hotnews-resource")

// SetupWebhookWithManager sets up the webhook with the manager.
func (r *HotNews) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/validate-news-aggregator-com-teamdev-v1-hotnews,mutating=false,failurePolicy=fail,sideEffects=None,groups=news-aggregator.com.teamdev,resources=hotnews,verbs=create;update,versions=v1,name=vhotnews.kb.io,admissionReviewVersions=v1
var _ webhook.Validator = &HotNews{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *HotNews) ValidateCreate() (admission.Warnings, error) {
	hotnewslog.Info("validate create", "name", r.Name)

	if err := r.validateHotNews(); err != nil {
		return nil, err
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *HotNews) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	hotnewslog.Info("validate update", "name", r.Name)

	if err := r.validateHotNews(); err != nil {
		return nil, err
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *HotNews) ValidateDelete() (admission.Warnings, error) {
	hotnewslog.Info("validate delete", "name", r.Name)
	return nil, nil
}

// validateHotNews contains the core validation logic for HotNews.
func (r *HotNews) validateHotNews() error {
	if len(r.Spec.Keywords) == 0 {
		return errors.New("keywords must not be empty")
	}

	if r.Spec.DateStart != nil && r.Spec.DateEnd != nil {
		if !r.Spec.DateEnd.After(r.Spec.DateStart.Time) {
			return errors.New("dateEnd must be after dateStart")
		}
	} else if r.Spec.DateStart != nil && r.Spec.DateEnd == nil {
		return errors.New("dateEnd must be provided if dateStart is specified")
	}

	if err := r.validateFeedGroups(); err != nil {
		return err
	}

	return nil
}

// validateFeedGroups checks if the feedGroups exist in the specified ConfigMap.
func (r *HotNews) validateFeedGroups() error {
	configMapName := "hotnews-feeds-group"
	configMapNamespace := "news-aggregator-namespace"

	var cm v1.ConfigMap
	if err := k8sClient.Get(context.TODO(), client.ObjectKey{
		Namespace: configMapNamespace,
		Name:      configMapName,
	}, &cm); err != nil {
		return fmt.Errorf("failed to retrieve ConfigMap %s/%s: %v", configMapNamespace, configMapName, err)
	}

	for _, group := range r.Spec.FeedGroups {
		if _, exists := cm.Data[group]; !exists {
			return fmt.Errorf("feedGroup %s does not exist in ConfigMap %s/%s", group, configMapNamespace, configMapName)
		}
	}

	return nil
}
