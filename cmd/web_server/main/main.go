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

// Default values for environment variables.
const (
	// DefaultTimeout is the default timeout for the update scheduler.
	DefaultTimeout = "12h"

	// DefaultPort is the default port number for the server.
	DefaultPort = "8443"

	// DefaultManagerConfigPath is the default path to the feeds dictionary configuration file.
	DefaultManagerConfigPath = "config/feeds_dictionary.json"

	// DefaultStoragePath is the default path to the storage directory.
	DefaultStoragePath = "resources"

	// DefaultCertFilePath is the default path to the certificate file.
	DefaultCertFilePath = "/etc/tls/tls.crt"

	// DefaultKeyFilePath is the default path to the key file.
	DefaultKeyFilePath = "/etc/tls/tls.key"
)

func main() {
	basePath, err := getCurrentDirectory()

	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	managerConfigPath := getEnv("MANAGER_CONFIG_PATH", path.Join(basePath, DefaultManagerConfigPath))
	storagePath := getEnv("STORAGE_PATH", path.Join(basePath, DefaultStoragePath))
	m, err := createResourceManager(managerConfigPath, storagePath)

	if err != nil {
		log.Fatalf("failed to create resource manager: %v", err)
	}

	timeoutStr := getEnv("TIMEOUT", DefaultTimeout)
	timeout, err := time.ParseDuration(timeoutStr)

	if err != nil {
		log.Fatalf("Failed to parse TIMEOUT duration: %v", err)
	}

	scheduler := web_server.NewUpdateScheduler(m, timeout)
	scheduler.Start()

	port, err := getPort()

	if err != nil {
		log.Println(err)
	}

	certFilePath := getEnv("CERT_FILE_PATH", DefaultCertFilePath)
	keyFilePath := getEnv("KEY_FILE_PATH", DefaultKeyFilePath)

	startServer(port, certFilePath, keyFilePath, m)
}

// getCurrentDirectory retrieves the current working directory.
func getCurrentDirectory() (string, error) {
	basePath, err := os.Getwd()
	return basePath, err
}

// createResourceManager initializes and returns the resource manager.
func createResourceManager(managerConfigPath, storagePath string) (*manager.ResourceManager, error) {
	return manager.New(storagePath, managerConfigPath)
}

// getPort returns the port number to use for the server.
func getPort() (string, error) {
	port := getEnv("PORT", DefaultPort)
	if _, err := strconv.Atoi(port); err != nil {
		return "", err
	}
	return port, nil
}

// startServer initializes and starts the web server.
func startServer(port, certFilePath, keyFilePath string, m *manager.ResourceManager) {
	server := web_server.NewServerBuilder().
		SetPort(port).
		AddHandler("/news", handler.NewNewsHandler(m).Handle).
		AddHandler("/sources", handler.NewFeedsManagerHandler(m).Handle).
		AddHandler("/availableFeeds", handler.NewAvailableFeedsHandler(m).Handle).
		Build()

	log.Println("Starting server on port " + port + " ...")

	err := server.ListenAndServeTLS(certFilePath, keyFilePath)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
	log.Println("Server started successfully")
}

// getEnv retrieves the value of the environment variable named by the key or returns the default value if the
// variable is not present.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
