/*
Package controller provides the API that manages the lifecycle of the Feed and HotNews resources in the news-aggregator.

Feed Controller (feed_controller.go):

	– Creates, updates, or deletes the corresponding news feed in the news-aggregator when a Feed resource is created,

updated, or deleted.

	– Updates the status of the Feed resource with conditions such as whether the feed was successfully added to the

aggregator.

HotNews Controller (hotnews_controller.go):

	– Fetches and filters news based on the criteria provided in the HotNews resource (keywords, feeds, date range).
	– Automatically triggers the HotNews reconciler when:
		A feed defined in the HotNews resource is updated.
		The feed-group-source ConfigMap is updated.
	– Updates the status with a count of articles and a link to view the news in JSON format.
	– Stores a limited number of article titles (configurable via spec.summaryConfig.titlesCount) in the status.
*/
package controller
