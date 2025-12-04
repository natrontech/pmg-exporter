package graceful

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server interface {
	Label() string
	Start() error
	Shutdown(context.Context) error
}

func Context() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan

		slog.Info("shutdown signal received, cancelling context")
		cancel()
	}()

	return ctx, cancel
}

func Serve(ctx context.Context, servers ...Server) {
	var wg sync.WaitGroup
	wg.Add(len(servers))

	for _, server := range servers {
		go runServer(ctx, &wg, server)
	}

	// Wait for all servers to shut down
	wg.Wait()
	slog.Info("all servers stopped successfully")
}

func runServer(ctx context.Context, wg *sync.WaitGroup, server Server) {
	defer wg.Done()

	errChan := make(chan error, 1)

	// start the server in a goroutine
	go func() {
		slog.Info(fmt.Sprintf("Starting %s server...", server.Label()))
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("%s server error: %w", server.Label(), err)
		}
	}()

	// wait for either context cancellation or server error
	select {
	case <-ctx.Done():
		handleShutdown(server, "Context cancelled")
	case err := <-errChan:
		handleShutdown(server, fmt.Sprintf("Server error: %v", err))
	}
}

func handleShutdown(server Server, reason string) {
	slog.Info("Shutting down server", "server", server.Label(), "reason", reason)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to gracefully shut down server", "server", server.Label(), "error", err)
	}
}
