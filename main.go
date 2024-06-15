package main

import (
	"news-aggregator/cmd"
	"news-aggregator/console_printer"
)

// the main is the entry point of the application.
func main() {
	cli, err := cmd.New()

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
