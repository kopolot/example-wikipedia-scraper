package logger

import (
	"example-wikipedia-scraper/pkg/helpers"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LevelDebug = int8(slog.LevelDebug)
	LevelInfo  = int8(slog.LevelInfo)
	LevelWarn  = int8(slog.LevelWarn)
	LevelError = int8(slog.LevelError)
	LevelFatal = int8(slog.LevelError)
)

var (
	logger *Logger
)

type Logger struct {
	logFileWriter io.Writer
	logger        *slog.Logger
	errorsCount   int64
}

func NewLogger(loggerName string, logLevel int8, cliLogging bool) *Logger {
	dir := helpers.GetAbsoluteFilePath(filepath.Join("..", "..", "var", "log"))
	filename := filepath.Join(dir, loggerName+".log")
	logFile := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	}
	writers := []io.Writer{logFile}
	if cliLogging {
		writers = append(writers, os.Stdout)
	}
	multiWriter := io.MultiWriter(writers...)
	serviceName := strings.TrimSuffix(loggerName, "_app")
	level := slog.Level(logLevel)
	handler := newLogHandler(multiWriter, level)
	return &Logger{
		logger: slog.New(handler).With(
			slog.String("service", serviceName),
		),
		logFileWriter: logFile,
	}
}

func newLogHandler(w io.Writer, level slog.Level) slog.Handler {
	opts := &slog.HandlerOptions{Level: level}
	if os.Getenv("LOG_FORMAT") == "json" {
		return slog.NewJSONHandler(w, opts)
	}
	return NewPlainTextHandler(w, level)
}

func (l *Logger) GetLogWriter() io.Writer {
	return l.logFileWriter
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	atomic.AddInt64(&l.errorsCount, 1)
	l.logger.Error(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.logger.Error(msg, args...)
	os.Exit(1)
}

func (l *Logger) Log(level int8, msg string, args ...any) {
	switch level {
	case LevelDebug:
		l.Debug(msg, args...)
	case LevelInfo:
		l.Info(msg, args...)
	case LevelWarn:
		l.Warn(msg, args...)
	case LevelError:
		l.Error(msg, args...)
	}
}

func (l *Logger) GetErrorsCount() int {
	return int(atomic.LoadInt64(&l.errorsCount))
}

func GetLogger() *Logger {
	return logger
}

func Init(fileNamePrefix string, level int8, cliLogging bool) {
	logger = NewLogger(fileNamePrefix+"_app", level, cliLogging)
	slog.SetDefault(logger.logger)
}

func Info(msg string, args ...any) {
	Log(LevelInfo, msg, args...)
}

func Error(msg string, args ...any) {
	Log(LevelError, msg, args...)
}

func Warn(msg string, args ...any) {
	Log(LevelWarn, msg, args...)
}

func Debug(msg string, args ...any) {
	Log(LevelDebug, msg, args...)
}

func Fatal(msg string, args ...any) {
	Log(LevelError, msg, args...)
	os.Exit(1)
}

func Log(level int8, msg string, args ...any) {
	if logger == nil {
		panic("Logger is not initialized. Call Init() first.")
	}
	logger.Log(level, msg, args...)
}

func GetLogWriter() io.Writer {
	if logger == nil {
		panic("Logger is not initialized. Call Init() first.")
	}
	return logger.GetLogWriter()
}
