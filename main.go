package main

import (
	"news-aggregator/cmd"
)

// the main is the entry point of the application.
func main() {
	cli := cmd.New()
	cli.ParseFlags()
	cli.Run()
}
