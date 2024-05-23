package parser

import (
	"errors"
	"time"
)

// DateParser is a parser for date strings.
type DateParser struct {
	dateFormats []string
}

// NewDateParser creates a new DateParser instance with predefined date formats.
func NewDateParser() *DateParser {
	return &DateParser{
		dateFormats: []string{
			time.RFC1123Z,
			time.RFC3339,
			"3:04 p.m. ET January 2",
			"January 2, 2006",
			"Mon, 02 Jan 2006 15:04:05 GMT",
		},
	}
}

// Parse parses the given string into a time.Time value using predefined supported date formats.
func (dp *DateParser) Parse(dateStr string) (time.Time, error) {
	for _, layout := range dp.dateFormats {
		creationDate, err := time.Parse(layout, dateStr)
		if err == nil {
			if creationDate.Year() == 0 {
				// Set the year to the current year if not provided in the date string.
				creationDate = creationDate.AddDate(time.Now().Year(), 0, 0)
			}
			return creationDate, nil
		}
	}

	// Return a meaningful error message if unable to parse date.
	return time.Time{}, errors.New("unable to parse date")
}
