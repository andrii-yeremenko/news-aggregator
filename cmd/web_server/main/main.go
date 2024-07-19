package main

import (
	"log"
	"news-aggregator/cmd/web_server"
	"news-aggregator/cmd/web_server/handler"
	"news-aggregator/manager"
	"os"
	"path"
	"strconv"
	"time"
)

func main() {
	basePath := getCurrentDirectory()

	m := createResourceManager(basePath)

	timeoutStr := os.Getenv("TIMEOUT")
	if timeoutStr == "" {
		log.Println("TIMEOUT environment variable not set. Using default timeout of 12 hours.")
		timeoutStr = "12h"
	}

	timeout, err := time.ParseDuration(timeoutStr)

	if err != nil {
		log.Fatalf("Failed to parse TIMEOUT duration: %v", err)
	}

	scheduler := web_server.NewUpdateScheduler(m, timeout)
	scheduler.Start()

	p := getPort()

	startServer(p, m)
}

// getCurrentDirectory retrieves the current working directory.
func getCurrentDirectory() string {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}
	return basePath
}

// createResourceManager initializes and returns the resource manager.
func createResourceManager(basePath string) *manager.ResourceManager {
	managerConfigPath := path.Join(basePath, "/config/feeds_dictionary.json")
	storagePath := path.Join(basePath, "/resources")

	m, err := manager.New(storagePath, managerConfigPath)
	if err != nil {
		log.Fatalf("failed to create resource manager: %v", err)
	}
	return m
}

// getPort returns the port number to use for the server.
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

// startServer initializes and starts the web server.
func startServer(port string, m *manager.ResourceManager) {
	server := web_server.NewServerBuilder().
		SetPort(port).
		AddHandler("/news", handler.NewNewsHandler(m).Handle).
		AddHandler("/update", handler.NewUpdateHandler(m).Handle).
		AddHandler("/sources", handler.NewFeedsManagerHandler(m).Handle).
		Build()

	log.Println("Starting server on port " + port + " ...")

	err := server.ListenAndServeTLS("certificates/cert.pem", "certificates/key.pem")
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
	log.Println("Server started successfully")
}
