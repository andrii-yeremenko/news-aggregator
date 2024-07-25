# News Aggregator

The News Aggregator is a versatile tool designed to collect, filter, and present news articles from a wide range of
sources. Whether you're interested in technology, politics, or sports, this project has you covered. With its intuitive
interface and powerful backend, staying informed has never been easier.

## Features

- **Aggregation**: Collects news articles from different websites in JSON, RSS, and HTML formats.
- **Filtering**: Allows users to filter articles by sources, keywords, and date range.
- **Customization**: Tailor your news feed to your interests with personalized filters.
- **Presentation**: Display filtered articles with relevant metadata such as title, description, publication date,
  source, and author.

## API Design

The News Aggregator API provides several functionalities to interact with and manipulate news data:

### Aggregation

The API aggregates news articles from various sources. Each source can provide data in different formats such as JSON,
RSS, and HTML. The API handles parsing these formats using different parsers registered in the system.

### Filtering

Users can filter the aggregated articles based on several criteria:

- **Sources**: Filter articles by their source.
- **Keywords**: Filter articles by specific keywords.
- **Date Range**: Filter articles by a start and end date.

## Get Started:

## Command Line Interface
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

## Web Interface

### Run Locally

The image is available on Docker Hub, and you can run it locally using Docker.
Image supports both `amd64` and `arm64` architectures.
Follow the steps below to run the News Aggregator web interface:

1. **Pull Docker Image**: Pull the Docker image from the Docker Hub repository:
    ```bash
    docker pull ayeremenko/news-aggregator:(version you want)
    ```
2. **Run Docker Container**: Run the Docker container with the following command:
    ```bash
    docker run ayeremenko/news-aggregator:(version you pulled)
    ```
3. **Access the Web Interface**: Open your browser and navigate to `http://[::1]:8443`.

4. **To check if web server running**: Open your browser and navigate to `http://[::1]:8443/status`.

Also docker image provides the following environment variables to configure the application:

- `PORT` - port to run the web server on (default is 8443)
  To set the port to 8080, run the following command:
    ```bash
    docker run -e PORT=8080 ayeremenko/news-aggregator
    ```
- `TIMEOUT` - timeout for the web server (default is 12h)
  To set the timeout to 1h, run the following command:
     ```bash
     docker run -e TIMEOUT=1h ayeremenko/news-aggregator
     ```

## Web Server API Documentation

### Client API

1. **Fetch Articles**: Retrieve articles from the server.
    - **URL**: `/articles`
    - **Method**: `GET`
    - **Query Parameters**:
        - `source`: Filter articles by source.
        - `keywords`: Filter articles by keywords.
        - `date-start`: Filter articles by start date.
        - `date-end`: Filter articles by end date.
        - `sort`: Sort articles by a specific field.
    - **Response**: Returns a JSON formatted text of articles that match the specified criteria.

### Admin API

1. **Add Source**: Add a new source to the system.
    - **URL**: `/source`
    - **Method**: `POST`
    - **Request Body**: JSON object containing the source information.
        - `name`: The name of the source.
        - `url`: The URL of the source.
        - `format`: The format of the source data (JSON, RSS, HTML).
    - **Response**: `201 Created` if the source was added successfully.
      Example:
    ```json
    {
      "name": "abc-news",
      "url": "https://feeds.abcnews.com/abcnews/internationalheadlines",
      "format": "RSS"
    }
    ```
2. **Delete Source**: Remove a source from the system.
    - **URL**: `/source`
    - **Method**: `DELETE`
    - **Response**: `200 Ok` if the source was deleted successfully.

   Example:
    ```json
    {
      "name": "abc-news"
    }  
    ```   

3. **Update Source**: Update an existing source in the system.
    - **URL**: `/source`
    - **Method**: `PUT`
    - **Request Body**: JSON object containing the source information.
        - `name`: The name of the source.
        - `url`: The URL of the source.
        - `format`: The format of the source data (JSON, RSS, HTML).
    - **Response**: `200 Ok` if the source was updated successfully.

   Example:
    ```json
    {
      "name": "abc-news",
      "url": "https://feeds.abcnews.com/abcnews/internationalheadlines",
      "format": "RSS"
    }
    ```

4. **Update Feed**: Fetch articles from the internet and update the database.
    - **URL**: `/update`
    - **Method**: `GET`
    - **Response**: `200 Ok` if the feed was updated successfully.

# Contribution

We welcome contributions from the community to improve and extend the functionality of the News Aggregator project.
Whether it's fixing bugs, adding new features, or enhancing documentation, your contributions are highly appreciated.

# License

This project is licensed under the Apache 2.0.