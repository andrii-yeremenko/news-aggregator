package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HotNewsSpec defines the desired state of HotNews
type HotNewsSpec struct {
	// +kubebuilder:validation:Required
	Keywords []string `json:"keywords"`

	// +kubebuilder:validation:Required
	DateStart *metav1.Time `json:"dateStart,omitempty"`

	// +kubebuilder:validation:Required
	DateEnd *metav1.Time `json:"dateEnd,omitempty"`

	// +kubebuilder:validation:Required
	Feeds []string `json:"feeds,omitempty"`

	// +optional
	// +kubebuilder:validation:Required
	FeedGroups []string `json:"feedGroups,omitempty"`

	// +optional
	// +kubebuilder:validation:Required
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
	// +kubebuilder:validation:Required
	NewsLink string `json:"newsLink"`
	// +kubebuilder:validation:Required
	ArticlesTitles []string `json:"articlesTitles"`
	// +kubebuilder:validation:Required
	ArticlesCount int `json:"articlesCount"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HotNews is the Schema for the hotnews API
type HotNews struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec HotNewsSpec `json:"spec,omitempty"`

	// +kubebuilder:validation:Optional
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
