package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FeedSpec defines the desired state of Feed
type FeedSpec struct {
	// Name is the name of the news feed.
	// +kubebuilder:validation:Pattern=`^[a-zA-Z0-9-_]+$`
	Name string `json:"name"`

	// Link is the URL of the news feed.
	Link string `json:"link"`
}

// FeedStatus defines the observed state of Feed
type FeedStatus struct {
	// Conditions represent the latest available observations of an object's state.
	Conditions []Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Feed is the Schema for the feeds API
type Feed struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of Feed
	Spec FeedSpec `json:"spec,omitempty"`

	// Status defines the observed state of Feed
	Status FeedStatus `json:"status,omitempty"`
}

// FeedList contains a list of Feed
type FeedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Feed `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Feed{}, &FeedList{})
}
