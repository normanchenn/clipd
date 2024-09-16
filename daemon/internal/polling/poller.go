package polling

import (
	"github.com/normanchenn/clipd/daemon/internal/clipboard"
	"github.com/normanchenn/clipd/daemon/internal/logging"
	"github.com/normanchenn/clipd/daemon/internal/store"
)

type Poller struct {
	interval  int
	logger    logging.Logger
	clipboard clipboard.Clipboard
	store     store.Store
	prev      string
}
