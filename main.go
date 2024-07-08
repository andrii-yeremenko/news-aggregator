package main

import (
	"news-aggregator/cmd"
	"news-aggregator/console_printer"
)

// the main is the entry point of the application.
func main() {
	cli, err := cmd.New()

	printer := console_printer.New()

	if err != nil {
		printer.Error(err.Error())
		return
	}

	cli.ParseFlags()

	err = cli.Run()
	if err != nil {
		printer.Error("Error occurred during execution")
	}
}
