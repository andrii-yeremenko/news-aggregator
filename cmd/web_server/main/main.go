package main

import (
	"log"
	"news-aggregator/aggregator"
	"news-aggregator/cmd/web_server"
	"news-aggregator/cmd/web_server/handler"
	"news-aggregator/resource_manager"
	"os"
	"path"
)

func main() {
	parserPool := aggregator.NewParserFactory()
	a, err := aggregator.New(parserPool)
	if err != nil {
		log.Fatalf("failed to create aggregator: %v", err)
	}

	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	manager := resource_manager.New(path.Join(basePath, "/storage"))

	newsHandler := handler.NewNewsHandler(a, manager)

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
