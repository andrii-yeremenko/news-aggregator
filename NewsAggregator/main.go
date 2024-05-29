package main

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/aggregator/filter"
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/parser"
	"NewsAggregator/storage"
	"flag"
	"fmt"
	"strings"
)

func main() {
	var sourceArgument, keywordsArgument, startDateArgument, endDateArgument string

	parserSelector := parser.NewParserFactory()
	newsAggregator := aggregator.New(parserSelector)
	resourceLoader := storage.NewLoader(newsAggregator)

	flag.StringVar(&sourceArgument, "sources", "", "Comma-separated list of news sourceArgument\n"+
		"Available sources: "+resourceLoader.GetAvailableSources())
	flag.StringVar(&keywordsArgument, "keywords", "", "Comma-separated list of keywordsArgument"+
		" to filter news articles")
	flag.StringVar(&startDateArgument, "date-start", "", "Start date for filtering news articles"+
		" (format: yyyy-mm-dd)")
	flag.StringVar(&endDateArgument, "date-end", "", "End date for filtering news articles"+
		" (format: yyyy-mm-dd)")
	flag.Parse()

	flagCount := 0
	flag.Visit(func(f *flag.Flag) {
		flagCount++
	})

	// Check if the user has provided the correct number of flags
	if flagCount > 4 {
		flag.Usage()
		return
	}

	if flagCount == 0 {
		resourceLoader.LoadAllResources()
		printArticles(newsAggregator.GetAllArticles())
		return
	}

	if sourceArgument == "" {
		resourceLoader.LoadAllResources()
	} else {
		sources := strings.Split(sourceArgument, ",")
		resourceLoader.LoadSelectedResources(sources)
		newsAggregator.AddFilter(filter.NewSourceFilter(sources))
	}

	if startDateArgument != "" {
		newsAggregator.AddFilter(filter.NewStartDateFilter(startDateArgument))
	}

	if endDateArgument != "" {
		newsAggregator.AddFilter(filter.NewEndDateFilter(endDateArgument))
	}

	if keywordsArgument != "" {
		keywords := strings.Split(keywordsArgument, ",")
		newsAggregator.AddFilter(filter.NewKeywordFilter(keywords))
	}

	filteredArticles := newsAggregator.GetFilteredArticles()

	printArticles(filteredArticles)
}

func printArticles(articles []article.Article) {
	for _, art := range articles {
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("Title: %s\n", art.Title())
		fmt.Printf("Description: %s\n", art.Description())
		fmt.Printf("Date: %s\n", art.Date().HumanReadableString())
		fmt.Printf("Source: %s\n", art.Source())
		fmt.Printf("Author: %s\n", art.Author())
	}
}
