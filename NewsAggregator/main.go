package main

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/parser"
	"NewsAggregator/repository"
	"flag"
	"fmt"
	"strings"
)

func main() {
	var sourceArgument, keywordsArgument, startDateArgument, endDateArgument string
	flag.StringVar(&sourceArgument, "source", "", "Comma-separated list of news sourceArgument")
	flag.StringVar(&keywordsArgument, "keywords", "", "Comma-separated list of keywordsArgument")
	flag.StringVar(&startDateArgument, "date-start", "", "Start date for filtering news articles")
	flag.StringVar(&endDateArgument, "date-end", "", "End date for filtering news articles")
	flag.Parse()

	flagCount := 0
	flag.Visit(func(f *flag.Flag) {
		flagCount++
	})

	if flagCount > 2 || flagCount < 1 {
		flag.Usage()
		return
	}

	newsRepository := repository.NewInMemoryRepository()

	parserSelector := parser.NewParserFactory()
	parserSelector.RegisterParser("json", "nbc-news.com", &parser.JSONParser{})
	parserSelector.RegisterParser("rss", "abc-news.com", &parser.RSSParser{})
	parserSelector.RegisterParser("rss", "washington-times.com", &parser.RSSParser{})
	parserSelector.RegisterParser("rss", "bbc-world.com", &parser.RSSParser{})
	parserSelector.RegisterParser("html", "usa-today.com", &parser.USATodayHTMLParser{})

	newsAggregator := aggregator.NewAggregator(parserSelector)

	// Error handling is omitted for simplicity
	nbc, _ := newsRepository.ReadFile("nbc-news.com", "json", "repository/news-resources/nbc-news.json")
	_ = newsAggregator.LoadResource(nbc)
	abc, _ := newsRepository.ReadFile("abc-news.com", "rss", "repository/news-resources/abc-news.xml")
	_ = newsAggregator.LoadResource(abc)
	washington, _ := newsRepository.ReadFile("washington-times.com", "rss", "repository/news-resources/washington-times.xml")
	_ = newsAggregator.LoadResource(washington)
	bbc, _ := newsRepository.ReadFile("bbc-world.com", "rss", "repository/news-resources/bbc-world.xml")
	_ = newsAggregator.LoadResource(bbc)
	usaToday, _ := newsRepository.ReadFile("usa-today.com", "html", "repository/news-resources/usa-today-world-news.html")
	_ = newsAggregator.LoadResource(usaToday)

	var filteredArticles []article.Article

	if sourceArgument != "" {
		newsAggregator.FilterBySources([]string{sourceArgument})
	}
	if keywordsArgument != "" {
		keywords := strings.Split(keywordsArgument, ",")
		filteredArticles = newsAggregator.FilterByKeywords(keywords)
	}
	if startDateArgument != "" && endDateArgument != "" {
		filteredArticles = newsAggregator.FilterByDateRange(startDateArgument, endDateArgument)
	}

	for _, filteredArticle := range filteredArticles {
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("Title: %s\n", filteredArticle.Title())
		fmt.Printf("Description: %s\n", filteredArticle.Description())
		fmt.Printf("Date: %s\n", filteredArticle.Date().HumanReadableString())
		fmt.Printf("Source: %s\n", filteredArticle.Source())
		fmt.Printf("Author: %s\n", filteredArticle.Author())
	}
}
