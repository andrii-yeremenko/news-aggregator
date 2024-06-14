package main

import (
	"news-aggregator/cmd"
	"news-aggregator/logger"
)

// the main is the entry point of the application.
func main() {
	cli, err := cmd.New()

	if err != nil {
		logger.New().Error(err.Error())
		return
	}

	cli.ParseFlags()

	err = cli.Run()
	if err != nil {
		logger.New().Error("Error occurred during execution")
	}
}
