package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Logger levels
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

// Logger formats
const (
	TextFormat = "text"
	JSONFormat = "json"
)

// Logger outputs
const (
	StdoutOutput = "stdout"
	FileOutput   = "file"
)

// Logger is the main logger structure
type Logger struct {
	level  string
	format string
	output string
	logger *log.Logger
}

// Config holds the logger configuration
type Config struct {
	Level    string
	Format   string
	Output   string
	FilePath string
}

// New creates a new logger instance
func New(config Config) (*Logger, error) {
	var output *os.File
	var err error

	if config.Output == FileOutput && config.FilePath != "" {
		output, err = os.OpenFile(config.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
	} else {
		output = os.Stdout
	}

	logger := log.New(output, "", log.LstdFlags)

	return &Logger{
		level:  strings.ToLower(config.Level),
		format: strings.ToLower(config.Format),
		output: strings.ToLower(config.Output),
		logger: logger,
	}, nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level == DebugLevel {
		l.log("DEBUG", msg, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level == DebugLevel || l.level == InfoLevel {
		l.log("INFO", msg, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level == DebugLevel || l.level == InfoLevel || l.level == WarnLevel {
		l.log("WARN", msg, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log("ERROR", msg, args...)
}

// log formats and logs a message
func (l *Logger) log(level, msg string, args ...interface{}) {
	var logMsg string
	
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	
	if l.format == JSONFormat {
		logMsg = fmt.Sprintf("{\"level\":\"%s\",\"message\":\"%s\"}", level, msg)
	} else {
		logMsg = fmt.Sprintf("[%s] %s", level, msg)
	}
	
	l.logger.Println(logMsg)
} 