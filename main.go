package main

import (
	"news-aggregator/cmd"
	"news-aggregator/print"
)

// the main is the entry point of the application.
func main() {
	cli, err := cmd.New()

	printer := print.New()

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
