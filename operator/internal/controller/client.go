package controller

import (
	"crypto/tls"
	"io"
	"net/http"
)

type DefaultHTTPClient struct {
	client *http.Client
}

func NewDefaultHTTPClient() *DefaultHTTPClient {
	return &DefaultHTTPClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *DefaultHTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return c.client.Post(url, contentType, body)
}
