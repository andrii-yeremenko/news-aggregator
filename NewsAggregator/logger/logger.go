package logger

import (
	"NewsAggregator/aggregator/model/article"
	"fmt"
)

// Logger is a tool that records actions, measurements, print program output or other information.
type Logger struct {
}

// New creates a new Logger instance.
func New() *Logger {
	return &Logger{}
}

// PrintArticle prints article.Article to the console.
func (l *Logger) PrintArticle(art article.Article) {
	fmt.Printf("----------------------------------------\n")
	fmt.Printf("Title: %s\n", art.Title())
	fmt.Printf("Description: %s\n", art.Description())
	fmt.Printf("Date: %s\n", art.Date().HumanReadableString())
	fmt.Printf("Source: %s\n", art.Source())
	fmt.Printf("Author: %s\n", art.Author())
}

// Log logs the given message.
func (l *Logger) Log(message string) {
	fmt.Printf("[Log] %s\n", message)
}

// Error logs the given error.
func (l *Logger) Error(error string) {
	fmt.Printf("[Error] %s\n", error)
}

// Warn logs the given warning.
func (l *Logger) Warn(warning string) {
	fmt.Printf("[Warning] %s\n", warning)
}
