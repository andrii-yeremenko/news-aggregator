package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// FeedSpec defines the desired state of Feed
type FeedSpec struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

// FeedStatus defines the observed state of Feed
type FeedStatus struct {
	LastUpdated  string `json:"lastUpdated"`
	CurrentState string `json:"state"`
	Message      string `json:"message"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Feed is the Schema for the feeds API
type Feed struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FeedSpec   `json:"spec,omitempty"`
	Status FeedStatus `json:"status,omitempty"`
}

func (in *Feed) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(Feed)
	in.DeepCopyInto(out)
	return out
}

func (in *FeedList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(FeedList)
	in.DeepCopyInto(out)
	return out
}

// +kubebuilder:object:root=true

// FeedList contains a list of Feed
type FeedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Feed `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Feed{}, &FeedList{})
}
