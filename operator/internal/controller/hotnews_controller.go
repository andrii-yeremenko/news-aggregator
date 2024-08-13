package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sort"
	"strings"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	newsaggregatorv1 "com.teamdev/news-aggregator/api/v1"
)

const newsEndpoint = "/news"

// HotNewsReconciler reconciles a HotNews object
type HotNewsReconciler struct {
	client.Client
	HTTPClient         HTTPClient
	Scheme             *runtime.Scheme
	NewsAggregatorURL  string
	ConfigMapName      string
	ConfigMapNamespace string
}

type Article struct {
	Title string `json:"title"`
}

type ArticlesResponse []Article

// +kubebuilder:rbac:groups=news-aggregator.com.teamdev,resources=hotnews,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=news-aggregator.com.teamdev,resources=hotnews/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=news-aggregator.com.teamdev,resources=hotnews/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Starting reconciliation", "HotNews", req.NamespacedName)

	var hotNews newsaggregatorv1.HotNews
	if err := r.Get(ctx, req.NamespacedName, &hotNews); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("HotNews resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get HotNews resource")
		return ctrl.Result{}, err
	}

	url, err := r.buildRequestURL(hotNews.Spec)
	if err != nil {
		logger.Error(err, "Failed to build request URL")
		return ctrl.Result{}, err
	}

	titles, err := r.fetchNews(url)
	if err != nil {
		logger.Error(err, "Failed to fetch news")
		return ctrl.Result{}, err
	}

	if err := r.updateHotNewsStatus(ctx, &hotNews, titles, url); err != nil {
		logger.Error(err, "Failed to update HotNews status")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully reconciled HotNews", "HotNews", hotNews.Name, "ArticlesCount", hotNews.Status.ArticlesCount)
	return ctrl.Result{}, nil
}

func (r *HotNewsReconciler) buildRequestURL(spec newsaggregatorv1.HotNewsSpec) (string, error) {
	if len(spec.Feeds) == 0 && len(spec.FeedGroups) == 0 {
		return "", fmt.Errorf("both Feeds and FeedGroups are empty; at least one must be specified")
	}

	baseURL := fmt.Sprintf("%s%s", r.NewsAggregatorURL, newsEndpoint)

	sources, err := r.getFeedSources(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get feed sources: %w", err)
	}

	uniqueSources := r.collectUniqueSources(spec, sources)

	var params []string
	if len(spec.Keywords) > 0 {
		params = append(params, "keywords="+strings.Join(spec.Keywords, ","))
	}
	if len(uniqueSources) > 0 {
		params = append(params, "sources="+strings.Join(uniqueSources, ","))
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

func (r *HotNewsReconciler) collectUniqueSources(spec newsaggregatorv1.HotNewsSpec, sources map[string][]string) []string {
	uniqueSources := make(map[string]struct{})
	for _, feed := range spec.Feeds {
		uniqueSources[feed] = struct{}{}
	}
	for _, group := range spec.FeedGroups {
		if groupFeeds, ok := sources[group]; ok {
			for _, feed := range groupFeeds {
				uniqueSources[feed] = struct{}{}
			}
		}
	}
	var allFeeds []string
	for feed := range uniqueSources {
		allFeeds = append(allFeeds, feed)
	}
	sort.Strings(allFeeds)
	return allFeeds
}

func formatDateForURL(t time.Time) string {
	return t.Format("2006-01-02")
}

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
func (r *HotNewsReconciler) updateHotNewsStatus(ctx context.Context, hotNews *newsaggregatorv1.HotNews, titles []string, url string) error {
	titlesCount := hotNews.Spec.SummaryConfig.TitlesCount
	if len(titles) > titlesCount {
		titles = titles[:titlesCount]
	}

	hotNews.Status = newsaggregatorv1.HotNewsStatus{
		ArticlesCount:  len(titles),
		NewsLink:       url,
		ArticlesTitles: titles,
	}

	if err := r.Client.Status().Update(ctx, hotNews); err != nil {
		return fmt.Errorf("failed to update HotNews status: %w", err)
	}
	return nil
}

func (r *HotNewsReconciler) getFeedSources(ctx context.Context) (map[string][]string, error) {
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

// SetupWithManager sets up the controller with the Manager.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&newsaggregatorv1.HotNews{}).
		Watches(
			&newsaggregatorv1.Feed{},
			&handler.EnqueueRequestForObject{},
		).
		Watches(
			&v1.ConfigMap{},
			&handler.EnqueueRequestForObject{},
		).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				return true
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				return e.ObjectNew.GetGeneration() != e.ObjectOld.GetGeneration()
			},
		}).
		Complete(r)
}
