package console_printer

import (
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/fatih/color"
	"github.com/reiver/go-porterstemmer"
	"news-aggregator/aggregator/model/article"
	"os"
	"path"
	"strings"
	"text/template"
)

// Logger is a tool that records actions, measurements, print program output or other information.
type Logger struct {
	templatePath string
}

// FilterParams is a struct that holds the parameters for filtering articles.
type FilterParams struct {
	SourceArg    string
	KeywordsArg  string
	StartDateArg string
	EndDateArg   string
	OrderArg     string
}

// New creates a new Logger instance.
func New() *Logger {
	basePath, _ := os.Getwd()
	return &Logger{
		templatePath: path.Join(basePath, "console_printer/template/article_template.txt"),
	}
}

func highlightKeywords(text string, keywordsArg string) string {
	keywords := strings.Split(keywordsArg, ",")
	stemmedKeywords := make([]string, len(keywords))

	for i, keyword := range keywords {
		stemmedKeywords[i] = porterstemmer.StemString(keyword)
	}

	words := strings.Fields(text)
	stemmedWords := make([]string, len(words))
	for i, word := range words {
		stemmedWords[i] = porterstemmer.StemString(word)
	}

	for i, stemmedWord := range stemmedWords {
		for _, stemmedKeyword := range stemmedKeywords {
			if stemmedWord == stemmedKeyword {
				highlighted := color.New(color.Underline).SprintFunc()(words[i])
				words[i] = highlighted
				break
			}
		}
	}

	return strings.Join(words, " ")
}

func groupBySource(articles []article.Article) map[string][]article.Article {
	sourceGroups := make(map[string][]article.Article)
	for _, a := range articles {
		sourceGroups[string(a.Source())] = append(sourceGroups[string(a.Source())], a)
	}
	return sourceGroups
}

// PrintArticles prints a slice of article.Article to the console in predefined template.
func (l *Logger) PrintArticles(articles []article.Article, params FilterParams) error {

	data := struct {
		Articles []article.Article
		Params   FilterParams
	}{
		Articles: articles,
		Params:   params,
	}

	funcMap := template.FuncMap{
		"highlight":     highlightKeywords,
		"groupBySource": groupBySource,
	}

	tmpl, err := template.New("main").Funcs(funcMap).Funcs(sprig.FuncMap()).ParseFiles(l.templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	err = tmpl.ExecuteTemplate(os.Stdout, "main", data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
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

// SetTemplatePath sets the template for the logger.
func (l *Logger) SetTemplatePath(path string) {
	l.templatePath = path
}
