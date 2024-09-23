package controller

import (
	"com.teamdev/news-aggregator/internal/controller/predicates"
	"context"
	"encoding/json"
	"fmt"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sort"
	"strings"
	"time"

	newsaggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const newsEndpoint = "/news"

// HotNewsReconciler reconciles a HotNews object
type HotNewsReconciler struct {
	client.Client
	HTTPClient         HTTPClient
	Scheme             *runtime.Scheme
	Finalizer          string
	NewsAggregatorURL  string
	Namespace          string
	ConfigMapName      string
	ConfigMapNamespace string
}

// Article represents a news article.
type Article struct {
	Title string `json:"title"`
}

// ArticlesResponse represents a response from the News Aggregator HTTPS server.
type ArticlesResponse []Article

// +kubebuilder:rbac:groups=news-aggregator.com.teamdev,resources=hotnews,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=news-aggregator.com.teamdev,resources=hotnews/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=news-aggregator.com.teamdev,resources=hotnews/finalizers,verbs=update
// +kubebuilder:rbac:groups=news-aggregator.com.teamdev,resources=feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile reconciles the HotNews object.
// It fetches news articles from the News Aggregator service based on the HotNews spec input,
// and updates the HotNews status with the fetched articles.
func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Starting reconciliation", "HotNews", req.NamespacedName)

	var hotNews newsaggregatorv1.HotNews
	if err := r.Get(ctx, req.NamespacedName, &hotNews); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("HotNews resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !controllerutil.ContainsFinalizer(&hotNews, r.Finalizer) {
		controllerutil.AddFinalizer(&hotNews, r.Finalizer)
		if err := r.Update(ctx, &hotNews); err != nil {
			return ctrl.Result{}, err
		}
	}

	if hotNews.DeletionTimestamp != nil {
		if err := r.finalizeHotNews(ctx, &hotNews); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if err := r.addOwnerReferencesOnFeeds(ctx, &hotNews); err != nil {
		logger.Error(err, "Failed to set owner references on related feeds")
		return ctrl.Result{}, err
	}

	url, err := r.buildRequestURL(hotNews.Spec)
	if err != nil {
		logger.Error(err, "Failed to build request URL")
		statusErr := r.updateStatus(&hotNews, newsaggregatorv1.ConditionFailed)
		if statusErr != nil {
			logger.Error(statusErr, "Failed to update HotNews status")
		}
		return ctrl.Result{}, err
	}

	titles, err := r.fetchNews(url)
	if err != nil {
		logger.Error(err, "Failed to fetch news")
		statusErr := r.updateStatus(&hotNews, newsaggregatorv1.ConditionFailed)
		if statusErr != nil {
			logger.Error(statusErr, "Failed to update HotNews status")
		}
		return ctrl.Result{}, err
	}

	if err := r.updateHotNewsArticles(ctx, &hotNews, titles, url); err != nil {
		logger.Error(err, "Failed to update HotNews status")
		statusErr := r.updateStatus(&hotNews, newsaggregatorv1.ConditionFailed)
		if statusErr != nil {
			logger.Error(statusErr, "Failed to update HotNews status")
		}
		return ctrl.Result{}, err
	}

	logger.Info("Successfully reconciled HotNews", "HotNews", hotNews.Name, "ArticlesCount",
		hotNews.Status.ArticlesCount)
	return ctrl.Result{}, r.updateStatus(&hotNews, newsaggregatorv1.ConditionUpdated)
}

// finalizeHotNews handles the finalizer logic for the HotNews resource.
func (r *HotNewsReconciler) finalizeHotNews(ctx context.Context, hotNews *newsaggregatorv1.HotNews) error {
	logger := log.FromContext(ctx)

	if err := r.removeOwnerReferencesFromFeeds(ctx, hotNews); err != nil {
		return err
	}

	controllerutil.RemoveFinalizer(hotNews, r.Finalizer)
	if err := r.Update(ctx, hotNews); err != nil {
		logger.Error(err, "Failed to remove finalizer from HotNews")
		return err
	}

	logger.Info("Finalizer successfully handled and removed", "HotNews", hotNews.Name)
	return nil
}

// removeOwnerReferencesFromFeeds removes the owner references from all related Feed resources.
func (r *HotNewsReconciler) removeOwnerReferencesFromFeeds(ctx context.Context, hotNews *newsaggregatorv1.HotNews) error {

	allFeeds, err := r.getAllFeeds(ctx, hotNews.Spec)

	if err != nil {
		return fmt.Errorf("failed to get all feeds: %w", err)
	}

	for _, feedName := range allFeeds {
		var feed newsaggregatorv1.Feed
		if err := r.Get(ctx, client.ObjectKey{Name: feedName, Namespace: hotNews.Namespace}, &feed); err != nil {
			if errors.IsNotFound(err) {
				log.Log.Info("Feed resource not found. Skipping", "Feed", feedName)
				continue
			}
			return err
		}

		currentOwnerRef := metav1.NewControllerRef(hotNews, newsaggregatorv1.GroupVersion.WithKind("HotNews"))

		allOwnerRefs := feed.GetOwnerReferences()
		var newOwnerRefs []metav1.OwnerReference

		for _, ownerRef := range allOwnerRefs {
			if ownerRef.Name != currentOwnerRef.Name {
				newOwnerRefs = append(newOwnerRefs, ownerRef)
			}
		}

		feed.SetOwnerReferences(newOwnerRefs)

		if err := r.Update(ctx, &feed); err != nil {
			return fmt.Errorf("failed to remove owner reference from Feed: %w", err)
		}
	}

	return nil
}

// addOwnerReferencesOnFeeds sets the owner reference on all related Feed resources
func (r *HotNewsReconciler) addOwnerReferencesOnFeeds(ctx context.Context, hotNews *newsaggregatorv1.HotNews) error {

	allFeeds, err := r.getAllFeeds(context.Background(), hotNews.Spec)

	if err != nil {
		return fmt.Errorf("failed to get all feeds: %w", err)
	}

	for _, feedName := range allFeeds {
		var feed newsaggregatorv1.Feed
		if err := r.Get(ctx, client.ObjectKey{Name: feedName, Namespace: hotNews.Namespace}, &feed); err != nil {
			if errors.IsNotFound(err) {
				log.Log.Info("Feed resource not found. Skipping", "Feed", feedName)
				continue
			}
			return err
		}

		allOwnerRefs := feed.GetOwnerReferences()
		ownerRef := metav1.NewControllerRef(hotNews, newsaggregatorv1.GroupVersion.WithKind("HotNews"))

		allOwnerRefs = append(allOwnerRefs, *ownerRef)

		feed.SetOwnerReferences(allOwnerRefs)

		if err := r.Update(ctx, &feed); err != nil {
			return fmt.Errorf("failed to update Feed with owner reference: %w", err)
		}
	}
	return nil
}

// buildRequestURL builds the URL for the request to the News Aggregator service.
func (r *HotNewsReconciler) buildRequestURL(spec newsaggregatorv1.HotNewsSpec) (string, error) {

	baseURL := fmt.Sprintf("%s%s", r.NewsAggregatorURL, newsEndpoint)

	allFeeds, err := r.getAllFeeds(context.Background(), spec)

	if err != nil {
		return "", fmt.Errorf("failed to get all feeds: %w", err)
	}

	var params []string
	if len(spec.Keywords) > 0 {
		params = append(params, "keywords="+strings.Join(spec.Keywords, ","))
	}
	if len(allFeeds) > 0 {
		params = append(params, "sources="+strings.Join(allFeeds, ","))
	}
	if spec.DateStart != nil {
		params = append(params, "date-start="+formatDateForURL(spec.DateStart.Time))
	}
	if spec.DateEnd != nil {
		params = append(params, "date-end="+formatDateForURL(spec.DateEnd.Time))
	}

	url := fmt.Sprintf("%s?%s", baseURL, strings.Join(params, "&"))
	log.Log.Info("Built request URL", "URL", url)
	return url, nil
}

// getAllFeeds returns all feed sources from the HotNews spec and feed groups.
func (r *HotNewsReconciler) getAllFeeds(ctx context.Context, hotNews newsaggregatorv1.HotNewsSpec) ([]string, error) {

	if hotNews.FeedGroups != nil {
		feeds, err := r.getFeedSourcesFromConfigMap(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get feed sources: %w", err)
		}

		return r.collectUniqueSources(hotNews, feeds), nil
	}

	return hotNews.Feeds, nil
}

// collectUniqueSources collects unique sources from the HotNews spec and feed groups.
func (r *HotNewsReconciler) collectUniqueSources(spec newsaggregatorv1.HotNewsSpec, sources map[string][]string) []string {
	uniqueSources := make(map[string]struct{})

	for _, feed := range spec.Feeds {
		uniqueSources[feed] = struct{}{}
	}

	for _, group := range spec.FeedGroups {
		for _, feed := range sources[group] {
			uniqueSources[feed] = struct{}{}
		}
	}

	allFeeds := make([]string, 0, len(uniqueSources))
	for feed := range uniqueSources {
		allFeeds = append(allFeeds, feed)
	}

	sort.Strings(allFeeds)
	return allFeeds
}

// formatDateForURL formats the time to the required format for the URL.
func formatDateForURL(t time.Time) string {
	return t.Format("2006-02-01")
}

// fetchNews fetches news articles from the given URL.
func (r *HotNewsReconciler) fetchNews(url string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if closeErr := Body.Close(); closeErr != nil {
			log.Log.Error(closeErr, "Error closing response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	var articles ArticlesResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&articles); err != nil {
		return nil, err
	}

	var titles []string
	for _, article := range articles {
		titles = append(titles, article.Title)
	}

	return titles, nil
}
func (r *HotNewsReconciler) updateHotNewsArticles(ctx context.Context, hotNews *newsaggregatorv1.HotNews, titles []string,
	url string) error {
	actualFeedCount := len(titles)
	titlesCount := hotNews.Spec.SummaryConfig.TitlesCount
	if len(titles) > titlesCount {
		actualFeedCount = titlesCount
		titles = titles[:titlesCount]
	}

	hotNews.Status = newsaggregatorv1.HotNewsStatus{
		ArticlesCount:  actualFeedCount,
		NewsLink:       url,
		ArticlesTitles: titles,
	}

	if err := r.Client.Status().Update(ctx, hotNews); err != nil {
		return fmt.Errorf("failed to update HotNews status: %w", err)
	}
	return nil
}

func (r *HotNewsReconciler) getFeedSourcesFromConfigMap(ctx context.Context) (map[string][]string, error) {
	var configMap v1.ConfigMap
	if err := r.Get(ctx, client.ObjectKey{
		Name:      r.ConfigMapName,
		Namespace: r.ConfigMapNamespace,
	}, &configMap); err != nil {
		return nil, fmt.Errorf("failed to get ConfigMap: %w", err)
	}

	sources := make(map[string][]string)
	for group, feeds := range configMap.Data {
		sources[group] = strings.Split(feeds, ",")
	}

	return sources, nil
}

// updateStatus updates the status of the HotNews object.
func (r *HotNewsReconciler) updateStatus(hotnews *newsaggregatorv1.HotNews, conditionType newsaggregatorv1.ConditionType) error {

	if hotnews.DeletionTimestamp != nil {
		return nil
	}

	hotnews.Status.Conditions = append(hotnews.Status.Conditions, conditionType)
	return r.Client.Status().Update(context.Background(), hotnews)
}

// SetupWithManager sets up the controller with the Manager.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {

	feedPredicate := predicates.NewFeedPredicate(r.Namespace)
	configMapPredicate := predicates.NewConfigMapPredicate(r.ConfigMapNamespace, r.ConfigMapName)

	return ctrl.NewControllerManagedBy(mgr).
		For(&newsaggregatorv1.HotNews{}).
		Owns(&newsaggregatorv1.Feed{}, builder.WithPredicates(feedPredicate)).
		Watches(
			&v1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(r.reconcileAllHotNews),
			builder.WithPredicates(configMapPredicate),
		).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

// reconcileAllHotNews reconciles all HotNews resources.
func (r *HotNewsReconciler) reconcileAllHotNews(context.Context, client.Object) []reconcile.Request {
	var hotNewsList newsaggregatorv1.HotNewsList

	log.Log.Info("Reconciling all HotNews resources")

	if err := r.Client.List(context.TODO(), &hotNewsList, client.InNamespace(r.Namespace)); err != nil {
		log.Log.Error(err, "Failed to list HotNews resources")
		return nil
	}

	var requests []ctrl.Request
	for _, hotNews := range hotNewsList.Items {

		log.Log.Info("Enqueueing HotNews resource", "HotNews", hotNews.Name)

		requests = append(requests, ctrl.Request{
			NamespacedName: client.ObjectKey{
				Name:      hotNews.Name,
				Namespace: hotNews.Namespace,
			},
		})
	}

	return requests
}
