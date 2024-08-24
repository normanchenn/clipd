package main

import (
	"fmt"

	"github.com/normanchenn/clipd/daemon/internal/config"
	"github.com/normanchenn/clipd/daemon/internal/logging"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		// TODO: log
		fmt.Println("error with loading config")
		panic("")
	}

	logger, err := logging.NewLogLogger(config)
	if err != nil {
		// TODO: log
		fmt.Println("error with loading logger")
		panic("")
	}
	defer logger.Close()
	// logger.Info("testing info")
	// logger.Error("testing error", nil)
	// logger.Debug("testing debug")
}
