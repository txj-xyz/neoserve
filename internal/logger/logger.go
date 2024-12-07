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
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: false,
	})
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)

	return &Logger{
		Logger:  logger,
		level:   slog.LevelDebug,
		handler: jsonHandler,
	}
}
