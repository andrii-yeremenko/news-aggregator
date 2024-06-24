package web_server

import (
	"crypto/tls"
	"net/http"
	"news-aggregator/cmd/web_server/handler"
)

type ServerBuilder struct {
	port     string
	handlers map[string]http.HandlerFunc
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		port:     "443",
		handlers: make(map[string]http.HandlerFunc),
	}
}

func (sb *ServerBuilder) SetPort(port string) *ServerBuilder {
	sb.port = port
	return sb
}

func (sb *ServerBuilder) AddHandler(path string, handler http.HandlerFunc) *ServerBuilder {
	sb.handlers[path] = handler
	return sb
}

func (sb *ServerBuilder) Build() *http.Server {
	mux := http.NewServeMux()

	for path, handler := range sb.handlers {
		mux.HandleFunc(path, handler)
	}

	mux.HandleFunc("/status", handler.NewServerStatusHandler("1.0").Handle)

	return &http.Server{
		Addr:    ":" + sb.port,
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
}
