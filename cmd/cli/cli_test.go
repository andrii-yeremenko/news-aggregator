package cli

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

// originalDir is the original testing directory
var originalDir string

// resetFlags resets the command-line flags
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

// changeToProjectRoot changes the current working directory to the project root
func changeToProjectRoot() error {
	var err error
	originalDir, err = os.Getwd()
	if err != nil {
		return err
	}

	dir, err := filepath.Abs("../../")
	if err != nil {
		return err
	}
	return os.Chdir(dir)
}

// returnToTestDir changes the current working directory to the test directory
func returnToTestDir() error {
	if originalDir == "" {
		return nil
	}
	return os.Chdir(originalDir)
}

// TestNew checks if the CLI instance is created successfully.
// This test runs in the project root directory to test the relative paths.
func TestNew(t *testing.T) {
	resetFlags()
	if err := changeToProjectRoot(); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}
	cli, err := New("/config/feeds_dictionary.json", "/resources")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if cli.parserFactory == nil || cli.aggregator == nil || cli.resourceManager == nil {
		t.Fatal("Expected CLI fields to be initialized")
	}
	if returnToTestDir() != nil {
		t.Fatalf("Failed to return to test directory")
	}
}

// TestParseFlags checks if command line flags are parsed correctly.
// This test runs in the project root directory to test the relative paths.
func TestParseFlags(t *testing.T) {
	resetFlags()
	if err := changeToProjectRoot(); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}
	cli, err := New("config/feeds_dictionary.json", "resources")
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
	if returnToTestDir() != nil {
		t.Fatalf("Failed to return to test directory")
	}
}

// TestRun checks the Run method of CLI.
// This test runs in the project root directory to test the relative paths.
func TestRun(t *testing.T) {
	resetFlags()
	if err := changeToProjectRoot(); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}
	cli, err := New("config/feeds_dictionary.json", "resources")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	os.Args = []string{"cmd"}
	cli.ParseFlags()

	err = cli.Run()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if returnToTestDir() != nil {
		t.Fatalf("Failed to return to test directory")
	}
}

// TestRunWithParams tests the Run method with parameters.
// This test runs in the project root directory to test the relative paths.
func TestRunWithParams(t *testing.T) {
	resetFlags()
	if err := changeToProjectRoot(); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}
	cli, err := New("config/feeds_dictionary.json", "resources")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	os.Args = []string{"-sources=source", "-keywords=keyword", "-date-start=2024-01-01", "-date-end=2024-31-12", "-sort-order=asc"}
	cli.ParseFlags()

	err = cli.Run()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if returnToTestDir() != nil {
		t.Fatalf("Failed to return to test directory")
	}
}
