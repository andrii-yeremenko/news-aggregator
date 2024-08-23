package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HotNewsSpec defines the desired state of HotNews
type HotNewsSpec struct {
	// Keywords are the search terms for the news.
	Keywords []string `json:"keywords"`

	// DateStart is the start date for the news search.
	DateStart *metav1.Time `json:"dateStart,omitempty"`

	// DateEnd is the end date for the news search.
	DateEnd *metav1.Time `json:"dateEnd,omitempty"`

	// Feeds are the news sources.
	// +optional
	Feeds []string `json:"feeds,omitempty"`

	// FeedGroups are the groups of news sources.
	// +optional
	FeedGroups []string `json:"feedGroups,omitempty"`

	// SummaryConfig defines how the status will show the summary of observed hot news.
	// +optional
	SummaryConfig SummaryConfig `json:"summaryConfig"`
}

// SummaryConfig defines how the status will show the summary of observed hot news.
type SummaryConfig struct {
	// +optional
	// +kubebuilder:validation:Required
	// +kubebuilder:default=10
	TitlesCount int `json:"titlesCount"`
}

// HotNewsStatus defines the observed state of HotNews
type HotNewsStatus struct {
	// This is the link to the news source.
	NewsLink string `json:"newsLink"`
	// This is the titles of fetched news articles.
	ArticlesTitles []string `json:"articlesTitles"`
	// This is the count of fetched news articles.
	ArticlesCount int `json:"articlesCount"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HotNews is the Schema for the hotnews API
type HotNews struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of HotNews
	Spec HotNewsSpec `json:"spec,omitempty"`

	// Status defines the observed state of HotNews
	Status HotNewsStatus `json:"status,omitempty"`
}

// HotNewsList contains a list of HotNews
type HotNewsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HotNews `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HotNews{}, &HotNewsList{})
}
