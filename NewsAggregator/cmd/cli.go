package cmd

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/aggregator/filter"
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/logger"
	"NewsAggregator/storage"
	"flag"
	"strings"
)

// CLI is the command line interface for the news aggregator.
type CLI struct {
	sourceArg     string
	keywordsArg   string
	startDateArg  string
	endDateArg    string
	parserFactory *aggregator.Factory
	aggregator    *aggregator.Aggregator
	loader        *storage.ResourceLoader
}

// New creates a new CLI instance.
func New() *CLI {
	fact := aggregator.NewParserFactory()
	agr := aggregator.New(fact)
	loader := storage.NewLoader(agr)

	return &CLI{
		parserFactory: fact,
		aggregator:    agr,
		loader:        loader,
	}
}

// ParseFlags parses the command line flags.
func (cli *CLI) ParseFlags() {
	flag.StringVar(&cli.sourceArg, "sources", "", "Comma-separated list of news sources\n"+
		"Available sources: "+cli.loader.GetAvailableSources())
	flag.StringVar(&cli.keywordsArg, "keywords", "",
		"Comma-separated list of keywords to filter news articles")
	flag.StringVar(&cli.startDateArg, "date-start", "",
		"Start date for filtering news articles (format: yyyy-dd-mm)")
	flag.StringVar(&cli.endDateArg, "date-end", "",
		"End date for filtering news articles (format: yyyy-dd-mm)")
	flag.Parse()
}

// Run runs the CLI.
func (cli *CLI) Run() {

	if cli.loader.GetAvailableSources() == "" {
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
		cli.loader.LoadAllResources()
		cli.printArticles(cli.aggregator.GetAllArticles())
		return
	}

	if cli.sourceArg == "" {
		cli.loader.LoadAllResources()
	} else {
		sources := strings.Split(cli.sourceArg, ",")
		cli.loader.LoadSelectedResources(sources)
		cli.aggregator.AddFilter(filter.NewSourceFilter(sources))
	}

	if cli.startDateArg != "" {
		cli.aggregator.AddFilter(filter.NewStartDateFilter(cli.startDateArg))
	}

	if cli.endDateArg != "" {
		cli.aggregator.AddFilter(filter.NewEndDateFilter(cli.endDateArg))
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
}
