package storage

import (
	"errors"
	"sync"

	"github.com/normanchenn/clipd/daemon/internal/config"
)

type StorageManager struct {
	mutex       sync.RWMutex
	cache       Storage
	persistence Storage
	batch       []Entry
	cacheSize   int
}

func NewStorageManager(config config.Config, cache Storage, persistence Storage) StorageManager {
	return StorageManager{
		cache:       cache,
		persistence: persistence,
		batch:       make([]Entry, config.CacheSize),
		cacheSize:   config.CacheSize,
	}
}

func (s *StorageManager) flushBatch() error {
	for len(s.batch) > 0 {
		entry := s.batch[0]

		if _, err := s.persistence.Save(entry); err != nil {
			return err
		}
		s.batch = s.batch[1:]
	}
	return nil
}

func (s *StorageManager) Save(entry Entry) (*Entry, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	e, err := s.cache.Save(entry)
	if err != nil {
		return nil, err
	}

	// TODO: add write ahead logging here
	s.batch = append(s.batch, *e)
	if len(s.batch) >= s.cacheSize {
		// TODO: add logging here
		_ = s.flushBatch()
	}
	return e, nil
}

func (s *StorageManager) GetByIndex(index int) (Entry, error) {
	s.mutex.RLock()
	defer s.mutex.Unlock()

	entry, err := s.cache.GetByIndex(index)
	if err != nil && !errors.Is(err, ErrIndexOutOfRange) {
		return Entry{}, err
	} else if err == nil {
		return entry, nil
	}
	return s.persistence.GetByIndex(index)
}

func (s *StorageManager) GetByID(id string) (Entry, error) {
	s.mutex.RLock()
	defer s.mutex.Unlock()

	entry, err := s.cache.GetByID(id)
	if err != nil && !errors.Is(err, ErrIdDoesntExist) {
		return Entry{}, err
	} else if err == nil {
		return entry, err
	}
	return s.persistence.GetByID(id)
}

func (s *StorageManager) GetRange(start int, end int) ([]Entry, error) {
	s.mutex.RLock()
	defer s.mutex.Unlock()

	// TODO: fix this to work with the batching
	entries, err := s.cache.GetRange(start, end)
	if err != nil && !errors.Is(err, ErrIncompleteRange) {
		return []Entry{}, err
	} else if err == nil {
		return entries, nil
	}
	return s.persistence.GetRange(start, end)
}

func (s *StorageManager) GetAll() ([]Entry, error) {
	s.mutex.RLock()
	defer s.mutex.Unlock()

	persistentEntries, err := s.persistence.GetAll()
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, len(s.batch)+len(persistentEntries))
	copy(entries, s.batch)
	copy(entries[len(s.batch):], persistentEntries)
	return entries, nil
}

func (s *StorageManager) Size() (int, error) {
	s.mutex.RLock()
	defer s.mutex.Unlock()

	persistentSize, err := s.persistence.Size()
	if err != nil {
		return 0, err
	}
	return len(s.batch) + persistentSize, nil
}

func (s *StorageManager) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.cache.Close(); err != nil {
		return err
	}
	return s.persistence.Close()
}
