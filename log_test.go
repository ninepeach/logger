package logger

import (
	"bytes"
	"os"
	"testing"
)

// Helper function to create a new standard logger for testing
func newTestStdLogger(time, debug, trace, colors, pid bool) *Logger {
	return NewStdLogger(time, debug, trace, colors, pid)
}

// Helper function to create a new file logger for testing
func newTestFileLogger(filename string, time, debug, trace, pid bool) *Logger {
	return NewFileLogger(filename, time, debug, trace, pid)
}

// Test NewStdLogger function
func TestNewStdLogger(t *testing.T) {
	l := newTestStdLogger(true, true, false, false, true)

	if l == nil {
		t.Fatalf("expected a new logger, got nil")
	}

	// Test logging at different levels
	var buf bytes.Buffer
	l.logger.SetOutput(&buf)

	l.Noticef("This is a notice log")
	if !bytes.Contains(buf.Bytes(), []byte("[INF] This is a notice log")) {
		t.Errorf("expected 'Notice' log output, got %s", buf.String())
	}

	l.Warnf("This is a warning log")
	if !bytes.Contains(buf.Bytes(), []byte("[WRN] This is a warning log")) {
		t.Errorf("expected 'Warning' log output, got %s", buf.String())
	}
}

// Test NewFileLogger function
func TestNewFileLogger(t *testing.T) {
	// Create a temporary file for testing
	tmpFile := "./test.log"
	defer os.Remove(tmpFile)

	l := newTestFileLogger(tmpFile, true, true, false, true)

	if l == nil {
		t.Fatalf("expected a new file logger, got nil")
	}

	// Test logging at different levels
	var buf bytes.Buffer
	l.logger.SetOutput(&buf)

	l.Noticef("This is a notice log")
	if !bytes.Contains(buf.Bytes(), []byte("[INF] This is a notice log")) {
		t.Errorf("expected 'Notice' log output, got %s", buf.String())
	}

	// Check if the log file is created
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Errorf("expected log file to be created, but got error: %v", err)
	}
}

// Test Logger methods with Debug and Trace options enabled
func TestLoggerWithDebugAndTrace(t *testing.T) {
	l := newTestStdLogger(true, true, true, false, true)

	var buf bytes.Buffer
	l.logger.SetOutput(&buf)

	l.Debugf("This is a debug log")
	if !bytes.Contains(buf.Bytes(), []byte("[DBG] This is a debug log")) {
		t.Errorf("expected 'Debug' log output, got %s", buf.String())
	}

	l.Tracef("This is a trace log")
	if !bytes.Contains(buf.Bytes(), []byte("[TRC] This is a trace log")) {
		t.Errorf("expected 'Trace' log output, got %s", buf.String())
	}
}

// Test Logger methods when Debug and Trace options are disabled
func TestLoggerWithoutDebugAndTrace(t *testing.T) {
	l := newTestStdLogger(true, false, false, false, true)

	var buf bytes.Buffer
	l.logger.SetOutput(&buf)

	l.Debugf("This debug log should not be printed")
	if bytes.Contains(buf.Bytes(), []byte("[DBG] This debug log should not be printed")) {
		t.Errorf("expected no 'Debug' log output, got %s", buf.String())
	}

	l.Tracef("This trace log should not be printed")
	if bytes.Contains(buf.Bytes(), []byte("[TRC] This trace log should not be printed")) {
		t.Errorf("expected no 'Trace' log output, got %s", buf.String())
	}
}

// Test Logger with file rotation and max file size limit
func TestLoggerFileRotation(t *testing.T) {
	tmpFile := "./test_rotate.log"
	defer os.Remove(tmpFile)

	l := newTestFileLogger(tmpFile, true, true, true, true)

	if l == nil {
		t.Fatalf("expected a new file logger, got nil")
	}

	// Simulate logging with file rotation
	for i := 0; i < 10; i++ {
		l.Noticef("Log message number %d", i)
	}

	// Check if the file has been written
	fileInfo, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("expected log file to be created, but got error: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Errorf("expected log file to contain logs, got empty file")
	}
}

// Test Fatal log level (should exit the program)
func TestLoggerFatal(t *testing.T) {
	tmpFile := "./test_fatal.log"
	defer os.Remove(tmpFile)

	// Capture the output of the logger
	var buf bytes.Buffer
	l := newTestFileLogger(tmpFile, true, false, false, true)
	l.logger.SetOutput(&buf)

	// Trigger a fatal error (this will cause the test to fail)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on Fatalf, but did not panic")
		}
	}()

	l.Fatalf("This is a fatal error log")

	// Ensure that the Fatal log has been written to the buffer
	if !bytes.Contains(buf.Bytes(), []byte("[FTL] This is a fatal error log")) {
		t.Errorf("expected 'Fatal' log output, got %s", buf.String())
	}
}
