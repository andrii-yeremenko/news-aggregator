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

	manager, err := resource_manager.New(path.Join(basePath, "/storage"),
		path.Join(basePath, "/resource_manager/setup/resources.json"))

	if err != nil {
		log.Fatalf("failed to create resource manager: %v", err)
	}

	server := web_server.NewServerBuilder().
		SetPort("8443").
		AddHandler("/news", handler.NewNewsHandler(manager).Handle).
		AddHandler("/update", handler.NewUpdateHandler(manager).Handle).
		AddHandler("/sources", handler.NewControlHandler(manager).Handle).
		Build()

	log.Println("Starting server on port 8443")
	err = server.ListenAndServeTLS("cmd/web_server/cert/cert.pem", "cmd/web_server/cert/key.pem")
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
