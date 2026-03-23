package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"fooder-backend/core/constants"
	"github.com/labstack/echo/v4"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

type LogConfig struct {
	Level         LogLevel
	JSONFormat    bool
	DailyRotation bool
	EnableFile    bool
}

type Logger struct {
	*slog.Logger
	config     LogConfig
	currentDay string
	logFile    *os.File
	level      slog.Level
	mu         sync.Mutex
}

var defaultLogger *Logger

func NewLogger(config LogConfig) (*Logger, error) {
	var level slog.Level
	switch config.Level {
	case LogLevelDebug:
		level = slog.LevelDebug
	case LogLevelInfo:
		level = slog.LevelInfo
	case LogLevelWarn:
		level = slog.LevelWarn
	case LogLevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   "timestamp",
					Value: a.Value,
				}
			}
			return a
		},
	}

	logger := &Logger{
		config: config,
		level:  level,
	}

	var writer io.Writer
	if config.EnableFile {
		logDir := filepath.Join("logs")
		if err := os.MkdirAll(logDir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create logs directory: %w", err)
		}

		if config.DailyRotation {
			if err := logger.cleanupOldLogs(logDir); err != nil {
				fmt.Printf("Warning: failed to cleanup old logs: %v\n", err)
			}
		}

		logPath := logger.getLogFilePath(logDir)
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			fmt.Printf("Warning: cannot open log file %s, fallback to stdout: %v\n", logPath, err)
			writer = os.Stdout
		} else {
			logger.logFile = file
			logger.currentDay = time.Now().Format("2006-01-02")
			writer = file
		}
	} else {
		writer = os.Stdout
	}

	var handler slog.Handler
	if config.JSONFormat {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}

	logger.Logger = slog.New(handler)
	return logger, nil
}

func (l *Logger) getLogFilePath(logDir string) string {
	if l.config.DailyRotation {
		today := time.Now().Format("2006-01-02")
		filename := constants.LogFileName
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)
		dailyFilename := fmt.Sprintf("%s_%s%s", name, today, ext)
		return filepath.Join(logDir, dailyFilename)
	}
	return filepath.Join(logDir, constants.LogFileName)
}

func (l *Logger) cleanupOldLogs(logDir string) error {
	cutoffDate := time.Now().AddDate(0, 0, -constants.LogRetentionDays)

	entries, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("failed to read logs directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if l.isLogFileWithDate(filename) {
			fileDate, err := l.extractDateFromFilename(filename)
			if err != nil {
				continue
			}
			if fileDate.Before(cutoffDate) {
				logPath := filepath.Join(logDir, filename)
				if err := os.Remove(logPath); err != nil {
					fmt.Printf("Warning: failed to remove old log file %s: %v\n", filename, err)
				} else {
					fmt.Printf("Removed old log file: %s\n", filename)
				}
			}
		}
	}

	return nil
}

func (l *Logger) isLogFileWithDate(filename string) bool {
	baseName := strings.TrimSuffix(constants.LogFileName, filepath.Ext(constants.LogFileName))
	return strings.Contains(filename, baseName+"_") && strings.Contains(filename, "-")
}

func (l *Logger) extractDateFromFilename(filename string) (time.Time, error) {
	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return time.Time{}, fmt.Errorf("invalid filename format")
	}

	datePart := strings.TrimSuffix(parts[len(parts)-1], filepath.Ext(filename))
	return time.Parse("2006-01-02", datePart)
}

func (l *Logger) rotateLogIfNeeded() error {
	if !l.config.DailyRotation || l.logFile == nil {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if l.currentDay == today {
		return nil
	}

	if err := l.logFile.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	logDir := filepath.Join("logs")
	logPath := l.getLogFilePath(logDir)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	l.logFile = file
	l.currentDay = today

	opts := &slog.HandlerOptions{
		Level: l.level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   "timestamp",
					Value: a.Value,
				}
			}
			return a
		},
	}

	var handler slog.Handler
	if l.config.JSONFormat {
		handler = slog.NewJSONHandler(file, opts)
	} else {
		handler = slog.NewTextHandler(file, opts)
	}

	l.Logger = slog.New(handler)

	go func() {
		if err := l.cleanupOldLogs(logDir); err != nil {
			fmt.Printf("Warning: failed to cleanup old logs: %v\n", err)
		}
	}()

	return nil
}

func Init(config LogConfig) error {
	logger, err := NewLogger(config)
	if err != nil {
		return err
	}
	defaultLogger = logger
	slog.SetDefault(logger.Logger)
	return nil
}

func (l *Logger) Debug(msg string, args ...any) {
	_ = l.rotateLogIfNeeded()
	l.Logger.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	_ = l.rotateLogIfNeeded()
	l.Logger.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	_ = l.rotateLogIfNeeded()
	l.Logger.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	_ = l.rotateLogIfNeeded()
	l.Logger.Error(msg, args...)
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger:     l.Logger.With(args...),
		config:     l.config,
		currentDay: l.currentDay,
		logFile:    l.logFile,
		level:      l.level,
	}
}

func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

func Error(msg string, args ...any) {
	if defaultLogger == nil {
		fmt.Printf("ERROR: %s %v\n", msg, args)
		return
	}
	defaultLogger.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	if defaultLogger == nil {
		fmt.Printf("INFO: %s %v\n", msg, args)
		return
	}
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	if defaultLogger == nil {
		fmt.Printf("WARN: %s %v\n", msg, args)
		return
	}
	defaultLogger.Warn(msg, args...)
}

func With(args ...any) *Logger {
	if defaultLogger == nil {
		tempLogger, _ := NewLogger(LogConfig{
			Level: LogLevelInfo,
		})
		return tempLogger.With(args...)
	}
	return defaultLogger.With(args...)
}

func Close() error {
	if defaultLogger != nil {
		return defaultLogger.Close()
	}
	return nil
}

func RequestLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)

			status := c.Response().Status
			if status == 0 {
				status = 200
			}
			if err != nil && status < 400 {
				status = 500
			}

			fields := []any{
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", status,
				"latency_ms", time.Since(start).Milliseconds(),
				"remote_ip", c.RealIP(),
			}
			if rid := c.Response().Header().Get(echo.HeaderXRequestID); rid != "" {
				fields = append(fields, "request_id", rid)
			}
			if err != nil {
				fields = append(fields, "error", err.Error())
			}

			switch {
			case status >= 500:
				Error("http request completed", fields...)
			case status >= 400:
				Warn("http request completed", fields...)
			default:
				Info("http request completed", fields...)
			}
			return err
		}
	}
}
