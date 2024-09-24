package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Condition is current state of the resource
type Condition struct {
	// Type of the condition (Added, Updated, Deleted, Failed)
	Type ConditionType `json:"type"`
	// Status of the condition (True or False)
	Status bool `json:"status"`
	// Reason is a brief readable explanation for the condition's last transition
	Reason string `json:"reason,omitempty"`
	// Message is a human-readable message indicating details about the last transition
	Message string `json:"message,omitempty"`
	// LastUpdateTime is the last time the condition was updated
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
}
