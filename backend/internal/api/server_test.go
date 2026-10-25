package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"calculator/internal/calculator"
)

func TestRouterRoutesAndErrors(t *testing.T) {
	handler := NewRouter(calculator.NewDefaultRegistry())

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantType   string
	}{
		{"health", http.MethodGet, "/healthz", http.StatusOK, "application/json"},
		{"operations", http.MethodGet, "/api/v1/operations", http.StatusOK, "application/json"},
		{"unknown api route", http.MethodGet, "/api/v1/unknown", http.StatusNotFound, "application/json"},
		{"wrong method", http.MethodGet, "/api/v1/calculate", http.StatusMethodNotAllowed, "application/json"},
		{"cors preflight", http.MethodOptions, "/api/v1/calculate", http.StatusNoContent, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if test.wantType != "" && response.Header().Get("Content-Type") != test.wantType {
				t.Fatalf("content type = %q, want %q", response.Header().Get("Content-Type"), test.wantType)
			}
		})
	}
}

func TestServeStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Serve(ctx, "127.0.0.1:0", calculator.NewDefaultRegistry()); err != nil {
		t.Fatalf("Serve() error = %v", err)
	}
}

func TestServeStartsAndStops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- Serve(ctx, "127.0.0.1:0", calculator.NewDefaultRegistry())
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Serve() error = %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		cancel()
		select {
		case err := <-done:
			if err != nil {
				t.Fatalf("Serve() shutdown error = %v", err)
			}
		case <-time.After(time.Second):
			t.Fatal("Serve() did not stop after cancellation")
		}
	}
}

func TestRouterServesStaticFilesAndSPA(t *testing.T) {
	staticDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(staticDir, "index.html"), []byte("app"), 0o600); err != nil {
		t.Fatal(err)
	}
	handler := NewRouterWithStatic(calculator.NewDefaultRegistry(), staticDir)

	for _, path := range []string{"/", "/history"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK || response.Body.String() != "app" {
			t.Fatalf("%s: status = %d, body = %q", path, response.Code, response.Body.String())
		}
	}
}
