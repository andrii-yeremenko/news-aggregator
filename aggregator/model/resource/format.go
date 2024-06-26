package resource

import "fmt"

// Format the way in which information is arranged and stored on the resource.
type Format int

const (
	UNKNOWN = iota
	RSS
	HTML
	JSON
)

// FormatToString converts a Format to a string.
func FormatToString(f Format) string {
	switch f {
	case RSS:
		return "RSS"
	case HTML:
		return "HTML"
	case JSON:
		return "JSON"
	default:
		return "UNKNOWN"
	}
}

// ParseFormat converts a string to a Format.
func ParseFormat(formatStr string) (Format, error) {
	switch formatStr {
	case "RSS":
		return RSS, nil
	case "HTML":
		return HTML, nil
	case "JSON":
		return JSON, nil
	default:
		return UNKNOWN, fmt.Errorf("unknown format: %s", formatStr)
	}
}
