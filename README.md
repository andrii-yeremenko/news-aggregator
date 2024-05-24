# News Aggregator

The News Aggregator is a versatile tool designed to collect, filter, and present news articles from a wide range of sources. Whether you're interested in technology, politics, or sports, this project has you covered. With its intuitive interface and powerful backend, staying informed has never been easier.

## Features

- **Aggregation**: Collects news articles from different websites in JSON, RSS, and HTML formats.
- **Filtering**: Allows users to filter articles by sources, keywords, and date range.
- **Customization**: Tailor your news feed to your interests with personalized filters.
- **Presentation**: Display filtered articles with relevant metadata such as title, description, publication date, source, and author.

## API Design

The News Aggregator API provides several functionalities to interact with and manipulate news data:

### Aggregation

The API aggregates news articles from various sources. Each source can provide data in different formats such as JSON, RSS, and HTML. The API handles parsing these formats using different parsers registered in the system.

### Filtering

Users can filter the aggregated articles based on several criteria:
- **Sources**: Filter articles by their source.
- **Keywords**: Filter articles by specific keywords.
- **Date Range**: Filter articles by a start and end date.

## Build

### For macOS and Linux

#### Build Script (.sh)

1. **Build the executable**: Open Terminal in the directory containing `main.go` file:
    ```bash
    go build -o news-aggregator main.go
    ```

2. **Make the application executable**: In your terminal run the following command:
    ```bash
    chmod +x news-aggregator.sh
    ```

3. **Usage**: To run the app, execute the following command in your terminal:
    ```bash
    ./news-aggregator.sh -keywords=technology, science
    ```

### For Windows

#### Command Prompt

1. **Open Command Prompt**: Open Command Prompt in the directory containing `main.go` file.

2. **Build the executable**: Run the following command to build the executable:
    ```bash
    go build -o news-aggregator.exe main.go
    ```

3. **Usage**: To run the executable with provided arguments, execute the following command in Command Prompt:
    ```bash
    news-aggregator.exe -source="nbc-news.com"
    news-aggregator.exe -keywords=technology,science
    news-aggregator.exe -date-start=2024-01-01 -date-end=2024-05-01
    ```



# Deployment

## Local Deployment
To deploy the News Aggregator locally, follow the installation instructions above and run the application using the go run command.

## Cloud Deployment
In the future, we aim to implement cloud solutions for deploying the News Aggregator project. Stay tuned for updates on cloud deployment options.

# Contribution
We welcome contributions from the community to improve and extend the functionality of the News Aggregator project. Whether it's fixing bugs, adding new features, or enhancing documentation, your contributions are highly appreciated.

# License

This project is licensed under the Apache 2.0.
