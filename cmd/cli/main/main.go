package main

import (
	"news-aggregator/cmd/cli"
	"news-aggregator/console_printer"
)

// the main is the entry point of the application.
func main() {
	cli, err := cli.New()

	if err != nil {
		console_printer.New().Error(err.Error())
		return
	}

	cli.ParseFlags()

	err = cli.Run()
	if err != nil {
		console_printer.New().Error("Error occurred during execution")
	}
}
