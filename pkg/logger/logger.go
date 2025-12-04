package logger

import (
	"log/slog"
	"os"
	"strings"

	"koda/pkg/env"

	"github.com/lmittmann/tint"
)

var (
	logger *slog.Logger
)

func ParseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func SetLevel(level string) {
	parsedLevel := ParseLevel(level)
	w := os.Stdout

	logger = slog.New(tint.NewHandler(w, &tint.Options{
		Level:     parsedLevel,
		AddSource: true,
	}))

	slog.SetDefault(logger)
}

func InitDefaultLogger(name string, version string) {
	logger.Info("initializing service",
		"name", name,
		"version", version,
		"environment", env.LookupOrDefault("ENV", "development"),
		"log_level", env.LookupOrDefault("LOG_LEVEL", "info"),
	)
}
