package store

import (
	"github.com/normanchenn/clipd/daemon/internal/config"
)

type CompositeStore struct {
}

func NewCompositeStore(config *config.Config) *CompositeStore {
	return &CompositeStore{}
}

func (c *CompositeStore) SaveEntry(entry Entry) error {
	return nil
}

func (c *CompositeStore) GetEntryByIndex(index int) (Entry, error) {
	return Entry{}, nil
}

func (c *CompositeStore) GetEntryByID(id string) (Entry, error) {
	return Entry{}, nil
}

func (c *CompositeStore) GetAllEntries() ([]Entry, error) {
	return []Entry{}, nil
}

func (c *CompositeStore) GetEntryRange(start_index int, end_index int) ([]Entry, error) {
	return []Entry{}, nil
}

func (c *CompositeStore) Clear() error {
	return nil
}

func (c *CompositeStore) Size() (int, error) {
	return 0, nil
}

func (c *CompositeStore) Close() error {
	return nil
}
