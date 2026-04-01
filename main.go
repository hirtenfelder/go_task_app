package main

import (
	"awesomeProject/db"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "awesomeProject/task"
)

var (
	serverPort string
)

func init() {
	serverPort = os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}
}

func main() {
	addr := fmt.Sprintf(":%s", serverPort)
	server := &http.Server{Addr: addr}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	defer close(stop)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start a server in a goroutine
	go func() {
		slog.Info("Starting server on port: " + serverPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error: %v", err)
		}
	}()

	// Wait for the interrupt signal
	<-stop
	slog.Info("Shutdown signal received, closing resources...")

	// Close database connection here
	db.Close()

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown error: %v", err)
	}

	slog.Info("Server stopped gracefully")
}
