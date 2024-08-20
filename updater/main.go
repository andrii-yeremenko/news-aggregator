package main

import (
	"flag"
	"log"
	"updater/storage"
	"updater/updater"
)

const (
	defaultResourcesPath = "./resources"
	defaultFeedsConfig   = "./config/feeds_dictionary.json"
)

func main() {
	resource := flag.String("resource", "", "[Optional] Name of the resource to update")
	feedsConfig := flag.String("feeds-config", "", "[Optional] Path to the feeds config file")
	resourcesPath := flag.String("resources-path", "", "[Optional] Path to the resources directory")
	flag.Usage = printUsage
	flag.Parse()

	if *feedsConfig == "" {
		*feedsConfig = defaultFeedsConfig
	}

	if *resourcesPath == "" {
		*resourcesPath = defaultResourcesPath
	}

	s, err := storage.New(*resourcesPath)

	if err != nil {
		log.Fatalf("Error of storage creation: %v", err)
	}

	u, err := updater.New(*feedsConfig, s)

	if err != nil {
		log.Fatalf("Error of updater creation: %v", err)
	}

	if *resource == "" {
		u.UpdateAllFeeds()
	} else {
		err := u.UpdateFeed(*resource)
		if err != nil {
			log.Fatalf("Error of resource updation: %v", err)
		}
	}

	log.Println("Update successful!")
}

func printUsage() {
	log.Println("Usage: updater [options]")
	log.Println("Options:")
	flag.PrintDefaults()
	log.Println("Pay attention: If resource is not specified, all resources will be updated!")
	log.Println("If you didn't specify the feeds-config and resources-path, the default values will be used.")
	log.Println("Example: updater -resource=example -feeds-config=feeds.json -resources-path=./resources")
}
