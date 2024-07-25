package cli

import (
	"flag"
	"os"
	"testing"
)

// resetFlags resets the command-line flags
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

// TestNew checks if the CLI instance is created successfully.
func TestNew(t *testing.T) {
	resetFlags()
	cli, err := New("../../config/feeds_dictionary.json", "../../resources")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if cli.parserFactory == nil || cli.aggregator == nil || cli.resourceManager == nil {
		t.Fatal("Expected CLI fields to be initialized")
	}
}

// TestParseFlags checks if command line flags are parsed correctly.
func TestParseFlags(t *testing.T) {
	resetFlags()
	cli, err := New("../../config/feeds_dictionary.json", "../../resources")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	os.Args = []string{"cmd", "-sources=source1,source2", "-keywords=keyword1,keyword2", "-date-start=2024-01-01", "-date-end=2024-12-31", "-sort-order=asc"}
	cli.ParseFlags()

	if cli.sourceArg != "source1,source2" {
		t.Errorf("Expected sourceArg to be 'source1,source2', got '%v'", cli.sourceArg)
	}
	if cli.keywordsArg != "keyword1,keyword2" {
		t.Errorf("Expected keywordsArg to be 'keyword1,keyword2', got '%v'", cli.keywordsArg)
	}
	if cli.startDateArg != "2024-01-01" {
		t.Errorf("Expected startDateArg to be '2024-01-01', got '%v'", cli.startDateArg)
	}
	if cli.endDateArg != "2024-12-31" {
		t.Errorf("Expected endDateArg to be '2024-12-31', got '%v'", cli.endDateArg)
	}
	if cli.sortOrderArg != "asc" {
		t.Errorf("Expected sortOrderArg to be 'asc', got '%v'", cli.sortOrderArg)
	}
}

// TestRun checks the Run method of CLI.
func TestRun(t *testing.T) {
	resetFlags()
	cli, err := New("../../config/feeds_dictionary.json", "../../resources")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	os.Args = []string{"cmd"}
	cli.ParseFlags()

	err = cli.Run()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// TestRunWithParams tests the Run method with parameters.
func TestRunWithParams(t *testing.T) {
	resetFlags()
	cli, err := New("../../config/feeds_dictionary.json", "../../resources")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	os.Args = []string{"-sources=source", "-keywords=keyword", "-date-start=2024-01-01", "-date-end=2024-12-31", "-sort-order=asc"}
	cli.ParseFlags()

	err = cli.Run()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
