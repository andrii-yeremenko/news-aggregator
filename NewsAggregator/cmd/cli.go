package cmd

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/aggregator/filter"
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/storage"
	"flag"
	"fmt"
	"strings"
)

type CLI struct {
	sourceArg     string
	keywordsArg   string
	startDateArg  string
	endDateArg    string
	parserFactory *aggregator.Factory
	aggregator    *aggregator.Aggregator
	loader        *storage.ResourceLoader
}

func NewCLI() *CLI {
	fact := aggregator.NewParserFactory()
	agr := aggregator.New(fact)
	loader := storage.NewLoader(agr)

	return &CLI{
		parserFactory: fact,
		aggregator:    agr,
		loader:        loader,
	}
}

func (cli *CLI) ParseFlags() {
	flag.StringVar(&cli.sourceArg, "sources", "", "Comma-separated list of news sources\n"+
		"Available sources: "+cli.loader.GetAvailableSources())
	flag.StringVar(&cli.keywordsArg, "keywords", "",
		"Comma-separated list of keywords to filter news articles")
	flag.StringVar(&cli.startDateArg, "date-start", "",
		"Start date for filtering news articles (format: yyyy-mm-dd)")
	flag.StringVar(&cli.endDateArg, "date-end", "",
		"End date for filtering news articles (format: yyyy-mm-dd)")
	flag.Parse()
}

func (cli *CLI) Run() {

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
	for _, art := range articles {
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("Title: %s\n", art.Title())
		fmt.Printf("Description: %s\n", art.Description())
		fmt.Printf("Date: %s\n", art.Date().HumanReadableString())
		fmt.Printf("Source: %s\n", art.Source())
		fmt.Printf("Author: %s\n", art.Author())
	}
}
