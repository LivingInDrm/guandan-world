package log

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestFileLogger(t *testing.T) {
	dir := t.TempDir()
	logger, err := NewFileLogger(dir, LevelDebug)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Debug("debug message", "key1", "value1")
	logger.Info("info message", "key2", 123)
	logger.Warn("warn message", "key3", true)
	logger.Error("error message", "key4", 3.14)

	date := time.Now().In(beijingLoc).Format("2006-01-02")
	logFile := filepath.Join(dir, date+".log")
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	t.Logf("Log content:\n%s", string(content))
}

func TestConcurrentWrite(t *testing.T) {
	dir := t.TempDir()
	logger, err := NewFileLogger(dir, LevelInfo)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Close()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			logger.Info("concurrent log", "goroutine_id", id)
		}(i)
	}
	wg.Wait()

	date := time.Now().In(beijingLoc).Format("2006-01-02")
	logFile := filepath.Join(dir, date+".log")
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	lines := 0
	for _, b := range content {
		if b == '\n' {
			lines++
		}
	}
	if lines != 100 {
		t.Errorf("expected 100 lines, got %d", lines)
	}
}

func TestNoOpLogger(t *testing.T) {
	logger := DefaultLogger()
	logger.Debug("should not panic")
	logger.Info("should not panic")
	logger.Warn("should not panic")
	logger.Error("should not panic")
	if err := logger.Close(); err != nil {
		t.Errorf("Close() should return nil, got %v", err)
	}
}

func TestGlobalLogger(t *testing.T) {
	dir := t.TempDir()
	if err := Init(dir, LevelDebug); err != nil {
		t.Fatalf("failed to init global logger: %v", err)
	}
	defer Close()

	Debug("global debug", "key", "value")
	Info("global info", "count", 42)
	Warn("global warn")
	Error("global error", "err", "something wrong")

	date := time.Now().In(beijingLoc).Format("2006-01-02")
	logFile := filepath.Join(dir, date+".log")
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	t.Logf("Global log content:\n%s", string(content))
}

func TestSetLogger(t *testing.T) {
	SetLogger(DefaultLogger())
	Info("should not panic with noop logger")

	SetLogger(nil)
	Info("should not panic with nil logger")
}
