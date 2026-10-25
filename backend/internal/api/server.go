package api

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"calculator/internal/calculator"
)

func NewRouter(registry *calculator.Registry) http.Handler {
	return NewRouterWithStatic(registry, "")
}

func NewRouterWithStatic(registry *calculator.Registry, staticDir string) http.Handler {
	handlers := NewHandlers(registry)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /api/v1/operations", handlers.Operations)
	mux.HandleFunc("POST /api/v1/calculate", handlers.Calculate)
	if staticDir != "" {
		static := http.FileServer(http.Dir(staticDir))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/healthz" {
				http.NotFound(w, r)
				return
			}
			relativePath := strings.TrimPrefix(filepath.Clean(r.URL.Path), string(filepath.Separator))
			filePath := filepath.Join(staticDir, relativePath)
			if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
				static.ServeHTTP(w, r)
				return
			}
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
		})
	}
	return withMiddleware(routeErrors(mux))
}

func routeErrors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) < 5 || r.URL.Path[:5] != "/api/" {
			next.ServeHTTP(w, r)
			return
		}
		buffer := &responseBuffer{header: make(http.Header), status: http.StatusOK}
		next.ServeHTTP(buffer, r)
		if buffer.status == http.StatusNotFound {
			writeError(w, http.StatusNotFound, ErrCodeNotFound, "resource not found")
			return
		}
		if buffer.status == http.StatusMethodNotAllowed {
			if allow := buffer.header.Get("Allow"); allow != "" {
				w.Header().Set("Allow", allow)
			}
			writeError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed")
			return
		}
		for key, values := range buffer.header {
			w.Header()[key] = values
		}
		w.WriteHeader(buffer.status)
		_, _ = w.Write(buffer.body.Bytes())
	})
}

type responseBuffer struct {
	header http.Header
	status int
	body   bytes.Buffer
}

func (b *responseBuffer) Header() http.Header            { return b.header }
func (b *responseBuffer) WriteHeader(status int)         { b.status = status }
func (b *responseBuffer) Write(body []byte) (int, error) { return b.body.Write(body) }

func withMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		defer func() {
			if recovered := recover(); recovered != nil {
				slog.Error("request panic recovered", "panic", recovered, "path", r.URL.Path)
				writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Serve(ctx context.Context, address string, registry *calculator.Registry) error {
	return serve(ctx, address, registry, "")
}

func ServeWithStatic(ctx context.Context, address string, registry *calculator.Registry, staticDir string) error {
	return serve(ctx, address, registry, staticDir)
}

func serve(ctx context.Context, address string, registry *calculator.Registry, staticDir string) error {
	server := &http.Server{
		Addr:              address,
		Handler:           NewRouterWithStatic(registry, staticDir),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errorsFromServer := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			errorsFromServer <- nil
			return
		}
		errorsFromServer <- err
	}()

	select {
	case err := <-errorsFromServer:
		return err
	case <-ctx.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownContext)
	}
}
