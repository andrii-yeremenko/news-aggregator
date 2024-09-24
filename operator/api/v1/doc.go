/*
Package v1 contains API schema definitions for the Feed and HotNews CRDs in the
teamdev.com/v1 group.

Overview:

This package defines two Kubernetes Custom Resources (CRs) – Feed and HotNews.
These resources are used to interact with a news aggregator and manage
news feeds and hot news topics within a Kubernetes cluster.

The operator, built using Kubebuilder, manages the lifecycle of these CRs,
allowing CRUD operations (Create, Read, Update, Delete) on the news feeds and
hot news topics via Kubernetes native API. All operations are handled by
controller-manager, which watches for changes in the CRs and triggers the
corresponding actions in the news-aggregator.

Custom Resources:

1. Feed

The Feed custom resource represents a news source in the aggregator. When a
Feed resource is created, it triggers the creation of a new feed in the
news-aggregator. If a Feed is updated or deleted, the corresponding source in
the aggregator will be updated or removed, respectively.

Feed Resource Schema:

	– apiVersion: teamdev.com/v1
	– kind: Feed
	– metadata:
	  – name: The name of the Feed resource.
	  – namespace: The namespace where the resource is created.
	– spec:
	  – name: The name of the news source.
	  – link: The URL or link of the news source.
	– status:
	  – conditions: A list of conditions representing the status of the feed in
	  the news-aggregator.
	  – type: The condition type, for example, "Added".
	  – status: Whether the feed was successfully added to the aggregator (True/False).
	  – reason: If the status is False, a reason for failure.
	  – message: If the status is False, a message providing additional context.
	  – lastUpdateTime: The timestamp when the status last changed.

2. HotNews

The HotNews custom resource is used to aggregate and filter hot news based on
specific criteria such as keywords, feeds, and feed groups. It pulls news from
the available feeds and feedGroups, and filters them
based on keywords, start and end dates, and other configurations.

Reconciliation Triggers:

	– Updating any feed in the HotNews spec triggers the HotNews reconciler.
	– Changes in the feed-group-source ConfigMap also trigger the HotNews reconciler.

HotNews Resource Schema:

	– apiVersion: teamdev.com/v1
	– kind: HotNews
	– metadata:
	  – name: The name of the HotNews resource.
	  – namespace: The namespace where the resource is created.
	– spec:
	  – keywords: A list of search terms for the news.
	  – dateStart: The start date for the news search.
	  – dateEnd: The end date for the news search.
	  – feeds: A list of news sources to pull news from.
	  – feedGroups: A list of groups of news sources.
	  – summaryConfig: Configuration for the summary of observed hot news.
	– status:
	  – newsLink: The link to the news source.
	  – articlesTitles: Titles of fetched news articles.
	  – articlesCount: The count of fetched news articles.
	  – conditions: A list of ConditionType representing the status of the hot news
	  topic.
*/
package v1
