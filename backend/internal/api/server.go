package api

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"calculator/internal/calculator"
)

func NewRouter(registry *calculator.Registry) http.Handler {
	handlers := NewHandlers(registry)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /api/v1/operations", handlers.Operations)
	mux.HandleFunc("POST /api/v1/calculate", handlers.Calculate)
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
	server := &http.Server{
		Addr:              address,
		Handler:           NewRouter(registry),
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
