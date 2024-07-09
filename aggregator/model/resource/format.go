package resource

// Format the way in which information is arranged and stored on the resource.
type Format int

const (
	UNKNOWN = iota
	RSS
	HTML
	JSON
)
