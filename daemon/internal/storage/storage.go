package storage

import (
	"errors"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"github.com/normanchenn/clipd/daemon/internal/logging"
)

var (
	ErrIndexOutOfRange = errors.New("index out of range")
	ErrIdDoesntExist   = errors.New("ID doesn't exist")
	ErrIncompleteRange = errors.New("range is incomplete")
	ErrInvalidRange    = errors.New("range is invalid")
)

type Entry struct {
	ID      string
	Created time.Time
	Data    string
}

func NewEntry(data string) Entry {
	return Entry{
		ID:      shortuuid.New(),
		Created: time.Now(),
		Data:    data,
	}
}

type Storage interface {
	Save(entry Entry) (*Entry, error)
	// index 0 means the most recent value
	GetByIndex(index int) (Entry, error)
	GetByID(id string) (Entry, error)
	GetRange(start int, end int) ([]Entry, error)
	GetAll() ([]Entry, error)
	Size() (int, error)
	Close() error
}

func Update(logger logging.Logger, store Storage, updates <-chan string) {
	for update := range updates {
		entry := NewEntry(update)
		_, err := store.Save(entry)
		if err != nil {
			logger.Error("saving update to store", err)
		}
	}
}
