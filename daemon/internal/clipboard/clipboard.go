package clipboard

import (
	"bytes"
	"os/exec"

	"github.com/normanchenn/clipd/daemon/internal/logging"
)

type Clipboard struct {
	logger logging.Logger
}

func NewClipboard(logger logging.Logger) *Clipboard {
	return &Clipboard{
		logger: logger,
	}
}

func (c *Clipboard) GetCurrentClipboard() (string, error) {
	// macos only for now
	cmd := exec.Command("pbpaste")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
