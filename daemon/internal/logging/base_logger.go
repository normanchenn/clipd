package logging

import (
	"io"
	"log/slog"
	"os"

	"github.com/normanchenn/clipd/daemon/internal/config"
)

type baseLogger struct {
	logger *slog.Logger
	writer io.Writer
}

func (l baseLogger) Info(msg string, kvp ...interface{}) {
	l.logger.Info(msg, kvp...)
}

func (l baseLogger) Error(msg string, err error, kvp ...interface{}) {
	l.logger.Error(msg, append(kvp, "error", err)...)
}

func (l baseLogger) Debug(msg string, kvp ...interface{}) {
	l.logger.Debug(msg, kvp...)
}

func (l baseLogger) Fatal(msg string, err error, kvp ...interface{}) {
	kvp = append(kvp, "fatal", true)
	l.Error(msg, err, kvp...)
	os.Exit(1)
}

func (l *baseLogger) Close() error {
	if closer, ok := l.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func newBaseLogger(writer io.Writer, config config.Config) *baseLogger {
	options := &slog.HandlerOptions{
		Level: parseLogLevel(config.LogLevel),
	}
	logger := slog.New(slog.NewTextHandler(writer, options))
	return &baseLogger{
		logger: logger,
		writer: writer,
	}
}

func parseLogLevel(level config.LogLevel) slog.Level {
	switch level {
	case config.LogLevelDebug:
		return slog.LevelDebug
	case config.LogLevelInfo:
		return slog.LevelInfo
	case config.LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
