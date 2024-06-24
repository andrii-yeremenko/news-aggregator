package handler

import (
	"net/http"
	"news-aggregator/aggregator"
	"news-aggregator/aggregator/filter"
	"news-aggregator/aggregator/model/article"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/console_printer"
	"news-aggregator/resource_manager"
	"strings"
)

// NewsAggregatorHandler a Handler for aggregating news by provided filters and arguments.
type NewsAggregatorHandler struct {
	Aggregator      *aggregator.Aggregator
	ResourceManager *resource_manager.ResourceManager
}

// NewNewsHandler creates a new NewsAggregatorHandler instance.
func NewNewsHandler(aggregator *aggregator.Aggregator, resourceManager *resource_manager.ResourceManager) *NewsAggregatorHandler {
	return &NewsAggregatorHandler{
		Aggregator:      aggregator,
		ResourceManager: resourceManager,
	}
}

// Handle is responsible for handling the request and response for the news aggregator.
func (h *NewsAggregatorHandler) Handle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	sources := query.Get("sources")
	keywords := query.Get("keywords")
	startDate := query.Get("date-start")
	endDate := query.Get("date-end")
	sortOrder := query.Get("sort-order")

	resources := h.getResources(sources)
	err := h.applyFilters(keywords, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	articles, err := h.Aggregator.AggregateMultiple(resources)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	articles = h.sortArticles(articles, sortOrder)
	h.printArticles(w, articles)
}

func (h *NewsAggregatorHandler) getResources(sources string) []resource.Resource {
	if sources == "" {
		resources, err := h.ResourceManager.AllResources()
		if err != nil {
			console_printer.New().Error(err.Error())
		}

		return resources
	}

	sourceList := strings.Split(sources, ",")
	resources, err := h.ResourceManager.GetSelectedResources(sourceList)
	if err != nil {
		console_printer.New().Error(err.Error())
	}

	h.Aggregator.AddFilter(filter.NewSourceFilter(sourceList))
	return resources
}

func (h *NewsAggregatorHandler) applyFilters(keywords, startDate, endDate string) error {
	if startDate != "" {
		startDateFilter, err := filter.NewStartDateFilter(startDate)
		if err != nil {
			return err
		}
		h.Aggregator.AddFilter(startDateFilter)
	}

	if endDate != "" {
		endDateFilter, err := filter.NewEndDateFilter(endDate)
		if err != nil {
			return err
		}
		h.Aggregator.AddFilter(endDateFilter)
	}

	if keywords != "" {
		keywordList := strings.Split(keywords, ",")
		h.Aggregator.AddFilter(filter.NewKeywordFilter(keywordList))
	}

	return nil
}

func (h *NewsAggregatorHandler) sortArticles(articles []article.Article, sortOrder string) []article.Article {
	switch sortOrder {
	case "asc":
		return aggregator.SortArticlesByDateAsc(articles)
	case "desc":
		return aggregator.SortArticlesByDateDesc(articles)
	default:
		console_printer.New().Error("Unknown sort order")
		return articles
	}
}

func (h *NewsAggregatorHandler) printArticles(w http.ResponseWriter, articles []article.Article) {
	params := console_printer.FilterParams{
		SourceArg:    "",
		KeywordsArg:  "",
		StartDateArg: "",
		EndDateArg:   "",
		OrderArg:     "",
	}

	err := console_printer.New().PrintArticles(articles, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
