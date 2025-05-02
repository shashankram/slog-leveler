package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shashankram/slog-leveler/pkg/logger"
)

func main() {
	exitChan := make(chan os.Signal, 1)
	// Notify the channel of SIGINT and SIGTERM signals.
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)

	fooLogger := logger.New("foo")
	barLogger := logger.NewWithOptions("bar", logger.Options{
		AddSource: true,
		Level:     logger.LevelTrace,
	})

	go func() {
		for {
			select {
			case <-exitChan:
				slog.Info("received exit signal, exiting")
				return
			default:
				slog.Info("default info")
				slog.Debug("default debug")
				slog.Log(context.Background(), logger.LevelTrace, "default trace")
				fooLogger.Info("foo info")
				fooLogger.Debug("foo debug")
				fooLogger.Log(context.Background(), logger.LevelTrace, "foo trace")
				barLogger.Info("bar info")
				barLogger.Debug("bar debug")
				barLogger.Log(context.Background(), logger.LevelTrace, "bar trace")
				time.Sleep(5 * time.Second)
				fmt.Println()
			}
		}
	}()

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":15000",
		Handler: mux,
	}

	mux.HandleFunc("/logging", logger.HTTPLevelHandler)

	go func() {
		fmt.Println("Server listening on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-exitChan

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown server", "error", err)
	}
	slog.Info("exiting")
}
