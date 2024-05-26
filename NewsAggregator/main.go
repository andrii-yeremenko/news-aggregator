package main

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/aggregator/parser"
	"NewsAggregator/repository"
	"flag"
	"fmt"
	"strings"
)

func main() {
	var sourceArgument, keywordsArgument, startDateArgument, endDateArgument string
	flag.StringVar(&sourceArgument, "sources", "", "Comma-separated list of news sourceArgument")
	flag.StringVar(&keywordsArgument, "keywords", "", "Comma-separated list of keywordsArgument")
	flag.StringVar(&startDateArgument, "date-start", "", "Start date for filtering news articles")
	flag.StringVar(&endDateArgument, "date-end", "", "End date for filtering news articles")
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

	parserSelector := parser.NewParserFactory()

	newsAggregator := aggregator.NewAggregator(parserSelector)

	resourceLoader := repository.NewResourceLoader(newsAggregator)

	filterBuilder := aggregator.NewFilterBuilder()

	if sourceArgument == "" {
		resourceLoader.LoadAllResources()
	} else {
		sources := strings.Split(sourceArgument, ",")
		resourceLoader.LoadSelectedResources(sources)
		filterBuilder.WithSources(sources)
	}

	if startDateArgument != "" {
		filterBuilder.WithStartDate(startDateArgument)
	}

	if endDateArgument != "" {
		filterBuilder.WithEndDate(endDateArgument)
	}

	if keywordsArgument != "" {
		keywords := strings.Split(keywordsArgument, ",")
		filterBuilder.WithKeywords(keywords)
	}

	filter := filterBuilder.Build()

	filteredArticles := filter.Apply(newsAggregator.GetAllArticles())

	for _, filteredArticle := range filteredArticles {
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("Title: %s\n", filteredArticle.Title())
		fmt.Printf("Description: %s\n", filteredArticle.Description())
		fmt.Printf("Date: %s\n", filteredArticle.Date().HumanReadableString())
		fmt.Printf("Source: %s\n", filteredArticle.Source())
		fmt.Printf("Author: %s\n", filteredArticle.Author())
	}
}
