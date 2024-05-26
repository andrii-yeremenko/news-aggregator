- Project Name: NewsAggregator API
- Engineer name: Yeremenko Andrii

# Summary
The News Aggregator API is a tool designed to aggregate news articles from various sources. It offers functionality to filter articles by sources, keywords, and date ranges, providing users with a curated selection of relevant news content. The project utilizes parsers for different formats (such as JSON, RSS, and HTML) to extract article data, which is then stored and filtered based on user-defined criteria. This API aims to simplify the process of accessing and organizing news articles from diverse sources, enhancing the efficiency of news consumption and research.

# Motivation

The News Aggregator API addresses the increasing need for efficient news consumption and analysis in today's fast-paced digital landscape. By consolidating articles from multiple sources and offering flexible filtering options, it supports various use cases such as news monitoring, research, and content curation. Users can leverage the API to stay updated on specific topics or events, track trends, and gather insights from diverse perspectives. The expected outcome is to empower users with a convenient and efficient tool for accessing curated news content, thereby enhancing their ability to stay informed and up-to-date with current events.
# APIs design

### Repository:
The repository API is responsible for managing the storage and retrieval of news article data. It abstracts away the details of data storage implementation, allowing flexibility in choosing storage mediums (e.g., file system, database, in-memory storage). Input arguments for this API typically include the source name, format/type of data (e.g., JSON, RSS, HTML), and file path or resource location. The output is a Resource object.
Also it has a error handling mechanism, which allows to handle errors in a way that is convenient for the user.
#### Example:
```
newsRepository := repository.NewInMemoryRepository()
nbc, repositoryError := newsRepository.ReadFile("nbc-news.com", "json", "repository/news-resources/nbc-news.json")
```
### Parser Factory:
The parser factory API manages the registration and instantiation of parser objects based on the data format and source. It allows dynamic selection of parsers according to the source's format (e.g., JSON, RSS, HTML) and facilitates the extraction of structured data from raw sources. Input arguments typically include the data format and source URL. The output is an instantiated parser object capable of parsing articles from the specified source in the designated format.

#### Example:
```
parserSelector := parser.NewParserFactory()
parserSelector.RegisterParser("json", "json-news.com", &parser.JSONParser{})
parserSelector.RegisterParser("rss", "xml-news.com", &parser.RSSParser{})
```

### Aggregator:
The aggregator API orchestrates the aggregation and filtering of news articles from multiple sources. It leverages the registered parsers to extract article data from raw sources and aggregates them into a unified collection. Input arguments include the loaded news resources and optional filtering parameters such as sources, keywords, and date ranges. The output is a curated selection of news articles that meet the specified criteria.
Also, it has a error handling mechanism, which allows to handle errors in a way that is convenient for the user.

#### Example:
```
newsAggregator := aggregator.NewAggregator(parserSelector)
loadErr = newsAggregator.LoadResource(nbc)
articles, fetchError := newsAggregator.GetArticles()
```

### Filter:
The filter API provides functionality to filter news articles based on user-defined criteria such as sources, keywords, and date ranges. It accepts a collection of articles and filtering parameters as input and returns a subset of articles that match the specified criteria. This API enables users to customize their news feed according to their preferences and interests, enhancing the relevance and quality of the curated content.
Also, it has a builder pattern, which allows to create complex filters in a convenient way. 

#### Example:
```
filterBuilder := filter.NewFilterBuilder()
filterBuilder.SetSources([]string{"nbc-news.com", "cnn.com"})
filterBuilder.SetKeywords([]string{"politics", "economy"})
filterBuilder.SetDateRange("2022-01-01", "2022-01-31")
filter := filterBuilder.Build()

filteredArticles := filter.Filter(articles) # filteredArticles is a subset of given news articles that meet  
                                            # the specified criteria.
```

# Resume
These APIs interact seamlessly to fetch, parse, and filter news articles according to user-defined criteria, providing a comprehensive solution for news aggregation and consumption. The repository manages the storage and retrieval of article data, the parser factory dynamically selects parsers based on source format, and the aggregator coordinates the aggregation and filtering processes to deliver curated news content to users.

# Unresolved questions
1. Are there any specific requirements for error handling in the API?
2. How frequently should the news articles be updated and refreshed in the repository?
3. Do we need unstage the news articles from aggregator after a certain period of time?