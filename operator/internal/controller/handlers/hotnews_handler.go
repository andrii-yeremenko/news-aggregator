package handlers

import (
	"context"
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	newsaggregatorv1 "com.teamdev/news-aggregator/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// ReconcileAllHotNews reconciles all HotNews resources.
func ReconcileAllHotNews(ctx context.Context, obj client.Object, c client.Client, namespace string, configMapName string, configMapNamespace string) []reconcile.Request {
	var hotNewsList newsaggregatorv1.HotNewsList

	configMap, ok := obj.(*v1.ConfigMap)
	if !ok {
		log.Log.Error(fmt.Errorf("ConfigMap object is not of type v1.ConfigMap"), "Failed to reconcile all HotNews resources")
		return nil
	}

	if configMap.Name != configMapName || configMap.Namespace != configMapNamespace {
		log.Log.Info("ConfigMap does not match the expected ConfigMap. Skipping reconciliation of HotNews resources")
		return nil
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	log.Log.Info("Reconciling all HotNews resources")

	if err := c.List(timeoutCtx, &hotNewsList, client.InNamespace(namespace)); err != nil {
		log.Log.Error(err, "Failed to list HotNews resources")
		return nil
	}

	var requests []ctrl.Request
	for _, hotNews := range hotNewsList.Items {
		if hotnewsUsingConfigMap(&hotNews) {
			requests = append(requests, enqueueHotNewsResource(&hotNews))
		}
	}

	return requests
}

// hotnewsUsingConfigMap checks if the HotNews resource is using FeedGroups.
func hotnewsUsingConfigMap(hotNews *newsaggregatorv1.HotNews) bool {
	return hotNews.Spec.FeedGroups != nil
}

// enqueueHotNewsResource enqueues the HotNews resource for reconciliation.
func enqueueHotNewsResource(hotNews *newsaggregatorv1.HotNews) ctrl.Request {
	log.Log.Info("Enqueueing HotNews resource", "HotNews", hotNews.Name)

	return reconcile.Request{
		NamespacedName: client.ObjectKey{
			Name:      hotNews.Name,
			Namespace: hotNews.Namespace,
		},
	}
}
