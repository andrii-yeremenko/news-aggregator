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

	manager := resource_manager.New(path.Join(basePath, "/storage"))

	newsHandler := handler.NewNewsHandler(manager)

	server := web_server.NewServerBuilder().
		SetPort("8443").
		AddHandler("/news", newsHandler.Handle).
		Build()

	log.Println("Starting server on port 8443")
	err = server.ListenAndServeTLS("cmd/web_server/cert/cert.pem", "cmd/web_server/cert/key.pem")
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
