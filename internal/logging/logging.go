package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

var logger *log.Logger
var logFile *os.File

// Init sets up file logging to ~/.kbd/kbd.log.
func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ".kbd")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(dir, "kbd.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	logFile = f
	logger = log.New(f, "", log.LstdFlags)
	logger.Printf("--- kbd started at %s ---", time.Now().Format(time.RFC3339))
	return nil
}

// Close flushes and closes the log file.
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// Info logs an informational message.
func Info(format string, args ...any) {
	if logger != nil {
		logger.Printf("[INFO] "+format, args...)
	}
}

// Error logs an error message.
func Error(format string, args ...any) {
	if logger != nil {
		logger.Printf("[ERROR] "+format, args...)
	}
}

// RecoverPanic captures panics, logs them, and prints to stderr.
// Call as: defer logging.RecoverPanic()
func RecoverPanic() {
	if r := recover(); r != nil {
		stack := string(debug.Stack())
		msg := fmt.Sprintf("PANIC: %v\n%s", r, stack)
		Error("%s", msg)
		// Also print to stderr so user sees it after terminal restores
		fmt.Fprintf(os.Stderr, "\nkbd crashed. Details logged to ~/.kbd/kbd.log\n\n%s\n", msg)
	}
}
