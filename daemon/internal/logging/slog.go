package logging

import (
	"github.com/normanchenn/clipd/daemon/internal/config"
	"log/slog"
	"os"
	"path/filepath"
)

type LogLogger struct {
	logger  *slog.Logger
	logFile *os.File
}

func NewLogLogger(config *config.Config) (*LogLogger, error) {
	err := os.MkdirAll(filepath.Dir(config.LogPath), 0755)
	if err != nil {
		return nil, err
	}
	logFile, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	options := &slog.HandlerOptions{
		Level: parseLogLevel(config.LogLevel),
	}
	logger := slog.New(slog.NewTextHandler(logFile, options))

	return &LogLogger{
		logger:  logger,
		logFile: logFile,
	}, nil
}

func (l *LogLogger) Info(msg string, kvp ...interface{}) {
	l.logger.Info(msg, kvp...)
}

func (l *LogLogger) Error(msg string, err error, kvp ...interface{}) {
	l.logger.Error(msg, append(kvp, "error", err)...)
}

func (l *LogLogger) Debug(msg string, kvp ...interface{}) {
	l.logger.Debug(msg, kvp...)
}

func (l *LogLogger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
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
