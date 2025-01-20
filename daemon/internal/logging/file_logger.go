package logging

import (
	"os"
	"path/filepath"

	"github.com/normanchenn/clipd/daemon/internal/config"
)

type FileLogger struct {
	*baseLogger
}

func NewFileLogger(config config.Config) (FileLogger, error) {
	// 0755 : rwxr-xr-x
	err := os.MkdirAll(filepath.Dir(config.LogPath), 0755)
	if err != nil {
		return FileLogger{}, err
	}

	// 0644 : rw-r--r--
	logFile, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return FileLogger{}, err
	}

	return FileLogger{
		newBaseLogger(logFile, config),
	}, nil
}
