package main

import (
	"fmt"
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
	basePath, err := getCurrentDirectory()

	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	m, err := createResourceManager(basePath)

	if err != nil {
		log.Fatalf("failed to create resource manager: %v", err)
	}

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

	p, err := getPort()

	if err != nil {
		log.Println(err)
	}

	startServer(p, m)
}

// getCurrentDirectory retrieves the current working directory.
func getCurrentDirectory() (string, error) {
	basePath, err := os.Getwd()
	return basePath, err
}

// createResourceManager initializes and returns the resource manager.
func createResourceManager(basePath string) (*manager.ResourceManager, error) {
	managerConfigPath := path.Join(basePath, "/config/feeds_dictionary.json")
	storagePath := path.Join(basePath, "/resources")

	return manager.New(storagePath, managerConfigPath)
}

// getPort returns the port number to use for the server.
func getPort() (string, error) {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8443"
		fmt.Println("PORT environment variable not set. Using default port 8443.")
		return port, nil
	} else {
		_, err := strconv.Atoi(port)
		return "", err
	}
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
