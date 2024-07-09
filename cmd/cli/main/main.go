package main

import (
	"news-aggregator/cmd/cli"
	"news-aggregator/print"
)

// the main is the entry point of the application.
func main() {
	c, err := cli.New("/config/feeds_dictionary.json", "/resources")

	printer := print.New()

	if err != nil {
		printer.Error(err.Error())
		return
	}

	c.ParseFlags()

	err = c.Run()
	if err != nil {
		printer.Error("Error occurred during execution")
	}
}
