package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"calculator/internal/api"
	"calculator/internal/calculator"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("starting calculator server", "port", port)
	staticDir := os.Getenv("STATIC_DIR")
	if err := api.ServeWithStatic(ctx, ":"+port, calculator.NewDefaultRegistry(), staticDir); err != nil {
		slog.Error("calculator server stopped with error", "error", err)
		os.Exit(1)
	}
}
