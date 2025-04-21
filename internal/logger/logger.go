package logger

import (
	"log/slog"
	"os"
)

// ---------------------

type Logger struct {
	*slog.Logger // Wrap slog
	level        slog.Level
	handler      *slog.JSONHandler
}

// Create a new instance of a logger
func New() *Logger {
	return NewWithLevel("debug")
}

// NewWithLevel creates a logger with the specified level ("debug" or "info")
func NewWithLevel(level string) *Logger {
	var slogLevel slog.Level
	if level == "debug" {
		slogLevel = slog.LevelDebug
	} else {
		slogLevel = slog.LevelInfo
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
		AddSource: false,
	})
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)

	return &Logger{
		Logger:  logger,
		level:   slogLevel,
		handler: jsonHandler,
	}
}

