package main

import (
	"NewsAggregator/cmd"
)

// the main is the entry point of the application.
func main() {
	cli := cmd.New()
	cli.ParseFlags()
	cli.Run()
}
