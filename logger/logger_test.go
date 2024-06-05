package logger_test

import (
	"NewsAggregator/aggregator/model/article"
	"NewsAggregator/logger"
	"os"
	"testing"
	"time"
)

func TestPrintArticle(t *testing.T) {
	l := logger.New()

	articleBuilder := article.NewArticleBuilder()
	testDate := time.Date(2024, time.June, 5, 12, 0, 0, 0,
		time.FixedZone("CEST", 2*60*60))
	art, err := articleBuilder.
		SetTitle("Test Title").
		SetDescription("Test Description").
		SetDate(article.CreationDate(testDate)).
		SetSource("Test Source").
		SetAuthor("Test Author").
		SetLink("https://testlink.com").
		Build()

	if err != nil {
		t.Fatalf("Failed to create art: %v", err)
	}

	expectedOutput := "----------------------------------------\n" +
		"Title: Test Title\n" +
		"Description: Test Description\n" +
		"Date: 05 Jun 24 12:00 CEST\n" +
		"Source: Test Source\n" +
		"Author: Test Author\n" +
		"Link: https://testlink.com\n"

	output := captureOutput(func() {
		l.PrintArticle(*art)
	})

	if output != expectedOutput {
		t.Errorf("PrintArticle() output is incorrect, got: %s, want: %s", output, expectedOutput)
	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	err := w.Close()
	if err != nil {
		return ""
	}

	capturedOutput := make(chan string)

	go func() {
		var buf [1024]byte
		n, _ := r.Read(buf[:])
		capturedOutput <- string(buf[:n])
	}()

	os.Stdout = old

	return <-capturedOutput
}

func TestLog(t *testing.T) {
	l := logger.New()

	output := captureOutput(func() {
		l.Log("Test Log Message")
	})

	expectedOutput := "[Log] Test Log Message\n"
	if output != expectedOutput {
		t.Errorf("Log() output is incorrect, got: %s, want: %s", output, expectedOutput)
	}
}

func TestError(t *testing.T) {
	l := logger.New()

	output := captureOutput(func() {
		l.Error("Test Error Message")
	})

	expectedOutput := "[Error] Test Error Message\n"
	if output != expectedOutput {
		t.Errorf("Error() output is incorrect, got: %s, want: %s", output, expectedOutput)
	}
}

func TestWarn(t *testing.T) {
	l := logger.New()

	output := captureOutput(func() {
		l.Warn("Test Warning Message")
	})

	expectedOutput := "[Warning] Test Warning Message\n"
	if output != expectedOutput {
		t.Errorf("Warn() output is incorrect, got: %s, want: %s", output, expectedOutput)
	}
}
