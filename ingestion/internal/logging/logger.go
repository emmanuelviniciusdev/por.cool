package logging

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/porcool/ingestion/internal/config"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	// LevelDebug is for debug messages
	LevelDebug LogLevel = iota
	// LevelInfo is for informational messages
	LevelInfo
	// LevelWarn is for warning messages
	LevelWarn
	// LevelError is for error messages
	LevelError
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging with OpenSearch support and fallback to stdout
type Logger struct {
	opensearch   *OpenSearchClient
	stdLogger    *log.Logger
	serviceName  string
	hostname     string
	mu           sync.RWMutex
	minLevel     LogLevel
	reconnectMu  sync.Mutex
	lastReconnect time.Time
}

// NewLogger creates a new logger with OpenSearch support
func NewLogger(cfg config.OpenSearchConfig, serviceName string) *Logger {
	hostname, _ := os.Hostname()

	logger := &Logger{
		stdLogger:   log.New(os.Stdout, "", log.LstdFlags),
		serviceName: serviceName,
		hostname:    hostname,
		minLevel:    LevelInfo,
	}

	if cfg.Enabled && cfg.URL != "" {
		logger.opensearch = NewOpenSearchClient(cfg)
	}

	return logger
}

// Initialize connects to OpenSearch if configured
func (l *Logger) Initialize(ctx context.Context) error {
	if l.opensearch == nil {
		l.stdLogger.Println("OpenSearch logging disabled, using stdout only")
		return nil
	}

	if err := l.opensearch.Connect(ctx); err != nil {
		l.stdLogger.Printf("Warning: Failed to connect to OpenSearch, using fallback logging: %v", err)
		return nil // Don't fail initialization, just use fallback
	}

	l.stdLogger.Println("OpenSearch logging initialized successfully")
	return nil
}

// SetMinLevel sets the minimum log level
func (l *Logger) SetMinLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

// log is the internal logging method
func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	l.mu.RLock()
	if level < l.minLevel {
		l.mu.RUnlock()
		return
	}
	l.mu.RUnlock()

	timestamp := time.Now()

	// Always log to stdout as fallback
	if fields != nil && len(fields) > 0 {
		l.stdLogger.Printf("[%s] %s %v", level.String(), message, fields)
	} else {
		l.stdLogger.Printf("[%s] %s", level.String(), message)
	}

	// Try to log to OpenSearch if available
	if l.opensearch != nil && l.opensearch.IsAvailable() {
		entry := LogEntry{
			Timestamp: timestamp,
			Level:     level.String(),
			Message:   message,
			Service:   l.serviceName,
			Host:      l.hostname,
			Fields:    fields,
		}
		l.opensearch.IndexLogAsync(entry)
	} else if l.opensearch != nil {
		// Try to reconnect periodically
		l.tryReconnect()
	}
}

// tryReconnect attempts to reconnect to OpenSearch if enough time has passed
func (l *Logger) tryReconnect() {
	l.reconnectMu.Lock()
	defer l.reconnectMu.Unlock()

	// Only try to reconnect every 30 seconds
	if time.Since(l.lastReconnect) < 30*time.Second {
		return
	}

	l.lastReconnect = time.Now()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := l.opensearch.Reconnect(ctx); err != nil {
			// Silently fail - we're already using fallback
			return
		}
		l.stdLogger.Println("Reconnected to OpenSearch successfully")
	}()
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LevelDebug, message, f)
}

// Info logs an informational message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LevelInfo, message, f)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LevelWarn, message, f)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LevelError, message, f)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(LevelDebug, fmt.Sprintf(format, args...), nil)
}

// Infof logs a formatted informational message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(LevelInfo, fmt.Sprintf(format, args...), nil)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(LevelWarn, fmt.Sprintf(format, args...), nil)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(LevelError, fmt.Sprintf(format, args...), nil)
}

// WithFields returns a new logger entry with the specified fields
func (l *Logger) WithFields(fields map[string]interface{}) *LoggerEntry {
	return &LoggerEntry{
		logger: l,
		fields: fields,
	}
}

// Close closes the logger and its OpenSearch connection
func (l *Logger) Close() error {
	if l.opensearch != nil {
		return l.opensearch.Close()
	}
	return nil
}

// Writer returns an io.Writer that writes to the logger at INFO level
// This can be used to redirect standard library log output
func (l *Logger) Writer() io.Writer {
	return &logWriter{logger: l, level: LevelInfo}
}

// LoggerEntry represents a log entry with pre-set fields
type LoggerEntry struct {
	logger *Logger
	fields map[string]interface{}
}

// Debug logs a debug message with the entry's fields
func (e *LoggerEntry) Debug(message string) {
	e.logger.log(LevelDebug, message, e.fields)
}

// Info logs an info message with the entry's fields
func (e *LoggerEntry) Info(message string) {
	e.logger.log(LevelInfo, message, e.fields)
}

// Warn logs a warning message with the entry's fields
func (e *LoggerEntry) Warn(message string) {
	e.logger.log(LevelWarn, message, e.fields)
}

// Error logs an error message with the entry's fields
func (e *LoggerEntry) Error(message string) {
	e.logger.log(LevelError, message, e.fields)
}

// logWriter implements io.Writer for redirecting log output
type logWriter struct {
	logger *Logger
	level  LogLevel
}

// Write implements io.Writer
func (w *logWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	// Remove trailing newline if present
	if len(message) > 0 && message[len(message)-1] == '\n' {
		message = message[:len(message)-1]
	}
	w.logger.log(w.level, message, nil)
	return len(p), nil
}

// Global logger instance
var defaultLogger *Logger
var defaultLoggerOnce sync.Once

// SetDefaultLogger sets the global default logger
func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}

// GetDefaultLogger returns the global default logger
func GetDefaultLogger() *Logger {
	return defaultLogger
}

// RedirectStdLog redirects the standard library's log output to the logger
func RedirectStdLog(logger *Logger) {
	log.SetOutput(logger.Writer())
	log.SetFlags(0) // Remove timestamp since our logger adds it
}
