## Project Name: NewsAggregator API

## Engineer name: `Yeremenko Andrii`

## Summary

The `News Aggregator` API is a tool designed to aggregate news articles from various sources. It offers functionality to
filter articles by sources, keywords, and date ranges, providing users with a curated selection of relevant news
content. The project utilizes parsers for different formats (such as `JSON`, `RSS`, and `HTML`) to extract article data,
which
is then stored and filtered based on user-defined criteria. This API aims to simplify the process of accessing and
organizing news articles from diverse sources, enhancing the efficiency of news consumption and research.

## Motivation

The `News Aggregator` API addresses the increasing need for efficient news consumption and analysis in today's
fast-paced
digital landscape. By consolidating articles from multiple sources and offering flexible filtering options, it supports
various use cases such as news monitoring, research, and content curation. Users can leverage the API to stay updated on
specific topics or events, track trends, and gather insights from diverse perspectives. The expected outcome is to
empower users with a convenient and efficient tool for accessing curated news content, thereby enhancing their ability
to stay informed and up-to-date with current events.

# APIs design

## Parser:

Parsers are a crucial component of the News Aggregator API, responsible for extracting structured data from raw news
sources in various formats.
Each parser implements the Parser interface and is designed to handle a specific data format (
e.g., `JSON`, `RSS`, `HTML`).

### Available Parsers

1. JSON Parser (json_parser.go)
    * Description: Parses news articles from JSON data.
    * Args:
        * content (resource.Resource): The JSON resource to parse.
    * Returns:
        * []article.Article: A list of parsed articles.
        * error: Error object in case of failure.
    * Errors:
        * JSONParseError: Error occurred while parsing JSON data.
    * Docs:
        * Description: Parses JSON data to extract news articles.
    * Usage:
      ```
      articles, err := jsonParser.Parse(resource)
      ```
2. RSS Parser (rss_parser.go)
    * Description: Parses news articles from RSS data.
    * Args:
        * content (resource.Resource): The RSS resource to parse.
    * Returns:
        * []article.Article: A list of parsed articles.
        * error: Error object in case of failure.
    * Errors:
        * RSSParseError: Error occurred while parsing RSS data.
    * Docs:
        * Description: Parses RSS data to extract news articles.
    * Usage:
      ```
      articles, err := rssParser.Parse(resource)
      ```
3. USA Today HTML Parser (usa_today_html_parser.go)
    * Description: Parses news articles from USA Today HTML data.
    * Args:
        * content (resource.Resource): The USA Today HTML resource to parse.
    * Returns:
        * []article.Article: A list of parsed articles.
        * error: Error object in case of failure.
    * Errors:
        * HTMLParseError: Error occurred while parsing HTML data.
    * Docs:
        * Description: Parses USA Today HTML data to extract news articles.
    * Usage:
    ```
    articles, err := usaTodayHTMLParser.Parse(resource)
    ```

## Parser Factory:

The parser factory API manages the registration and instantiation of parser objects based on the data format and source.
It allows dynamic selection of parsers according to the source's format (e.g., JSON, RSS, HTML) and facilitates the
extraction of structured data from raw sources. Input arguments typically include the data format and source URL. The
output is an instantiated parser object capable of parsing articles from the specified source in the designated format.

### Methods

1. RegisterParser: Registers a parser for a specific format and source.
    * Args:
        * format (string): Format of the data (e.g., "json", "rss").
        * source (string): Source URL or identifier.
        * parser (Parser): Parser object capable of parsing the specified format.
    * Returns:
        * error: Error object in case of failure.
    * Errors:
        * DuplicateParserError: A parser is already registered for the specified format and source.
    * Docs:
        * Description: Registers a new parser for the given format and source.

   #### Usage:

    ```
    pf := parser.NewParserFactory()
    pf.RegisterParser("json", "json-news.com", &parser.JSONParser{})
    pf.RegisterParser("rss", "xml-news.com", &parser.RSSParser{})
    ```
2. GetParser: Retrieves a parser for a specific format and source.
    * Args:
        * format (string): Format of the data (e.g., "json", "rss").
        * source (string): Source URL or identifier.
    * Returns:
        * Parser: Parser object capable of parsing the specified format.
        * error: Error object in case of failure.
    * Errors:
        * ParserNotFoundError: No parser is registered for the specified format and source.
    * Docs:
        * Description: Retrieves the parser for the given format and source.

   #### Usage:

     ```
     pf := parser.NewParserFactory()
     jsonParser, err := pf.GetParser("json", "json-news.com")
     rssParser, err := pf.GetParser("rss", "xml-news.com")
     ```

## Storage:

The storage API is responsible for managing the storage and retrieval of news article data. It abstracts away the
details of data storage implementation, allowing flexibility in choosing storage mediums (e.g., file system, database,
in-memory storage). Input arguments for this API typically include the source name, format/type of data (e.g., `JSON`,
`RSS`, `HTML`), and file path or resource location. The output is a Resource object.
Also, it has an error handling mechanism, which allows to handle errors in a way that is convenient for the user.

### Methods

1. ReadFile: Reads and loads articles from a file into the storage system.

    * Args:
        * sourceName (string): Name of the news source.
        * format (string): Format of the data (e.g., "json", "rss", "html").
        * filePath (string): Path to the file containing the news data.

    * Returns:
        * Resource object containing the loaded data.
        * error: Error object in case of failure.

    * Errors:
        * FileNotFoundError: The specified file does not exist.
        * FormatError: The data format is not supported.

    * Docs:
        * Description: Reads a file from the given path and loads the data into the repository.

   #### Usage:

    ```
    newsStorage := storage.New()
    nbc, err := newsStorage.ReadFile("nbc-news.com", "json", "storage/resources/nbc-news.json")
    ```

## Aggregator:

The aggregator API orchestrates the aggregation and filtering of news articles from multiple sources. It leverages the
registered parsers to extract article data from raw sources and aggregates them into a unified collection. Input
arguments include the loaded news resources and optional filtering parameters such as sources, keywords, and date
ranges. The output is a curated selection of news articles that meet the specified criteria.
Also, it has a error handling mechanism, which allows to handle errors in a way that is convenient for the user.

### Methods

1. LoadResource: Loads a news resource into the aggregator.

    * Args:
        * resource (Resource): Resource object containing the news data.
    * Returns:
        * error: Error object in case of failure.
    * Errors:
        * ResourceLoadError: Error occurred while loading the resource.
    * Docs:
        * Description: Loads a given resource into the aggregator for processing.

   #### Usage:
    ```
    newsAggregator := aggregator.NewAggregator(parserSelector)
    loadErr = newsAggregator.LoadResource(nbc)
    ```    

2. GetArticles: Retrieves aggregated articles.

    * Returns:
        * []Article: A list of aggregated articles matching the filters.
        * error: Error object in case of failure.
    * Errors:
        * FetchError: Error occurred while fetching the articles.
    * Docs:
        * Description: Fetches and returns a list of aggregated articles based on specified filters.

   #### Usage:
    ```
    newsAggregator := aggregator.NewAggregator(parserSelector)
    articles, fetchError := newsAggregator.GetArticles()
    ```

## Filter:

The filter API provides mechanisms to filter news articles based on specific criteria, such as keywords and sources. The
package includes a Filter interface and concrete implementations for keyword-based, source-based and date-based
filtering.

### Best usage practices

To use a filter, create an instance of a specific filter type (e.g., KeywordFilter or SourceFilter) and call the Apply
method with a slice of article.Article objects.

### Methods

1. Apply: Applies the filter to a list of articles.

    * Args:
        * articles ([]Article): List of articles to be filtered.
    * Returns:
        * []Article: A list of articles that match the filter criteria.
    * Errors:
        * FilterError: Error occurred while applying the filter.
    * Docs:
        * Description: Filters the given list of articles based on the implemented criteria.

   #### Example:
     ```
     keywords := []string{"technology", "science"}
     keywordFilter := filter.NewKeywordFilter(keywords)
     filteredByKeywords := keywordFilter.Apply(articles)
     fmt.Println("Articles filtered by keywords:", filteredByKeywords)
  
     sources := []string{"TechCrunch", "Wired"}
     sourceFilter := filter.NewSourceFilter(sources)
     filteredBySources := sourceFilter.Apply(articles)
     fmt.Println("Articles filtered by sources:", filteredBySources)
     ```
   This example demonstrates how to create and use both KeywordFilter and SourceFilter to filter articles based on
   specified criteria.

# Resume

These APIs interact seamlessly to fetch, parse, and filter news articles according to user-defined criteria, providing a
comprehensive solution for news aggregation and consumption. The repository manages the storage and retrieval of article
data, the parser factory dynamically selects parsers based on source format, and the aggregator coordinates the
aggregation and filtering processes to deliver curated news content to users.

# Unresolved questions

1. Are there any specific requirements for error handling in the API?
2. How frequently should the news articles be updated and refreshed in the repository?
3. Do we need unstage the news articles from aggregator after a certain period of time?
4. How should the API provide an error response to the user in case of a failure?