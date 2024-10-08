package cli

import (
	"flag"
	"fmt"
	"log"
	"news-aggregator/aggregator"
	"news-aggregator/aggregator/filter"
	"news-aggregator/aggregator/model/article"
	"news-aggregator/aggregator/model/resource"
	"news-aggregator/manager"
	"news-aggregator/print"
	"os"
	"path"
	"strings"
)

// CLI is the command line interface for the news aggregator.
type CLI struct {
	sourceArg       string
	keywordsArg     string
	startDateArg    string
	endDateArg      string
	sortOrderArg    string
	parserFactory   *aggregator.ParserFactory
	aggregator      *aggregator.Aggregator
	resourceManager *manager.ResourceManager
	printer         *print.Logger
}

// New creates a new CLI instance.
func New(managerPath, storagePath string) (*CLI, error) {
	parserPool := aggregator.NewParserFactory()
	a, err := aggregator.New(parserPool)
	if err != nil {
		return nil, err
	}

	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	managerConfigPath := path.Join(basePath, managerPath)
	storagePath = path.Join(basePath, storagePath)

	m, err := manager.New(storagePath, managerConfigPath)

	if err != nil {
		return nil, err
	}

	return &CLI{
		parserFactory:   parserPool,
		aggregator:      a,
		resourceManager: m,
		printer:         print.New(),
	}, nil
}

// ParseFlags parses the command line flags for sources, keywords, and date ranges.
// Use cases include filtering news by specific sources, topics, or timeframes.
// Errors may occur due to incorrect date formats, unrecognized sources, or unknown arguments.
func (cli *CLI) ParseFlags() {
	flag.StringVar(&cli.sourceArg, "sources", "", "Comma-separated list of news sources\n"+
		"Available sources: "+cli.resourceManager.AvailableSources())
	flag.StringVar(&cli.keywordsArg, "keywords", "", "Comma-separated list of keywords to filter news articles")
	flag.StringVar(&cli.startDateArg, "date-start", "", "Start date for filtering news articles (format: yyyy-dd-mm)")
	flag.StringVar(&cli.endDateArg, "date-end", "", "End date for filtering news articles (format: yyyy-dd-mm)")
	flag.StringVar(&cli.sortOrderArg, "sort-order", "asc", "Sort order for articles by date (asc/desc)")
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
	if cli.resourceManager.AvailableSources() == "" {
		cli.printer.Warn("No sources available")
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
	resources, err := cli.resourceManager.GetAllResources()
	if err != nil {
		cli.printer.Error(err.Error())
	}

	articles, err := cli.aggregator.AggregateMultiple(resources)
	if err != nil {
		cli.printer.Error(err.Error())
	}

	articles = cli.sortArticles(articles)
	cli.printArticles(articles)
}

func (cli *CLI) showFilteredArticles() {
	resources := cli.getResources()

	err := cli.applyFilters()
	if err != nil {
		cli.printer.Error(err.Error())
		return
	}

	filteredArticles, err := cli.aggregator.AggregateMultiple(resources)

	if err != nil {
		cli.printer.Error(err.Error())
		return
	}

	filteredArticles = cli.sortArticles(filteredArticles)
	cli.printArticles(filteredArticles)
}

func (cli *CLI) sortArticles(articles []article.Article) []article.Article {

	if cli.sortOrderArg == "asc" {
		return aggregator.SortArticlesByDateAsc(articles)
	} else if cli.sortOrderArg == "desc" {
		return aggregator.SortArticlesByDateDesc(articles)
	} else {
		cli.printer.Error("Unknown sort order")
		return articles
	}
}

func (cli *CLI) getResources() []resource.Resource {
	if cli.sourceArg == "" {
		resources, err := cli.resourceManager.GetAllResources()
		if err != nil {
			cli.printer.Error(err.Error())
		}

		return resources
	}

	sources := strings.Split(cli.sourceArg, ",")
	resources, err := cli.resourceManager.GetSelectedResources(sources)
	if err != nil {
		cli.printer.Error(err.Error())
	}
	cli.aggregator.AddFilter(filter.NewSourceFilter(sources))
	return resources
}

func (cli *CLI) applyFilters() error {
	if cli.startDateArg != "" {
		startDateFilter, err := filter.NewStartDateFilter(cli.startDateArg)

		if err != nil {
			return err
		}

		cli.aggregator.AddFilter(startDateFilter)
	}

	if cli.endDateArg != "" {
		endDateFilter, err := filter.NewEndDateFilter(cli.endDateArg)

		if err != nil {
			return err
		}

		cli.aggregator.AddFilter(endDateFilter)
	}

	if cli.keywordsArg != "" {
		keywords := strings.Split(cli.keywordsArg, ",")
		cli.aggregator.AddFilter(filter.NewKeywordFilter(keywords))
	}

	return nil
}

func (cli *CLI) printArticles(articles []article.Article) {

	params := print.FilterParams{
		SourceArg:    cli.sourceArg,
		KeywordsArg:  cli.keywordsArg,
		StartDateArg: cli.startDateArg,
		EndDateArg:   cli.endDateArg,
		OrderArg:     cli.sortOrderArg,
	}

	err := cli.printer.PrintArticles(articles, params)

	if err != nil {
		cli.printer.Error(err.Error())
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
	fmt.Println("  NewsAggregator -sort-order=asc")
}
