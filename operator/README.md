# Operator Documentation

## Overview

The News Aggregator Operator manages the lifecycle of news resources in a Kubernetes cluster, including both `Feed`
and `HotNews` resources. The operator automates the addition, update, and deletion of news sources and provides
functionality to create and manage hot news based on existing feeds and predefined criteria.

## Description

The operator is designed to handle two types of resources:

- **`Feed`**: Represents individual news sources.
- **`HotNews`**: Represents aggregated hot news based on specified criteria and feeds.

## Getting Started

### Prerequisites

Ensure you have the following tools installed:

- **Go**: Version 1.22.0+
- **Docker**: Version 17.03+
- **kubectl**: Version 1.11.3+
- **Kubernetes Cluster**: Version 1.11.3+ (Ensure you have access to a Kubernetes cluster)

For better development experience, consider use **`devbox`** to setup the development environment.

### Deployment

1. **Build and Push the Docker Image**

   Build and push your operator image to your container registry:

   ```sh
   make docker-build docker-push IMG=<some-registry>/operator:tag
   ```

   **Note:** Ensure the image is published in the registry specified and that you have access to pull the image.

2. **Install Custom Resource Definitions (CRDs)**

   Apply the Custom Resource Definitions (CRDs) to your cluster:

   ```sh
   make install
   ```

3. **Deploy the Operator**

   Deploy the operator using the image built in the previous step:

   ```sh
   make deploy IMG=<some-registry>/operator:tag
   ```

   **Note:** If you encounter RBAC (Role-Based Access Control) errors, you might need to grant cluster-admin privileges
   or ensure you are logged in as an admin.

4. **Create `Feed` Instances**

   Apply sample `Feed` resources to test the operator:

   ```sh
   kubectl apply -k config/samples/
   ```

5. **Create `HotNews` Instances**

   Apply sample `HotNews` resources to test the operator:

   ```sh
   kubectl apply -f config/samples/hotnews.yaml
   ```

   **Note:** Ensure the `feed-group-source` ConfigMap is set up correctly for `HotNews` to function properly.

### Uninstalling

1. **Delete Instances**

   Remove the `Feed` and `HotNews` resources from the cluster:

   ```sh
   kubectl delete -k config/samples/
   ```

2. **Remove CRDs**

   Uninstall the CRDs:

   ```sh
   make uninstall
   ```

3. **Undeploy the Operator**

   Remove the operator from the cluster:

   ```sh
   make undeploy
   ```

## CRD Definitions

### `Feed` Custom Resource

**Spec:**

- `name` (string): The name of the news source.
- `link` (string): The URL of the news feed.

**Status:**

- `conditions` (array of `Condition`): The status conditions of the `Feed`.

**Condition Types:**

- `Added`: The feed has been added successfully.
- `Updated`: The feed has been updated successfully.
- `Deleted`: The feed has been deleted successfully.

### `HotNews` Custom Resource

**Spec:**

- `keywords` (array of strings): List of keywords to filter hot news.
- `dateStart` (string, optional): Start date for filtering news.
- `dateEnd` (string, optional): End date for filtering news.
- `feeds` (array of strings, optional): List of feed names to include in the hot news.
- `feedGroups` (array of strings, optional): List of feed groups from the `feed-group-source` ConfigMap.
- `summaryConfig` (SummaryConfig): Configuration for displaying the summary of hot news.

**SummaryConfig:**

- `titlesCount` (integer): Number of article titles to display in the summary (default: 10).

**Status:**

- `newsLink` (string): Link to the news aggregator HTTPs server for the filtered news in JSON format.
- `articlesTitles` (array of strings): Titles of the articles, sorted by feed name.
- `articlesCount` (integer): Total number of articles matching the criteria.

### `feed-group-source` ConfigMap

The `feed-group-source` ConfigMap defines groups of feeds. Each key represents a feed group name, and each value is a
comma-separated list of feed names.

**Example:**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
   name: feed-group-source
data:
   sport: "usa-sports,nba"
   westMedias: "bbc,nbc,washingtontimes,nytimes"
```

## Supported Scenarios

1. **Feed Updates**

   Any update to `Feed` resources will trigger the `HotNews` reconciler to reprocess and update `HotNews` resources
   based on the updated feeds.

2. **ConfigMap Changes**

   Changes to the `feed-group-source` ConfigMap will trigger the `HotNews` reconciler to update `HotNews` resources that
   depend on the modified feed groups.

## License

Copyright 2024.

