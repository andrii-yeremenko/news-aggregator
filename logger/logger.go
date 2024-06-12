package logger

import (
	"NewsAggregator/aggregator/model/article"
	"fmt"
	"github.com/fatih/color"
	"os"
	"regexp"
	"strings"
	"text/template"
)

// Logger is a tool that records actions, measurements, print program output or other information.
type Logger struct {
}

// FilterParams is a struct that holds the parameters for filtering articles.
type FilterParams struct {
	SourceArg    string
	KeywordsArg  string
	StartDateArg string
	EndDateArg   string
}

// New creates a new Logger instance.
func New() *Logger {
	return &Logger{}
}

func highlightKeywords(text string, keywordsArg string) string {
	keywords := strings.Split(keywordsArg, ",")
	for _, keyword := range keywords {
		pattern := `\b` + regexp.QuoteMeta(keyword) + `\b`
		re := regexp.MustCompile(pattern)
		highlighted := color.New(color.Underline).SprintFunc()(keyword)
		text = re.ReplaceAllString(text, highlighted)
	}
	return text
}

// PrintArticlesInTemplate prints a slice of article.Article to the console in predefined template.
func (l *Logger) PrintArticlesInTemplate(articles []article.Article, params FilterParams, templatePath string) error {

	data := struct {
		Articles []article.Article
		Params   FilterParams
	}{
		Articles: articles,
		Params:   params,
	}

	funcMap := template.FuncMap{
		"highlight": highlightKeywords,
	}

	tmpl, err := template.New("main").Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(os.Stdout, "main", data)
	if err != nil {
		return err
	}

	return nil
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
