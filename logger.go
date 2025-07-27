package main

import (
	"fmt"
	"io"
	"os"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorGray   = "\033[90m"
)

// Logger handles formatted output with optional colors
type Logger struct {
	verbose bool
	noColor bool
	writer  io.Writer
}

// NewLogger creates a new Logger instance
func NewLogger(verbose, noColor bool) *Logger {
	return &Logger{
		verbose: verbose,
		noColor: noColor,
		writer:  os.Stdout,
	}
}

// NewLoggerWithWriter creates a new Logger with a custom writer
func NewLoggerWithWriter(verbose, noColor bool, writer io.Writer) *Logger {
	return &Logger{
		verbose: verbose,
		noColor: noColor,
		writer:  writer,
	}
}

// colorize adds color to text if colors are enabled
func (l *Logger) colorize(color, text string) string {
	if l.noColor {
		return text
	}
	return color + text + ColorReset
}

// Success logs a success message
func (l *Logger) Success(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	prefix := l.colorize(ColorGreen, "✓")
	fmt.Fprintf(l.writer, "%s %s\n", prefix, message)
}

// Info logs an informational message
func (l *Logger) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.writer, "  %s\n", message)
}

// Warning logs a warning message
func (l *Logger) Warning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	prefix := l.colorize(ColorYellow, "⚠")
	fmt.Fprintf(l.writer, "%s %s\n", prefix, message)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	prefix := l.colorize(ColorRed, "✗")
	fmt.Fprintf(os.Stderr, "%s %s\n", prefix, message)
}

// Debug logs a debug message if verbose mode is enabled
func (l *Logger) Debug(format string, args ...interface{}) {
	if !l.verbose {
		return
	}
	message := fmt.Sprintf(format, args...)
	prefix := l.colorize(ColorGray, "[DEBUG]")
	fmt.Fprintf(l.writer, "%s %s\n", prefix, message)
}

// Verbose logs a message only if verbose mode is enabled
func (l *Logger) Verbose(format string, args ...interface{}) {
	if !l.verbose {
		return
	}
	message := fmt.Sprintf(format, args...)
	prefix := l.colorize(ColorBlue, "→")
	fmt.Fprintf(l.writer, "%s %s\n", prefix, message)
}

// FileCreated logs a file creation
func (l *Logger) FileCreated(path string) {
	l.Success("Created file: %s", path)
}

// FileSkipped logs a skipped file
func (l *Logger) FileSkipped(path string) {
	l.Info("Skipped existing file: %s", path)
}

// DirCreated logs a directory creation
func (l *Logger) DirCreated(path string) {
	l.Success("Created directory: %s", path)
}

// DirSkipped logs a skipped directory
func (l *Logger) DirSkipped(path string) {
	l.Info("Skipped existing directory: %s", path)
}

// Summary logs a summary of operations
func (l *Logger) Summary(created, skipped int) {
	l.writer.Write([]byte("\n"))
	
	if created > 0 {
		l.Success("Created %d %s", created, pluralize("item", created))
	}
	
	if skipped > 0 {
		l.Info("Skipped %d existing %s", skipped, pluralize("item", skipped))
	}
	
	if created == 0 && skipped > 0 {
		l.Info("All files and directories already exist")
	}
}


