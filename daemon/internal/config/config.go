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
		// TODO: error
		return fmt.Errorf("")
	}

	if c.CacheSize <= 0 {
		// TODO: error
		return fmt.Errorf("")
	}

	if err := validateLogLevel(c.LogLevel); err != nil {
		return err
	}

	return nil
}

func (c *Config) setDefaults() {
	if c.PollingInterval == 0 {
		c.PollingInterval = 10
	}
	if c.CacheSize == 0 {
		c.CacheSize = 50
	}
	if c.LogLevel == "" {
		c.LogLevel = "debug"
	}

	// TODO: fix filepath due to mkdir not having permissions
	// c.LogPath = "/var/log/clipd/clipd.log"
	c.LogPath = "/Users/normanchen/log/clipd/clipd.log"
	c.SocketPath = "/tmp/clipd.sock"
}

func getConfigPath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	paths := []string{
		filepath.Join(user.HomeDir, ".config", "clipd", "clipd.toml"),
		filepath.Join(user.HomeDir, ".clipd.toml"),
	}
	for _, path := range paths {
		_, err = os.Stat(path)
		if err == nil {
			return path, nil
		}
	}

	return "", nil
}

func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		// TODO: log
		return nil, err
	}

	var config Config
	if configPath != "" {
		_, err = toml.DecodeFile(configPath, &config)
		if err != nil {
			// TODO: log
			return nil, err
		}
	}
	config.setDefaults()
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
		// TODO: error
		return fmt.Errorf("")
	}

}
