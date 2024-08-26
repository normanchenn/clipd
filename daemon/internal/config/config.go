package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"os/user"
	"path/filepath"
)

func (c *Config) validate() error {
	if c.PollingInterval <= 0 {
		return fmt.Errorf("invalid polling interval: %d", c.PollingInterval)
	}

	if c.CacheSize <= 0 {
		return fmt.Errorf("invalid cache size: %d", c.CacheSize)
	}

	if err := validateLogLevel(c.LogLevel); err != nil {
		return err
	}
	return nil
}

func (c *Config) setDefaults(homeDir string) {
	if c.PollingInterval == 0 {
		c.PollingInterval = 10
	}
	if c.CacheSize == 0 {
		c.CacheSize = 50
	}
	if c.LogLevel == "" {
		c.LogLevel = "debug"
	}

	c.LogPath = filepath.Join(homeDir, "log", "clipd", "clipd.log")
	c.SocketPath = filepath.Join("/", "tmp", "clipd.sock")
}

func getConfigPath(homeDir string) (string, error) {
	paths := []string{
		filepath.Join(homeDir, ".config", "clipd", "clipd.toml"),
		filepath.Join(homeDir, ".clipd.toml"),
	}
	for _, path := range paths {
		_, err := os.Stat(path)
		if err == nil {
			return path, nil
		}
	}

	return "", nil
}

func LoadConfig() (*Config, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	configPath, err := getConfigPath(user.HomeDir)
	if err != nil {
		return nil, err
	}

	var config Config
	if configPath != "" {
		_, err = toml.DecodeFile(configPath, &config)
		if err != nil {
			return nil, err
		}
	}

	config.setDefaults(user.HomeDir)
	err = config.validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func validateLogLevel(level LogLevel) error {
	switch level {
	case LogLevelDebug, LogLevelInfo, LogLevelError:
		return nil
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}
}
