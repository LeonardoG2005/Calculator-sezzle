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
	if err := api.Serve(ctx, ":"+port, calculator.NewDefaultRegistry()); err != nil {
		slog.Error("calculator server stopped with error", "error", err)
		os.Exit(1)
	}
}
