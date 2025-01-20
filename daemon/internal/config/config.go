package config

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelError LogLevel = "error"
)

var (
	ErrNoHomeDirectory     = errors.New("no home directory to get config path")
	ErrInvalidPollInterval = errors.New("poll interval must be larger than 0")
	ErrInvalidCacheSize    = errors.New("cache size must be larger than 0")
	ErrInvalidLogLevel     = errors.New("log level must be one of debug, info, and error")
)

type Config struct {
	PollInterval int      `yaml:"poll_interval"`
	CacheSize    int      `yaml:"cache_size"`
	LogLevel     LogLevel `yaml:"log_level"`

	LogPath    string `yaml:"-"`
	SocketPath string `yaml:"-"`
	DbPath     string `yaml:"-"`
}

type LogLevel string

// TODO: figure out sane defaults
func NewDefaultConfig(homeDir string) Config {
	return Config{
		PollInterval: 500,
		CacheSize:    1000,
		LogLevel:     LogLevelDebug,
		LogPath:      filepath.Join(homeDir, "clipd", "clipd.log"),
		SocketPath:   filepath.Join("/", "tmp", "clipd.sock"),
		DbPath:       filepath.Join(homeDir, "clipd", "db.sqlite3"),
	}
}

func NewConfig() (Config, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		return Config{}, err
	}

	config := NewDefaultConfig(homeDir)
	configPath, _ := getConfigPath(homeDir)
	if configPath == "" {
		return config, nil
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("couldn't read config file: %w", err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf("unmarshalling config file: %w", err)
	}

	err = config.validate()
	if err != nil {
		return config, fmt.Errorf("invalid config file: %w", err)
	}
	return config, nil
}

func getHomeDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	} else if user.HomeDir == "" {
		return "", ErrNoHomeDirectory
	}
	return user.HomeDir, nil
}

func getConfigPath(homeDir string) (string, error) {
	paths := []string{
		filepath.Join(homeDir, ".config", "clipd", "clipd.yaml"),
		filepath.Join(homeDir, ".clipd.yaml"),
	}

	for _, path := range paths {
		_, err := os.Stat(path)
		if err == nil {
			return path, nil
		}
	}
	return "", nil
}

func (c Config) validate() error {
	if c.PollInterval <= 0 {
		return ErrInvalidPollInterval
	}

	if c.CacheSize <= 0 {
		return ErrInvalidCacheSize
	}

	if !isValidLogLevel(c.LogLevel) {
		return ErrInvalidLogLevel
	}
	return nil
}

func isValidLogLevel(level LogLevel) bool {
	switch level {
	case LogLevelDebug, LogLevelInfo, LogLevelError:
		return true
	default:
		return false
	}
}
