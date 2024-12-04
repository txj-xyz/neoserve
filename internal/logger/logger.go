package logger

import (
	"log/slog"
	"os"
)

// type Logger struct {
// 	level   slog.Level
// 	logger  *slog.Logger
// 	handler *slog.JSONHandler
// }

// const (
// 	DEBUG string = "DEBUG"
// 	INFO  string = "INFO"
// 	ERROR string = "ERROR"
// )
//

// func NewLogger(level string) *Logger {
// 	level = strings.ToUpper(level)
// 	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
// 	if level == ERROR {
// 		logger.SetOutput(os.Stderr)
// 	}
// 	// Print out the logging level
// 	fmt.Printf("Logger created with '%s' level.\n", level)
//
// 	return &Logger{
// 		level: level,
// 		log:   logger,
// 	}
// }

// ---------------------

type Logger struct {
	*slog.Logger // Wrap slog
	level        slog.Level
	handler      *slog.JSONHandler
}

// Create a new instance of a logger
func New() *Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)

	return &Logger{
		Logger:  logger,
		level:   slog.LevelDebug,
		handler: jsonHandler,
	}
}
