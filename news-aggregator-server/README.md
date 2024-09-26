# News Aggregator Helm Chart

This Helm chart deploys the **News Aggregator** application, a powerful tool for aggregating, filtering, and presenting
news articles from various sources.

## Features

- **Aggregation**: Collects news articles from sources in JSON, RSS, and HTML formats.
- **Filtering**: Supports filtering by sources, keywords, and date range.
- **Customization**: Allows for personalized news feeds with tailored filters.
- **Presentation**: Displays articles with metadata including title, description, publication date, source, and author.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+

## Installation

```bash
helm install news-aggregator
```

## Accessing the Application

Once the deployment is complete, you can access the News Aggregator web interface:

- **Local Access**:
    - `http://[::1]:port`

## API Endpoints

### Client API

- **Fetch Articles**: `/articles` (GET)
    - **Query Parameters**:
        - `source`: Filter by source.
        - `keywords`: Filter by keywords.
        - `date-start`: Filter by start date.
        - `date-end`: Filter by end date.
        - `sort`: Sort by a specific field.

### Admin API

- **Add Source**: `/source` (POST)
- **Delete Source**: `/source` (DELETE)
- **Update Source**: `/source` (PUT)

## Uninstallation

```bash
helm uninstall news-aggregator
```

This command removes all the Kubernetes components associated with the chart and deletes the release.

## License

This project is licensed under the Apache 2.0 License.
