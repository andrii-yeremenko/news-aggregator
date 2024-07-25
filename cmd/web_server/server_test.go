package web_server

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestServerBuilder tests the ServerBuilder functionality.
func TestServerBuilder(t *testing.T) {
	// Test NewServerBuilder
	builder := NewServerBuilder()
	if builder == nil {
		t.Fatalf("expected non-nil builder")
	}

	// Test SetPort
	port := "8080"
	builder.SetPort(port)
	if builder.port != port {
		t.Errorf("expected port %s, got %s", port, builder.port)
	}

	// Test AddHandler
	path := "/test"
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	builder.AddHandler(path, handlerFunc)
	if len(builder.handlers) != 1 {
		t.Fatalf("expected 1 handler, got %d", len(builder.handlers))
	}

	// Test Build
	server := builder.Build()
	if server == nil {
		t.Fatalf("expected non-nil server")
	}
	if server.Addr != ":"+port {
		t.Errorf("expected address :%s, got %s", port, server.Addr)
	}
	if server.TLSConfig == nil || server.TLSConfig.MinVersion != tls.VersionTLS12 {
		t.Errorf("expected TLSConfig with MinVersion TLS12, got %+v", server.TLSConfig)
	}

	// Start the server in a goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			return // The server was closed
		}
	}()
	defer func(server *http.Server) {
		_ = server.Close()
	}(server)

	// Test if the /test handler is registered and working
	req := httptest.NewRequest("GET", "http://localhost:"+port+path, nil)
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	// Test if the /status handler is registered and working
	req = httptest.NewRequest("GET", "http://localhost:"+port+"/status", nil)
	rec = httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
