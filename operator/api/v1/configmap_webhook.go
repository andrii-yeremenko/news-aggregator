package v1

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
)

// +kubebuilder:webhook:path=/validate--v1-configmap,mutating=false,failurePolicy=fail,sideEffects=None,groups="",resources=configmaps,verbs=create;update,versions=v1,name=vconfigmap.kb.io,admissionReviewVersions=v1

// ConfigMapWebhook handles validation of ConfigMaps
type ConfigMapWebhook struct {
	client.Client
	ConfigMapName      string
	ConfigMapNamespace string
}

// ValidateCreate performs validation for ConfigMap creation
func (r *ConfigMapWebhook) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("expected a ConfigMap but got a %T", obj)
	}

	if cm.Namespace != r.ConfigMapNamespace || cm.Name != r.ConfigMapName {
		return nil, nil
	}

	if err := r.validateConfigMap(ctx, cm); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateUpdate performs validation for ConfigMap updates
func (r *ConfigMapWebhook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (warnings admission.Warnings, err error) {
	newCM, ok := newObj.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("expected a ConfigMap but got a %T", newObj)
	}

	if err := r.validateConfigMap(ctx, newCM); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateDelete performs validation for ConfigMap deletion
func (r *ConfigMapWebhook) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// validateConfigMap checks if any fields of ConfigMap are empty and if feeds exist
func (r *ConfigMapWebhook) validateConfigMap(ctx context.Context, cm *corev1.ConfigMap) error {
	for key, value := range cm.Data {
		if value == "" {
			return fmt.Errorf("data field '%s' is empty", key)
		}

		if err := r.validateFeeds(ctx, cm.Namespace, value); err != nil {
			return fmt.Errorf("validation error in field '%s': %v", key, err)
		}
	}

	return nil
}

// validateFeeds checks if the feeds listed in the ConfigMap exist as Kubernetes resources
func (r *ConfigMapWebhook) validateFeeds(ctx context.Context, namespace, feeds string) error {
	feedList := strings.Split(feeds, ",")
	for i := range feedList {
		feedList[i] = strings.TrimSpace(feedList[i])
	}

	availableFeeds, err := r.listAvailableFeeds(ctx, namespace)
	if err != nil {
		return fmt.Errorf("failed to list feeds in namespace %s: %w", namespace, err)
	}

	var notFoundFeeds []string
	for _, feed := range feedList {
		if !availableFeeds[feed] {
			notFoundFeeds = append(notFoundFeeds, feed)
		}
	}

	if len(notFoundFeeds) > 0 {
		errStr := strings.Join(notFoundFeeds, ", ")
		return fmt.Errorf("feeds \"%s\" do not exist in ", errStr)
	}

	return nil
}

// listAvailableFeeds lists all available feeds in the given namespace and returns a set of their names
func (r *ConfigMapWebhook) listAvailableFeeds(ctx context.Context, namespace string) (map[string]bool, error) {
	feedSet := make(map[string]bool)

	var feedList FeedList
	if err := r.Client.List(ctx, &feedList, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	for _, feed := range feedList.Items {
		feedSet[feed.Name] = true
	}

	return feedSet, nil
}
