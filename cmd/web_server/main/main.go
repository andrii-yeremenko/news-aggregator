package main

import (
	"log"
	"news-aggregator/cmd/web_server"
	"news-aggregator/cmd/web_server/handler"
	"news-aggregator/resource_manager"
	"os"
	"path"
)

func main() {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	managerConfigPath := path.Join(basePath, "/config/feeds_dictionary.json")
	storagePath := path.Join(basePath, "/resources")

	manager, err := resource_manager.New(storagePath, managerConfigPath)

	if err != nil {
		log.Fatalf("failed to create resource manager: %v", err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8443"
	}

	server := web_server.NewServerBuilder().
		SetPort(port).
		AddHandler("/news", handler.NewNewsHandler(manager).Handle).
		AddHandler("/update", handler.NewUpdateHandler(manager).Handle).
		AddHandler("/sources", handler.NewFeedsManagerHandler(manager).Handle).
		Build()

	log.Println("Starting server on port " + port + " ...")

	err = server.ListenAndServeTLS("certificates/cert.pem", "certificates/key.pem")
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
		return
	}

	log.Println("Server started successfully")
}
