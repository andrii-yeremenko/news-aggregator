package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"news-aggregator/aggregator"
	"news-aggregator/aggregator/filter"
	"news-aggregator/aggregator/model/article"
	"news-aggregator/aggregator/model/resource"
	"strings"
)

// NewsAggregatorHandler a Handler for aggregating news by provided filters and arguments.
type NewsAggregatorHandler struct {
	resourceManager ResourceManager
	parserPool      *aggregator.ParserFactory
}

// NewNewsHandler creates a new NewsAggregatorHandler instance.
func NewNewsHandler(resourceManager ResourceManager) *NewsAggregatorHandler {
	return &NewsAggregatorHandler{
		resourceManager: resourceManager,
		parserPool:      aggregator.NewParserFactory(),
	}
}

// Handle is responsible for handling the request and response for the news aggregator.
func (h *NewsAggregatorHandler) Handle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	sources := query.Get("sources")
	keywords := query.Get("keywords")
	startDate := query.Get("date-start")
	endDate := query.Get("date-end")
	sortOrder := query.Get("sort-order")

	a, err := aggregator.New(h.parserPool)
	if err != nil {
		log.Fatalf("failed to create aggregator: %v", err)
	}

	resources, err := h.getResources(sources)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.applyFilters(a, keywords, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	articles, err := a.AggregateMultiple(resources)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if sortOrder != "" {
		articles, err = h.sortArticles(articles, sortOrder)
		if err != nil {
			http.Error(w, "Invalid sort order", http.StatusBadRequest)
			return
		}
	}

	h.sendArticles(w, articles)
}

func (h *NewsAggregatorHandler) getResources(sources string) ([]resource.Resource, error) {
	if sources == "" {
		resources, err := h.resourceManager.GetAllResources()
		if err != nil {
			return resources, err
		}
		return resources, nil
	}

	sourceList := strings.Split(sources, ",")
	resources, err := h.resourceManager.GetSelectedResources(sourceList)
	if err != nil {
		return resources, err
	}

	return resources, nil
}

func (h *NewsAggregatorHandler) applyFilters(a *aggregator.Aggregator, keywords, startDate, endDate string) error {
	if startDate != "" {
		startDateFilter, err := filter.NewStartDateFilter(startDate)
		if err != nil {
			return err
		}
		a.AddFilter(startDateFilter)
	}

	if endDate != "" {
		endDateFilter, err := filter.NewEndDateFilter(endDate)
		if err != nil {
			return err
		}
		a.AddFilter(endDateFilter)
	}

	if keywords != "" {
		keywordList := strings.Split(keywords, ",")
		a.AddFilter(filter.NewKeywordFilter(keywordList))
	}

	return nil
}

func (h *NewsAggregatorHandler) sortArticles(articles []article.Article, sortOrder string) ([]article.Article, error) {
	switch sortOrder {
	case "asc":
		return aggregator.SortArticlesByDateAsc(articles), nil
	case "desc":
		return aggregator.SortArticlesByDateDesc(articles), nil
	default:
		return nil, errors.New("invalid sort order")
	}
}

func (h *NewsAggregatorHandler) sendArticles(w http.ResponseWriter, articles []article.Article) {
	var articlesJSON []map[string]interface{}

	for _, art := range articles {
		articleJSON := map[string]interface{}{
			"title":        art.TitleStr(),
			"description":  art.DescriptionStr(),
			"creationDate": art.Date().HumanReadableString(),
			"source":       art.Source(),
			"author":       art.Author(),
			"link":         art.Link(),
		}

		articlesJSON = append(articlesJSON, articleJSON)
	}

	w.Header().Set("Content-Type", "application/json")

	responseJSON, err := json.Marshal(articlesJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
