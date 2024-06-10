package parser

import (
	"errors"
	"time"
)

// DateParser is a parser for date strings.
type DateParser struct {
	supportedFormats  []string
	defaultDateFormat string
}

// NewDateParser creates a new DateParser instance with predefined date formats.
func NewDateParser() *DateParser {
	return &DateParser{
		supportedFormats: []string{
			time.RFC1123Z,
			time.RFC3339,
			"3:04 p.m. ET January 2",
			"January 2, 2006",
			"Mon, 02 Jan 2006 15:04:05 GMT",
		},
		defaultDateFormat: "2006-02-01",
	}
}

// Parse parses the given string into a time.Time value using predefined supported date formats.
func (p *DateParser) Parse(dateStr string) (time.Time, error) {
	for _, layout := range p.supportedFormats {
		creationDate, err := time.Parse(layout, dateStr)
		if err == nil {
			if creationDate.Year() == 0 {
				creationDate = creationDate.AddDate(time.Now().Year(), 0, 0)
			}
			return creationDate, nil
		}
	}

	// Return a meaningful error message if unable to parse date.
	return time.Time{}, errors.New("unable to parse date")
}

// ParseDefaultDateFormat parses the given string into a time.Time value using the default project date format.
func (p *DateParser) ParseDefaultDateFormat(dateStr string) (time.Time, error) {
	return time.Parse(p.defaultDateFormat, dateStr)
}
