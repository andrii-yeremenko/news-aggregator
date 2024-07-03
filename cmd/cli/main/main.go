package main

import (
	"news-aggregator/cmd/cli"
	"news-aggregator/console_printer"
)

// the main is the entry point of the application.
func main() {
	c, err := cli.New()

	if err != nil {
		console_printer.New().Error(err.Error())
		return
	}

	c.ParseFlags()

	err = c.Run()
	if err != nil {
		console_printer.New().Error("Error occurred during execution")
	}
}
