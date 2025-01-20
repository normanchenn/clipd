package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/normanchenn/clipd/daemon/internal/config"
	"github.com/normanchenn/clipd/daemon/internal/logging"
	"github.com/normanchenn/clipd/daemon/internal/polling"
	"github.com/normanchenn/clipd/daemon/internal/server"
	"github.com/normanchenn/clipd/daemon/internal/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	logger, err := logging.NewFileLogger(config)
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("clipd logger setup")
	logger.Info("clipd configuration loaded", "configuration", fmt.Sprintf("%+v", config))

	store, err := storage.NewSQLiteStorage(config)
	if err != nil {
		logger.Fatal("couldn't create store", err)
	}
	defer store.Close()

	logger.Info("clipd store initialized")

	var wg sync.WaitGroup

	updates := make(chan string)

	poller := polling.NewPoller(config, logger)
	wg.Add(1)
	go func() {
		defer wg.Done()
		poller.Start(ctx, updates)
	}()

	logger.Info("clipd polling started")

	wg.Add(1)
	go func() {
		defer wg.Done()
		storage.Update(logger, store, updates)
	}()

	logger.Info("clipd updating store started")

	server := server.NewUnixServer(config, logger)

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Start(ctx)
	}()

	logger.Info("clipd server started")

	cancel()
	wg.Wait()

	logger.Info("clipd exiting...")
}
