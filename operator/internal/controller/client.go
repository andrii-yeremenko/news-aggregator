package controller

import (
	"crypto/tls"
	"io"
	"net/http"
)

// DefaultHTTPClient is an interface for making HTTP requests.
type DefaultHTTPClient struct {
	client *http.Client
}

// NewDefaultHTTPClient creates a new DefaultHTTPClient.
func NewDefaultHTTPClient() *DefaultHTTPClient {
	return &DefaultHTTPClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

// Do sends an HTTP request and returns an HTTP response.
func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// Post sends an HTTP POST request and returns an HTTP response.
func (c *DefaultHTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return c.client.Post(url, contentType, body)
}
