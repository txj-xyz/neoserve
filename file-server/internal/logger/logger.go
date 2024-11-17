package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Logger struct {
	level string
	log   *log.Logger
}

const (
	DEBUG string = "DEBUG"
	INFO  string = "INFO"
	ERROR string = "ERROR"
)

func NewLogger(level string) *Logger {
	level = strings.ToUpper(level)
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	if level == ERROR {
		logger.SetOutput(os.Stderr)
	}
	// Print out the logging level
	fmt.Printf("Logger created with '%s' level.\n", level)

	return &Logger{
		level: level,
		log:   logger,
	}
}

func (l *Logger) Logf(format string, args ...any) {
	switch l.level {
	case DEBUG:
		l.log.Printf("[DEBUG]: "+format, args...)
	case INFO:
		l.log.Printf("[INFO]: "+format, args...)
	case ERROR:
		l.log.Printf("[ERROR]: "+format, args...)
	default:
		l.log.Printf("[LOG]: "+format, args...)
	}
}

func (l *Logger) Debug(format string, args ...any) {
	if l.level == DEBUG {
		l.log.Printf("[DEBUG]: "+format, args...)
	}
}

func (l *Logger) Info(message string) {
	if l.level == INFO {
		l.log.Println("INFO: " + message)
	}
}

func (l *Logger) Error(message string) {
	if l.level == ERROR {
		l.log.Println("ERROR: " + message)
	}
}
