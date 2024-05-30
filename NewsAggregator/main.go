package main

import "NewsAggregator/cmd"

// main is the entry point of the application.
func main() {
	cli := cmd.NewCLI()
	cli.ParseFlags()
	cli.Run()
}
