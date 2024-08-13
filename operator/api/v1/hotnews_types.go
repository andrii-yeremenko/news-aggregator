package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HotNewsSpec defines the desired state of HotNews
type HotNewsSpec struct {
	Keywords []string `json:"keywords"`

	DateStart *metav1.Time `json:"dateStart,omitempty"`

	DateEnd *metav1.Time `json:"dateEnd,omitempty"`

	Feeds []string `json:"feeds,omitempty"`

	FeedGroups []string `json:"feedGroups,omitempty"`

	SummaryConfig SummaryConfig `json:"summaryConfig"`
}

// SummaryConfig defines how the status will show the summary of observed hot news.
type SummaryConfig struct {
	TitlesCount int `json:"titlesCount"`
}

// HotNewsStatus defines the observed state of HotNews
type HotNewsStatus struct {
	// +optional
	NewsLink string `json:"newsLink"`
	// +optional
	ArticlesTitles []string `json:"articlesTitles"`
	// +optional
	ArticlesCount int `json:"articlesCount"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HotNews is the Schema for the hotnews API
type HotNews struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HotNewsSpec   `json:"spec,omitempty"`
	Status HotNewsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HotNewsList contains a list of HotNews
type HotNewsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HotNews `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HotNews{}, &HotNewsList{})
}
