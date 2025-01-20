package polling

import (
	"context"
	"time"

	"github.com/normanchenn/clipd/daemon/internal/clipboard"
	"github.com/normanchenn/clipd/daemon/internal/config"
	"github.com/normanchenn/clipd/daemon/internal/logging"
)

type Poller struct {
	logger       logging.Logger
	interval     int
	currentValue string
}

func NewPoller(config config.Config, logger logging.Logger) Poller {
	return Poller{
		logger:       logger,
		interval:     config.PollInterval,
		currentValue: "",
	}
}

func (p *Poller) Start(ctx context.Context, updates chan<- string) {
	// p.interval milliseconds
	duration := time.Duration(p.interval) * time.Millisecond
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	p.logger.Info("hello world")

	for {
		select {
		case <-ticker.C:
			value, err := clipboard.Get()
			if err != nil {
				p.logger.Error("getting clipboard value while polling", err)
				continue
			} else if value == p.currentValue {
				continue
			}
			p.logger.Debug("new clipboard", p.currentValue, value)
			p.currentValue = value
			updates <- value
		case <-ctx.Done():
			close(updates)
			return
		}
	}
}
