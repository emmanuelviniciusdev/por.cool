package logging

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/porcool/ingestion/internal/config"
)

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{LogLevel(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("LogLevel.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewLogger_DisabledOpenSearch(t *testing.T) {
	cfg := config.OpenSearchConfig{
		Enabled: false,
	}

	logger := NewLogger(cfg, "test-service")
	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}
	if logger.opensearch != nil {
		t.Error("OpenSearch client should be nil when disabled")
	}
	if logger.serviceName != "test-service" {
		t.Errorf("serviceName = %s, want test-service", logger.serviceName)
	}
}

func TestNewLogger_EnabledOpenSearch(t *testing.T) {
	cfg := config.OpenSearchConfig{
		Enabled:       true,
		URL:           "http://localhost:9200",
		Username:      "admin",
		Password:      "admin",
		IndexPrefix:   "test-logs",
		RetentionDays: 90,
	}

	logger := NewLogger(cfg, "test-service")
	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}
	if logger.opensearch == nil {
		t.Error("OpenSearch client should not be nil when enabled")
	}
}

func TestLogger_SetMinLevel(t *testing.T) {
	cfg := config.OpenSearchConfig{Enabled: false}
	logger := NewLogger(cfg, "test-service")

	logger.SetMinLevel(LevelWarn)
	if logger.minLevel != LevelWarn {
		t.Errorf("minLevel = %v, want %v", logger.minLevel, LevelWarn)
	}
}

func TestLogger_LogLevelFiltering(t *testing.T) {
	cfg := config.OpenSearchConfig{Enabled: false}
	logger := NewLogger(cfg, "test-service")

	// Capture stdout
	var buf bytes.Buffer
	logger.stdLogger = log.New(&buf, "", 0)

	// Set minimum level to WARN
	logger.SetMinLevel(LevelWarn)

	// These should not be logged
	logger.Debug("debug message")
	logger.Info("info message")

	// These should be logged
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Error("DEBUG message should not be logged when min level is WARN")
	}
	if strings.Contains(output, "info message") {
		t.Error("INFO message should not be logged when min level is WARN")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("WARN message should be logged")
	}
	if !strings.Contains(output, "error message") {
		t.Error("ERROR message should be logged")
	}
}

func TestLogger_LogMethods(t *testing.T) {
	cfg := config.OpenSearchConfig{Enabled: false}
	logger := NewLogger(cfg, "test-service")

	var buf bytes.Buffer
	logger.stdLogger = log.New(&buf, "", 0)
	logger.SetMinLevel(LevelDebug)

	// Test all log methods
	logger.Debug("debug test")
	logger.Info("info test")
	logger.Warn("warn test")
	logger.Error("error test")

	output := buf.String()
	if !strings.Contains(output, "[DEBUG] debug test") {
		t.Error("Debug message not formatted correctly")
	}
	if !strings.Contains(output, "[INFO] info test") {
		t.Error("Info message not formatted correctly")
	}
	if !strings.Contains(output, "[WARN] warn test") {
		t.Error("Warn message not formatted correctly")
	}
	if !strings.Contains(output, "[ERROR] error test") {
		t.Error("Error message not formatted correctly")
	}
}

func TestLogger_LogMethodsFormatted(t *testing.T) {
	cfg := config.OpenSearchConfig{Enabled: false}
	logger := NewLogger(cfg, "test-service")

	var buf bytes.Buffer
	logger.stdLogger = log.New(&buf, "", 0)
	logger.SetMinLevel(LevelDebug)

	// Test formatted log methods
	logger.Debugf("debug %s %d", "test", 123)
	logger.Infof("info %s %d", "test", 456)
	logger.Warnf("warn %s %d", "test", 789)
	logger.Errorf("error %s %d", "test", 101)

	output := buf.String()
	if !strings.Contains(output, "debug test 123") {
		t.Error("Debugf message not formatted correctly")
	}
	if !strings.Contains(output, "info test 456") {
		t.Error("Infof message not formatted correctly")
	}
	if !strings.Contains(output, "warn test 789") {
		t.Error("Warnf message not formatted correctly")
	}
	if !strings.Contains(output, "error test 101") {
		t.Error("Errorf message not formatted correctly")
	}
}

func TestLogger_WithFields(t *testing.T) {
	cfg := config.OpenSearchConfig{Enabled: false}
	logger := NewLogger(cfg, "test-service")

	var buf bytes.Buffer
	logger.stdLogger = log.New(&buf, "", 0)

	fields := map[string]interface{}{
		"user_id": "123",
		"action":  "test",
	}

	entry := logger.WithFields(fields)
	entry.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("Message not logged")
	}
	if !strings.Contains(output, "user_id") {
		t.Error("Fields not included in log")
	}
}

func TestLogger_Initialize_DisabledOpenSearch(t *testing.T) {
	cfg := config.OpenSearchConfig{Enabled: false}
	logger := NewLogger(cfg, "test-service")

	err := logger.Initialize(context.Background())
	if err != nil {
		t.Errorf("Initialize should not return error when OpenSearch is disabled: %v", err)
	}
}

func TestLogger_Writer(t *testing.T) {
	cfg := config.OpenSearchConfig{Enabled: false}
	logger := NewLogger(cfg, "test-service")

	var buf bytes.Buffer
	logger.stdLogger = log.New(&buf, "", 0)

	writer := logger.Writer()
	_, err := writer.Write([]byte("test message\n"))
	if err != nil {
		t.Errorf("Writer.Write returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("Message not written via Writer")
	}
}

func TestLogger_Close(t *testing.T) {
	cfg := config.OpenSearchConfig{Enabled: false}
	logger := NewLogger(cfg, "test-service")

	err := logger.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestSetAndGetDefaultLogger(t *testing.T) {
	cfg := config.OpenSearchConfig{Enabled: false}
	logger := NewLogger(cfg, "test-service")

	SetDefaultLogger(logger)
	got := GetDefaultLogger()

	if got != logger {
		t.Error("GetDefaultLogger did not return the set logger")
	}
}

func TestLogEntry_Structure(t *testing.T) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
		Service:   "test-service",
		Host:      "localhost",
		Fields: map[string]interface{}{
			"key": "value",
		},
	}

	if entry.Level != "INFO" {
		t.Errorf("Level = %s, want INFO", entry.Level)
	}
	if entry.Message != "test message" {
		t.Errorf("Message = %s, want test message", entry.Message)
	}
	if entry.Service != "test-service" {
		t.Errorf("Service = %s, want test-service", entry.Service)
	}
	if entry.Fields["key"] != "value" {
		t.Error("Fields not set correctly")
	}
}
