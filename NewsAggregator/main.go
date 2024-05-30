package main

import (
	"NewsAggregator/cmd"
	"NewsAggregator/logger"
)

// main is the entry point of the application.
func main() {
	logger.New().Log("Starting NewsAggregator")
	cli := cmd.New()
	cli.ParseFlags()
	cli.Run()
}
