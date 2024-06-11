package cmd

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/aggregator/filter"
	"NewsAggregator/aggregator/model/article"
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
	agr, err := aggregator.New(fact)

	if err != nil {
		logger.New().Error(err.Error())
	}

	basePath, err := os.Getwd()

	if err != nil {
		logger.New().Error(err.Error())
	}

	stg := storage.New(basePath + "/storage")

	return &CLI{
		parserFactory: fact,
		aggregator:    agr,
		storage:       stg,
	}
}

// ParseFlags parses the command line flags.
func (cli *CLI) ParseFlags() {
	flag.StringVar(&cli.sourceArg, "sources", "", "Comma-separated list of news sources\n"+
		"Available sources: "+cli.storage.GetAvailableSources())
	flag.StringVar(&cli.keywordsArg, "keywords", "",
		"Comma-separated list of keywords to filter news articles")
	flag.StringVar(&cli.startDateArg, "date-start", "",
		"Start date for filtering news articles (format: yyyy-dd-mm)")
	flag.StringVar(&cli.endDateArg, "date-end", "",
		"End date for filtering news articles (format: yyyy-dd-mm)")
	flag.Usage = func() {
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
	flag.Parse()
}

// Run runs the CLI.
func (cli *CLI) Run() {

	if cli.storage.GetAvailableSources() == "" {
		logger.New().Warn("No sources available")
		return
	}

	flagCount := 0
	flag.Visit(func(f *flag.Flag) {
		flagCount++
	})

	if flagCount > 4 {
		flag.Usage()
		return
	}

	if flagCount == 0 {
		resources, err := cli.storage.GetAllResources()

		if err != nil {
			logger.New().Error(err.Error())
			return
		}

		for _, res := range resources {
			cli.aggregator.LoadResource(res)
		}

		cli.printArticles(cli.aggregator.GetAllArticles())
		return
	}

	if cli.sourceArg == "" {
		resources, err := cli.storage.GetAllResources()
		if err != nil {
			logger.New().Error(err.Error())
			return
		}
		for _, res := range resources {
			cli.aggregator.LoadResource(res)
		}
	} else {
		sources := strings.Split(cli.sourceArg, ",")
		resources, err := cli.storage.GetSelectedResources(sources)
		if err != nil {
			logger.New().Error(err.Error())
			return
		}
		for _, res := range resources {
			cli.aggregator.LoadResource(res)
		}
		cli.aggregator.AddFilter(filter.NewSourceFilter(sources))
	}

	if cli.startDateArg != "" {
		f, err := filter.NewStartDateFilter(cli.startDateArg)
		if err != nil {
			logger.New().Error(err.Error())
		} else {
			cli.aggregator.AddFilter(f)
		}
	}

	if cli.endDateArg != "" {
		f, err := filter.NewEndDateFilter(cli.endDateArg)

		if err != nil {
			logger.New().Error(err.Error())
		} else {
			cli.aggregator.AddFilter(f)
		}
	}

	if cli.keywordsArg != "" {
		keywords := strings.Split(cli.keywordsArg, ",")
		cli.aggregator.AddFilter(filter.NewKeywordFilter(keywords))
	}

	filteredArticles := cli.aggregator.GetFilteredArticles()
	cli.printArticles(filteredArticles)
}

func (cli *CLI) printArticles(articles []article.Article) {

	log := logger.New()
	log.Log("Printing articles")

	for _, art := range articles {
		log.PrintArticle(art)
	}

	log.Log(fmt.Sprint(len(articles)) + " news articles were shown")
}
