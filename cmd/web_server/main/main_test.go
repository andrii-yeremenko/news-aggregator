package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestStartServer tests the main function's ability to start the server.
func TestStartServer(t *testing.T) {

	err := os.Setenv("PORT", "8443")
	if err != nil {
		t.Fatalf("failed to set PORT environment variable: %v", err)
	}
	err = os.Setenv("TIMEOUT", "3s")
	if err != nil {
		t.Fatalf("failed to set TIMEOUT environment variable: %v", err)
	}
	err = os.Setenv("MANAGER_CONFIG_PATH", "../../config/feeds_dictionary.json")
	if err != nil {
		t.Fatalf("failed to set MANAGER_CONFIG_PATH environment variable: %v", err)
	}
	err = os.Setenv("STORAGE_PATH", "../../resources")
	if err != nil {
		t.Fatalf("failed to set STORAGE_PATH environment variable: %v", err)
	}
	err = os.Setenv("CERT_FILE_PATH", "../../../certificates/cert.pem")
	if err != nil {
		t.Fatalf("failed to set CERT_FILE_PATH environment variable: %v", err)
	}
	err = os.Setenv("KEY_FILE_PATH", "../../../certificates/key.pem")
	if err != nil {
		t.Fatalf("failed to set KEY_FILE_PATH environment variable: %v", err)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("main function panicked: %v", r)
			}
		}()
		main()
	}()

	time.Sleep(3 * time.Second)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get("https://localhost:8443/news")
	if err != nil {
		t.Fatalf("server did not start: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("failed to close response body: %v", err)
		}
	}(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
