package requests

import (
	"net"

	"github.com/normanchenn/clipd/daemon/internal/logging"
)

func HandleConnection(conn net.Conn, logger logging.Logger) {
	defer conn.Close()
	// TODO: finish
	logger.Debug("handling connection...")
}
