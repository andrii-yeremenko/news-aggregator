package console_printer

import (
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/fatih/color"
	"news-aggregator/aggregator/model/article"
	"os"
	"path"
	"regexp"
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
	for _, keyword := range keywords {
		pattern := `\b` + regexp.QuoteMeta(keyword) + `\b`
		re := regexp.MustCompile(pattern)
		highlighted := color.New(color.Underline).SprintFunc()(keyword)
		text = re.ReplaceAllString(text, highlighted)
	}
	return text
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
