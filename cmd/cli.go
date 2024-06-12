package cmd

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/aggregator/filter"
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/aggregator/model/resource"
	"NewsAggregator/logger"
	"NewsAggregator/storage"
	"flag"
	"fmt"
	"os"
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
func New() *CLI {
	fact := aggregator.NewParserFactory()
	agg, err := aggregator.New(fact)
	handleError(err)

	basePath, err := os.Getwd()
	handleError(err)

	store := storage.New(basePath + "/storage")

	return &CLI{
		parserFactory: fact,
		aggregator:    agg,
		storage:       store,
	}
}

// ParseFlags parses the command line flags.
func (cli *CLI) ParseFlags() {
	flag.StringVar(&cli.sourceArg, "sources", "", "Comma-separated list of news sources\n"+
		"Available sources: "+cli.storage.GetAvailableSources())
	flag.StringVar(&cli.keywordsArg, "keywords", "", "Comma-separated list of keywords to filter news articles")
	flag.StringVar(&cli.startDateArg, "date-start", "", "Start date for filtering news articles (format: yyyy-dd-mm)")
	flag.StringVar(&cli.endDateArg, "date-end", "", "End date for filtering news articles (format: yyyy-dd-mm)")
	flag.Usage = cli.printUsage
	flag.Parse()
}

// Run runs the CLI.
func (cli *CLI) Run() {
	if cli.noSourcesAvailable() {
		return
	}

	flagCount := cli.countFlags()
	if flagCount > 4 {
		cli.printUsage()
		return
	}

	if flagCount == 0 {
		cli.showAllArticles()
	} else {
		cli.showFilteredArticles()
	}
}

func (cli *CLI) noSourcesAvailable() bool {
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
	handleError(err)

	articles, err := cli.aggregator.AggregateMultiple(resources)
	handleError(err)

	cli.printArticles(articles)
}

func (cli *CLI) showFilteredArticles() {
	resources := cli.getResources()

	cli.applyFilters()

	filteredArticles, err := cli.aggregator.AggregateMultiple(resources)
	handleError(err)

	cli.printArticles(filteredArticles)
}

func (cli *CLI) getResources() []resource.Resource {
	if cli.sourceArg == "" {
		resources, err := cli.storage.GetAllResources()
		handleError(err)
		return resources
	}

	sources := strings.Split(cli.sourceArg, ",")
	resources, err := cli.storage.GetSelectedResources(sources)
	handleError(err)

	cli.aggregator.AddFilter(filter.NewSourceFilter(sources))
	return resources
}

func (cli *CLI) applyFilters() {
	if cli.startDateArg != "" {
		startDateFilter, err := filter.NewStartDateFilter(cli.startDateArg)
		handleError(err)
		cli.aggregator.AddFilter(startDateFilter)
	}

	if cli.endDateArg != "" {
		endDateFilter, err := filter.NewEndDateFilter(cli.endDateArg)
		handleError(err)
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

func handleError(err error) {
	if err != nil {
		logger.New().Error(err.Error())
	}
}
