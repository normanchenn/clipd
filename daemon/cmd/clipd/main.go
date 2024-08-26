package main

import (
	"log"

	"github.com/normanchenn/clipd/daemon/internal/config"
	"github.com/normanchenn/clipd/daemon/internal/logging"
	"github.com/normanchenn/clipd/daemon/internal/store"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	logger, err := logging.NewFileLogger(config)
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	defer logger.Close()
	logger.Info("NEW INVOCATION")
	logger.Debug("file logger started successfully")

	storage := store.NewMemoryStore(config)
	defer storage.Close()
	logger.Debug("memory store started successfully")
}
