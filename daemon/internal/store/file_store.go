package store

import (
	"github.com/normanchenn/clipd/daemon/internal/config"
)

type FileStore struct {
}

func NewFileStore(config *config.Config) *FileStore {
	return &FileStore{}
}

func (f *FileStore) SaveEntry(entry Entry) error {
	return nil
}

func (f *FileStore) GetEntryByIndex(index int) (Entry, error) {
	return Entry{}, nil
}

func (f *FileStore) GetEntryByID(id string) (Entry, error) {
	return Entry{}, nil
}

func (f *FileStore) GetAllEntries() ([]Entry, error) {
	return []Entry{}, nil
}

func (f *FileStore) GetEntryRange(start_index int, end_index int) ([]Entry, error) {
	return []Entry{}, nil
}

func (f *FileStore) Clear() error {
	return nil
}

func (f *FileStore) Size() (int, error) {
	return 0, nil
}

func (f *FileStore) Close() error {
	return nil
}
