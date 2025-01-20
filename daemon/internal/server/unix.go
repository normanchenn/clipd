package server

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/normanchenn/clipd/daemon/internal/config"
	"github.com/normanchenn/clipd/daemon/internal/logging"
	"github.com/normanchenn/clipd/daemon/internal/requests"
)

type UnixServer struct {
	logger     logging.Logger
	socketPath string
	listener   net.Listener
}

func NewUnixServer(config config.Config, logger logging.Logger) *UnixServer {
	return &UnixServer{
		logger:     logger,
		socketPath: config.SocketPath,
	}
}

func (s *UnixServer) Start(ctx context.Context) error {
	if err := os.RemoveAll(s.socketPath); err != nil {
		return fmt.Errorf("cleaning up socket file %s: %w", s.socketPath, err)
	}

	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("creating unix socket: %w", err)
	}
	s.listener = listener

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("accepting unix connection", err)
			continue
		}
		go requests.HandleConnection(conn, s.logger)
	}
}

func (s *UnixServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
