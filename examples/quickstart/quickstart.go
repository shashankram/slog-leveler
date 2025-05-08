package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/shashankram/slog-leveler/pkg/logger"
)

func main() {
	// Default logger
	slog.Info("Hello world")
	slog.Debug("This won't be printed")
	logger.SetLevel(logger.DefaultComponent, slog.LevelDebug) // nolint: errcheck
	slog.Debug("This will be printed")
	fmt.Println()

	// Custom logger
	fooLogger := logger.NewWithOptions("foo", logger.Options{
		Format:    logger.JSONFormat,
		AddSource: true,
	})
	fooLogger.Info("Hello foo")
	fooLogger.Debug("This won't be printed")
	logger.SetLevel("foo", slog.LevelDebug) // nolint: errcheck
	fooLogger.Debug("This will be printed")
	fooLogger.Log(context.Background(), logger.LevelTrace, "This won't be printed")
	logger.SetLevel("foo", logger.LevelTrace) // nolint: errcheck
	fooLogger.Log(context.Background(), logger.LevelTrace, "This will be printed")

	// Managing short-lived loggers
	for range 3 {
		someFunc()
	}
}

func someFunc() {
	l := logger.New("tmp")
	defer logger.DeleteLeveler("tmp") // nolint: errcheck
	l.Info("tmp info")
}
