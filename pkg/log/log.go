package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var beijingLoc = time.FixedZone("CST", 8*60*60)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
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

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Close() error
}

type noOpLogger struct{}

func (noOpLogger) Debug(msg string, args ...any) {}
func (noOpLogger) Info(msg string, args ...any)  {}
func (noOpLogger) Warn(msg string, args ...any)  {}
func (noOpLogger) Error(msg string, args ...any) {}
func (noOpLogger) Close() error                  { return nil }

func DefaultLogger() Logger {
	return noOpLogger{}
}

type fileLogger struct {
	mu      sync.Mutex
	dir     string
	file    *os.File
	curDate string
	level   Level
}

func NewFileLogger(dir string, level Level) (Logger, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	l := &fileLogger{
		dir:   dir,
		level: level,
	}
	if err := l.rotateIfNeeded(); err != nil {
		return nil, err
	}
	return l, nil
}

func (l *fileLogger) rotateIfNeeded() error {
	now := time.Now().In(beijingLoc)
	date := now.Format("2006-01-02")
	if date == l.curDate && l.file != nil {
		return nil
	}
	if l.file != nil {
		l.file.Close()
	}
	filename := filepath.Join(l.dir, date+".log")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	l.file = f
	l.curDate = date
	return nil
}

func (l *fileLogger) log(level Level, msg string, args ...any) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.rotateIfNeeded(); err != nil {
		return
	}

	now := time.Now().In(beijingLoc)
	timestamp := now.Format("01-02T15:04:05.000")

	pkgName, typeName, methodName := extractCaller(3)

	var location string
	if typeName != "" {
		location = fmt.Sprintf("[%s.%s] [%s]", pkgName, typeName, methodName)
	} else {
		location = fmt.Sprintf("[%s] [%s]", pkgName, methodName)
	}

	var kvPairs string
	if len(args) >= 2 {
		var pairs []string
		for i := 0; i+1 < len(args); i += 2 {
			pairs = append(pairs, fmt.Sprintf("%v=%v", args[i], args[i+1]))
		}
		kvPairs = " " + strings.Join(pairs, " ")
	}

	line := fmt.Sprintf("%s %-5s %s %s%s\n", timestamp, level.String(), location, msg, kvPairs)
	l.file.WriteString(line)
}

func (l *fileLogger) Debug(msg string, args ...any) {
	l.log(LevelDebug, msg, args...)
}

func (l *fileLogger) Info(msg string, args ...any) {
	l.log(LevelInfo, msg, args...)
}

func (l *fileLogger) Warn(msg string, args ...any) {
	l.log(LevelWarn, msg, args...)
}

func (l *fileLogger) Error(msg string, args ...any) {
	l.log(LevelError, msg, args...)
}

func (l *fileLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

var globalLogger Logger = noOpLogger{}
var globalMu sync.RWMutex

func Init(dir string, level Level) error {
	logger, err := NewFileLogger(dir, level)
	if err != nil {
		return err
	}
	globalMu.Lock()
	globalLogger = logger
	globalMu.Unlock()
	return nil
}

func Close() error {
	globalMu.Lock()
	defer globalMu.Unlock()
	return globalLogger.Close()
}

func SetLogger(l Logger) {
	globalMu.Lock()
	defer globalMu.Unlock()
	if l == nil {
		globalLogger = noOpLogger{}
	} else {
		globalLogger = l
	}
}

func Debug(msg string, args ...any) {
	globalMu.RLock()
	l := globalLogger
	globalMu.RUnlock()
	if fl, ok := l.(*fileLogger); ok {
		fl.logWithSkip(LevelDebug, 3, msg, args...)
	} else {
		l.Debug(msg, args...)
	}
}

func Info(msg string, args ...any) {
	globalMu.RLock()
	l := globalLogger
	globalMu.RUnlock()
	if fl, ok := l.(*fileLogger); ok {
		fl.logWithSkip(LevelInfo, 3, msg, args...)
	} else {
		l.Info(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	globalMu.RLock()
	l := globalLogger
	globalMu.RUnlock()
	if fl, ok := l.(*fileLogger); ok {
		fl.logWithSkip(LevelWarn, 3, msg, args...)
	} else {
		l.Warn(msg, args...)
	}
}

func Error(msg string, args ...any) {
	globalMu.RLock()
	l := globalLogger
	globalMu.RUnlock()
	if fl, ok := l.(*fileLogger); ok {
		fl.logWithSkip(LevelError, 3, msg, args...)
	} else {
		l.Error(msg, args...)
	}
}

func (l *fileLogger) logWithSkip(level Level, skip int, msg string, args ...any) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.rotateIfNeeded(); err != nil {
		return
	}

	now := time.Now().In(beijingLoc)
	timestamp := now.Format("01-02T15:04:05.000")

	pkgName, typeName, methodName := extractCaller(skip)

	var location string
	if typeName != "" {
		location = fmt.Sprintf("[%s.%s] [%s]", pkgName, typeName, methodName)
	} else {
		location = fmt.Sprintf("[%s] [%s]", pkgName, methodName)
	}

	var kvPairs string
	if len(args) >= 2 {
		var pairs []string
		for i := 0; i+1 < len(args); i += 2 {
			pairs = append(pairs, fmt.Sprintf("%v=%v", args[i], args[i+1]))
		}
		kvPairs = " " + strings.Join(pairs, " ")
	}

	line := fmt.Sprintf("%s %-5s %s %s%s\n", timestamp, level.String(), location, msg, kvPairs)
	l.file.WriteString(line)
}

func extractCaller(skip int) (pkgName, typeName, methodName string) {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", "", "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown", "", "unknown"
	}

	fullName := fn.Name()

	lastSlash := strings.LastIndex(fullName, "/")
	if lastSlash >= 0 {
		fullName = fullName[lastSlash+1:]
	}

	parts := strings.Split(fullName, ".")
	if len(parts) < 2 {
		return "unknown", "", fullName
	}

	pkgName = parts[0]

	if len(parts) == 2 {
		return pkgName, "", parts[1]
	}

	typeAndMethod := parts[1:]
	if len(typeAndMethod) >= 2 {
		typeName = strings.Trim(typeAndMethod[0], "(*)")
		methodName = typeAndMethod[1]
	} else {
		methodName = typeAndMethod[0]
	}

	return pkgName, typeName, methodName
}
