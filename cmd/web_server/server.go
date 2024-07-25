package web_server

import (
	"crypto/tls"
	"net/http"
	"news-aggregator/cmd/web_server/handler"
)

const (
	// DefaultHttpsPort is the default port number for the server.
	DefaultHttpsPort = "8443"
)

// ServerBuilder is a builder pattern for creating a new http.Server instance.
type ServerBuilder struct {
	port     string
	handlers map[string]http.HandlerFunc
}

// NewServerBuilder creates a new ServerBuilder instance.
func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		port:     DefaultHttpsPort,
		handlers: make(map[string]http.HandlerFunc),
	}
}

// SetPort sets the port for the server.
func (sb *ServerBuilder) SetPort(port string) *ServerBuilder {
	sb.port = port
	return sb
}

// AddHandler adds a new handler to the server.
func (sb *ServerBuilder) AddHandler(path string, handler http.HandlerFunc) *ServerBuilder {
	sb.handlers[path] = handler
	return sb
}

// Build creates a new http.Server instance.
func (sb *ServerBuilder) Build() *http.Server {
	mux := http.NewServeMux()

	for path, hand := range sb.handlers {
		mux.HandleFunc(path, hand)
	}

	mux.HandleFunc("/status", handler.NewStatusHandler("1.0").Handle)

	return &http.Server{
		Addr:    ":" + sb.port,
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
}
