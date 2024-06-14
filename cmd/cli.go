package cmd

import (
	"flag"
	"fmt"
	"news-aggregator/aggregator"
	"news-aggregator/aggregator/filter"
	"news-aggregator/aggregator/model/article"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/logger"
	"news-aggregator/storage"
	"os"
	"path"
	"strings"
)

// CLI is the command line interface for the news aggregator.
type CLI struct {
	sourceArg     string
	keywordsArg   string
	startDateArg  string
	endDateArg    string
	parserFactory *aggregator.ParserFactory
	aggregator    *aggregator.Aggregator
	storage       *storage.Storage
}

// New creates a new CLI instance.
func New() (*CLI, error) {
	parserPool := aggregator.NewParserFactory()
	a, err := aggregator.New(parserPool)
	if err != nil {
		return nil, err
	}

	basePath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	store := storage.New(path.Join(basePath, "/storage"))

	return &CLI{
		parserFactory: parserPool,
		aggregator:    a,
		storage:       store,
	}, nil
}

// ParseFlags parses the command line flags for sources, keywords, and date ranges.
// Use cases include filtering news by specific sources, topics, or timeframes.
// Errors may occur due to incorrect date formats, unrecognized sources, or unknown arguments.
func (cli *CLI) ParseFlags() {
	flag.StringVar(&cli.sourceArg, "sources", "", "Comma-separated list of news sources\n"+
		"Available sources: "+cli.storage.GetAvailableSources())
	flag.StringVar(&cli.keywordsArg, "keywords", "", "Comma-separated list of keywords to filter news articles")
	flag.StringVar(&cli.startDateArg, "date-start", "", "Start date for filtering news articles (format: yyyy-dd-mm)")
	flag.StringVar(&cli.endDateArg, "date-end", "", "End date for filtering news articles (format: yyyy-dd-mm)")
	flag.Usage = cli.printUsage
	flag.Parse()
}

// Run executes the CLI application.
// This CLI will print the articles to the console based on the provided flags.
func (cli *CLI) Run() error {
	if cli.checkAvailableSources() {
		return fmt.Errorf("no sources available")
	}

	flagCount := cli.countFlags()
	if flagCount > 4 {
		cli.printUsage()
		return fmt.Errorf("too many flags provided")
	}

	if flagCount == 0 {
		cli.showAllArticles()
	} else {
		cli.showFilteredArticles()
	}

	return nil
}

func (cli *CLI) checkAvailableSources() bool {
	if cli.storage.GetAvailableSources() == "" {
		logger.New().Warn("No sources available")
		return true
	}
	return false
}

func (cli *CLI) countFlags() int {
	flagCount := 0
	flag.Visit(func(f *flag.Flag) {
		flagCount++
	})
	return flagCount
}

func (cli *CLI) showAllArticles() {
	resources, err := cli.storage.GetAllResources()
	if err != nil {
		logger.New().Error(err.Error())
	}

	articles, err := cli.aggregator.AggregateMultiple(resources)
	if err != nil {
		logger.New().Error(err.Error())
	}

	cli.printArticles(articles)
}

func (cli *CLI) showFilteredArticles() {
	resources := cli.getResources()

	cli.applyFilters()

	filteredArticles, err := cli.aggregator.AggregateMultiple(resources)

	if err != nil {
		logger.New().Error(err.Error())
	}

	cli.printArticles(filteredArticles)
}

func (cli *CLI) getResources() []resource.Resource {
	if cli.sourceArg == "" {
		resources, err := cli.storage.GetAllResources()
		if err != nil {
			logger.New().Error(err.Error())
		}
		return resources
	}

	sources := strings.Split(cli.sourceArg, ",")
	resources, err := cli.storage.GetSelectedResources(sources)
	if err != nil {
		logger.New().Error(err.Error())
	}
	cli.aggregator.AddFilter(filter.NewSourceFilter(sources))
	return resources
}

func (cli *CLI) applyFilters() {
	if cli.startDateArg != "" {
		startDateFilter, err := filter.NewStartDateFilter(cli.startDateArg)

		if err != nil {
			logger.New().Error(err.Error())
		}

		cli.aggregator.AddFilter(startDateFilter)
	}

	if cli.endDateArg != "" {
		endDateFilter, err := filter.NewEndDateFilter(cli.endDateArg)

		if err != nil {
			logger.New().Error(err.Error())
		}

		cli.aggregator.AddFilter(endDateFilter)
	}

	if cli.keywordsArg != "" {
		keywords := strings.Split(cli.keywordsArg, ",")
		cli.aggregator.AddFilter(filter.NewKeywordFilter(keywords))
	}
}

func (cli *CLI) printArticles(articles []article.Article) {

	params := logger.FilterParams{
		SourceArg:    cli.sourceArg,
		KeywordsArg:  cli.keywordsArg,
		StartDateArg: cli.startDateArg,
		EndDateArg:   cli.endDateArg,
	}

	err := logger.New().PrintArticlesInTemplate(articles, params, "logger/template/article_template.txt")

	if err != nil {
		logger.New().Error(err.Error())
		return
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage: NewsAggregator [options]")
	fmt.Println("If no options are provided, all available articles will be printed.")
	fmt.Println("If any option is provided, only filtered articles will be printed.")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nYou can use multiple flags in any order. Example usage:")
	fmt.Println("  NewsAggregator -sources=source1,source2 -keywords=keyword1,keyword2 -date-start=2024-01-01")
	fmt.Println("  NewsAggregator -keywords=keyword1,keyword2")
	fmt.Println("  NewsAggregator -date-start=2024-01-01 -date-end=2024-12-31")
}
