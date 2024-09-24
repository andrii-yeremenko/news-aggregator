package v1

// ConditionType defines the type of the condition
type ConditionType string

const (
	// ConditionAdded means the feed has been added successfully
	ConditionAdded ConditionType = "Added"
	// ConditionUpdated means the feed has been updated successfully
	ConditionUpdated ConditionType = "Updated"
	// ConditionDeleted means the feed has been deleted successfully
	ConditionDeleted ConditionType = "Deleted"
	// ConditionFailed means the feed has failed to be added, updated or deleted
	ConditionFailed ConditionType = "Failed"
)
