package main

import (
	"log"
	"news-aggregator/cmd/web_server"
	"news-aggregator/cmd/web_server/handler"
	"news-aggregator/manager"
	"os"
	"path"
	"strconv"
)

func main() {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	managerConfigPath := path.Join(basePath, "/config/feeds_dictionary.json")
	storagePath := path.Join(basePath, "/resources")

	m, err := manager.New(storagePath, managerConfigPath)

	if err != nil {
		log.Fatalf("failed to create resource m: %v", err)
	}

	port := getPort()

	server := web_server.NewServerBuilder().
		SetPort(port).
		AddHandler("/news", handler.NewNewsHandler(m).Handle).
		AddHandler("/update", handler.NewUpdateHandler(m).Handle).
		AddHandler("/sources", handler.NewFeedsManagerHandler(m).Handle).
		Build()

	log.Println("Starting server on port " + port + " ...")

	err = server.ListenAndServeTLS("certificates/cert.pem", "certificates/key.pem")
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
		return
	}

	log.Println("Server started successfully")
}

// getPort returns the port number from the environment variable, or the default one if port not specified or invalid.
func getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8443"
		log.Println("PORT environment variable not set. Using default port " + port)
	} else {
		_, err := strconv.Atoi(port)
		if err != nil {
			log.Println("Invalid PORT value set in environment variable. Using default port " + port)
			port = "8443"
		}
	}

	return port
}
