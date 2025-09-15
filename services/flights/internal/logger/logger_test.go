package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func parseLog(testHelper *testing.T, buffer *bytes.Buffer) map[string]interface{} {
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buffer.Bytes(), &logEntry); err != nil {
		testHelper.Fatalf("invalid JSON: %v", err)
	}
	return logEntry
}

func TestInit(testHelper *testing.T) {
	testCases := []struct {
		name        string
		environment string
		expectLevel slog.Level
	}{
		{"dev environment", "dev", slog.LevelDebug},
		{"prod environment", "prod", slog.LevelInfo},
		{"empty environment", "", slog.LevelInfo},
	}

	for _, testCase := range testCases {
		testHelper.Run(testCase.name, func(subTest *testing.T) {
			Logger = nil
			Init(testCase.environment)

			if Logger == nil {
				subTest.Fatal("Logger should not be nil after Init")
			}

			if !Logger.Enabled(context.Background(), testCase.expectLevel) {
				subTest.Errorf("Logger should be enabled for level %v", testCase.expectLevel)
			}
		})
	}
}

func TestPackageInit(testHelper *testing.T) {
	if Logger == nil {
		testHelper.Error("Logger should be initialized by package init")
	}
}

func TestInfo(testHelper *testing.T) {
	var outputBuffer bytes.Buffer
	Logger = slog.New(slog.NewJSONHandler(&outputBuffer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	Info("test message", "key", "value")

	logEntry := parseLog(testHelper, &outputBuffer)

	if logEntry["level"] != "INFO" {
		testHelper.Errorf("Expected level INFO, got %v", logEntry["level"])
	}
	if logEntry["msg"] != "test message" {
		testHelper.Errorf("Expected msg 'test message', got %v", logEntry["msg"])
	}
	if logEntry["key"] != "value" {
		testHelper.Errorf("Expected key 'value', got %v", logEntry["key"])
	}
}

func TestDebug(testHelper *testing.T) {
	var outputBuffer bytes.Buffer
	Logger = slog.New(slog.NewJSONHandler(&outputBuffer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	Debug("debug message", "debug_key", "debug_value")

	logEntry := parseLog(testHelper, &outputBuffer)

	if logEntry["level"] != "DEBUG" {
		testHelper.Errorf("Expected level DEBUG, got %v", logEntry["level"])
	}
	if logEntry["msg"] != "debug message" {
		testHelper.Errorf("Expected msg 'debug message', got %v", logEntry["msg"])
	}
}

func TestError(testHelper *testing.T) {
	var outputBuffer bytes.Buffer
	Logger = slog.New(slog.NewJSONHandler(&outputBuffer, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	Error("error message", "error_code", 500)

	logEntry := parseLog(testHelper, &outputBuffer)

	if logEntry["level"] != "ERROR" {
		testHelper.Errorf("Expected level ERROR, got %v", logEntry["level"])
	}
	if logEntry["msg"] != "error message" {
		testHelper.Errorf("Expected msg 'error message', got %v", logEntry["msg"])
	}
	if logEntry["error_code"].(float64) != 500 {
		testHelper.Errorf("Expected error_code 500, got %v", logEntry["error_code"])
	}
}

func TestWarn(testHelper *testing.T) {
	var outputBuffer bytes.Buffer
	Logger = slog.New(slog.NewJSONHandler(&outputBuffer, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))

	Warn("warning message", "warn_type", "deprecation")

	logEntry := parseLog(testHelper, &outputBuffer)

	if logEntry["level"] != "WARN" {
		testHelper.Errorf("Expected level WARN, got %v", logEntry["level"])
	}
	if logEntry["msg"] != "warning message" {
		testHelper.Errorf("Expected msg 'warning message', got %v", logEntry["msg"])
	}
}

func TestInfoContext(testHelper *testing.T) {
	var outputBuffer bytes.Buffer
	Logger = slog.New(slog.NewJSONHandler(&outputBuffer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	contextWithRequestID := context.WithValue(context.Background(), "request_id", "123")
	InfoContext(contextWithRequestID, "context message", "user_id", "456")

	logEntry := parseLog(testHelper, &outputBuffer)

	if logEntry["level"] != "INFO" {
		testHelper.Errorf("Expected level INFO, got %v", logEntry["level"])
	}
	if logEntry["msg"] != "context message" {
		testHelper.Errorf("Expected msg 'context message', got %v", logEntry["msg"])
	}
	if logEntry["user_id"] != "456" {
		testHelper.Errorf("Expected user_id '456', got %v", logEntry["user_id"])
	}
}

func TestErrorContext(testHelper *testing.T) {
	var outputBuffer bytes.Buffer
	Logger = slog.New(slog.NewJSONHandler(&outputBuffer, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	contextWithTraceID := context.WithValue(context.Background(), "trace_id", "abc123")
	ErrorContext(contextWithTraceID, "context error", "status", "failed")

	logEntry := parseLog(testHelper, &outputBuffer)

	if logEntry["level"] != "ERROR" {
		testHelper.Errorf("Expected level ERROR, got %v", logEntry["level"])
	}
	if logEntry["msg"] != "context error" {
		testHelper.Errorf("Expected msg 'context error', got %v", logEntry["msg"])
	}
}

func TestLogLevelFiltering(testHelper *testing.T) {
	var outputBuffer bytes.Buffer
	Logger = slog.New(slog.NewJSONHandler(&outputBuffer, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))

	Debug("debug message")
	Info("info message")

	firstOutput := outputBuffer.String()
	if strings.Contains(firstOutput, "debug message") {
		testHelper.Error("Debug message should be filtered out at WARN level")
	}
	if strings.Contains(firstOutput, "info message") {
		testHelper.Error("Info message should be filtered out at WARN level")
	}

	outputBuffer.Reset()
	Warn("warning message")
	Error("error message")

	secondOutput := outputBuffer.String()
	if !strings.Contains(secondOutput, "warning message") {
		testHelper.Error("Warning message should be logged at WARN level")
	}
	if !strings.Contains(secondOutput, "error message") {
		testHelper.Error("Error message should be logged at WARN level")
	}
}

func TestJSONOutput(testHelper *testing.T) {
	var outputBuffer bytes.Buffer
	Logger = slog.New(slog.NewJSONHandler(&outputBuffer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	Info("test", "string_val", "text", "int_val", 42, "bool_val", true)

	logEntry := parseLog(testHelper, &outputBuffer)

	expectedFields := []string{"time", "level", "msg", "string_val", "int_val", "bool_val"}
	for _, field := range expectedFields {
		if _, exists := logEntry[field]; !exists {
			testHelper.Errorf("Expected field %s not found in log output", field)
		}
	}
}

func TestGlobalLoggerSetDefault(testHelper *testing.T) {
	oldDefaultLogger := slog.Default()

	Init("test")

	if slog.Default() == oldDefaultLogger {
		testHelper.Error("slog.Default() should be updated after Init()")
	}

	if slog.Default() != Logger {
		testHelper.Error("slog.Default() should be the same as Logger after Init()")
	}
}
