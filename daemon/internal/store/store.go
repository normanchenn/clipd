package store

import (
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	ID        string
	CreatedAt time.Time
	Data      string
	Tags      []string
}

func NewEntry(data string, tags []string) Entry {
	return Entry{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Data:      data,
		Tags:      tags,
	}
}

type Store interface {
	SaveEntry(entry Entry) error
	GetEntryByIndex(index int) (Entry, error)
	GetEntryByID(id string) (Entry, error)
	GetAllEntries() ([]Entry, error)
	GetEntryRange(start_index int, end_index int) ([]Entry, error)
	Clear() error
	Size() (int, error)
	Close() error
}
