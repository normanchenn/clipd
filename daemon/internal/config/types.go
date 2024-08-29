package config

type Config struct {
	// editable settings
	PollingInterval int      `toml:"polling_interval"`
	CacheSize       int      `toml:"cache_size"`
	LogLevel        LogLevel `toml:"log_level"`

	// constant settings
	LogPath    string `toml:"-"`
	SocketPath string `toml:"-"`
	StoreDir   string `toml:"-"`
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelError LogLevel = "error"
)
