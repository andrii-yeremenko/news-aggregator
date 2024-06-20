package console_printer_test

import (
	"news-aggregator/aggregator/model/article"
	"news-aggregator/console_printer"
	"os"
	"testing"
	"time"
)

func TestPrintArticlesInTemplate(t *testing.T) {
	l := console_printer.New()

	articles := make([]article.Article, 0)
	testDate := time.Date(2024, time.June, 5, 12, 0, 0, 0,
		time.FixedZone("CEST", 2*60*60))
	art, _ := article.NewArticleBuilder().
		SetTitle("Test Title").
		SetDescription("Test Description").
		SetDate(article.CreationDate(testDate)).
		SetSource("Test Source").
		SetAuthor("Test Author").
		SetLink("https://testlink.com").
		Build()
	articles = append(articles, *art)

	params := console_printer.FilterParams{
		SourceArg:    "Test Source",
		KeywordsArg:  "Test",
		StartDateArg: "2024-06-05",
		EndDateArg:   "2024-06-05",
	}

	l.SetTemplatePath("testdata/article_simple_template.txt")
	err := l.PrintArticles(articles, params)
	if err != nil {
		t.Fatalf("Failed to print articles in template: %v", err)
	}

}

func TestPrintArticlesInTemplateError(t *testing.T) {

	l := console_printer.New()

	articles := make([]article.Article, 0)

	l.SetTemplatePath("testdata/template_with_error.txt")
	err := l.PrintArticles(articles, console_printer.FilterParams{})

	if err == nil {
		t.Errorf("PrintArticles() should return an error when template is invalid")
	}
}

func TestPrintArticlesInEmptyTemplate(t *testing.T) {

	l := console_printer.New()

	articles := make([]article.Article, 0)

	l.SetTemplatePath("testdata/empty_template.txt")
	err := l.PrintArticles(articles, console_printer.FilterParams{})

	if err == nil {
		t.Errorf("PrintArticles() should return an error when template is empty")
	}
}

func TestPrintArticlesInTemplate_WithIncorrectLocation(t *testing.T) {

	l := console_printer.New()

	articles := make([]article.Article, 0)

	l.SetTemplatePath("testdata/incorrect.txt")
	err := l.PrintArticles(articles, console_printer.FilterParams{})

	if err == nil {
		t.Errorf("PrintArticles() should return an error when template location is incorrect")
	}
}

func TestLog(t *testing.T) {
	l := console_printer.New()

	output, err := captureOutput(func() {
		l.Log("Test Log Message")
	})

	if err != nil {
		t.Fatalf("Failed to capture output: %v", err)
	}

	expectedOutput := "[Log] Test Log Message\n"
	if output != expectedOutput {
		t.Errorf("Log() output is incorrect, got: %s, want: %s", output, expectedOutput)
	}
}

func TestError(t *testing.T) {
	l := console_printer.New()

	output, err := captureOutput(func() {
		l.Error("Test Error Message")
	})

	if err != nil {
		t.Fatalf("Failed to capture output: %v", err)
	}

	expectedOutput := "[Error] Test Error Message\n"
	if output != expectedOutput {
		t.Errorf("Error() output is incorrect, got: %s, want: %s", output, expectedOutput)
	}
}

func TestWarn(t *testing.T) {
	l := console_printer.New()

	output, err := captureOutput(func() {
		l.Warn("Test Warning Message")
	})

	if err != nil {
		t.Fatalf("Failed to capture output: %v", err)
	}

	expectedOutput := "[Warning] Test Warning Message\n"
	if output != expectedOutput {
		t.Errorf("Warn() output is incorrect, got: %s, want: %s", output, expectedOutput)
	}
}

func captureOutput(f func()) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	err := w.Close()
	if err != nil {
		return "", err
	}

	capturedOutput := make(chan string)

	go func() {
		var buf [1024]byte
		n, _ := r.Read(buf[:])
		capturedOutput <- string(buf[:n])
	}()

	os.Stdout = old

	return <-capturedOutput, nil
}
