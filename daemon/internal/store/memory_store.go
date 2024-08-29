package store

import (
	"fmt"
	"sync"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/normanchenn/clipd/daemon/internal/config"
	"github.com/normanchenn/clipd/daemon/internal/errors"
	"github.com/normanchenn/clipd/daemon/internal/logging"
)

type MemoryStore struct {
	mutex   sync.Mutex
	tree    redblacktree.Tree
	idMap   map[string]*redblacktree.Node
	first   *redblacktree.Node
	limit   int
	counter int
	logger  logging.Logger
}

func NewMemoryStore(config *config.Config, logger logging.Logger) *MemoryStore {
	tree := redblacktree.NewWithIntComparator()
	return &MemoryStore{
		tree:    *tree,
		idMap:   make(map[string]*redblacktree.Node),
		limit:   config.CacheSize,
		counter: 0,
		logger:  logger,
	}
}

func (m *MemoryStore) SaveEntry(entry Entry) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.tree.Size() == m.limit {
		m.tree.Remove(m.tree.Right().Key.(int))
	}

	// smaller keys will be on the left
	m.tree.Put(-m.counter, entry)
	m.counter++
	var first = m.tree.Left()
	m.idMap[entry.ID] = first
	m.first = first

	return nil
}

func (m *MemoryStore) GetEntryByIndex(index int) (Entry, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if index == 0 {
		if m.first == nil {
			return Entry{}, fmt.Errorf("request when cache is empty: %w", errors.ErrNotFound)
		}
		return m.first.Value.(Entry), nil
	}

	var key = -m.counter + index + 1
	var entry, found = m.tree.Get(key)
	if found {
		return entry.(Entry), nil
	}
	return Entry{}, errors.ErrNotFound
}

func (m *MemoryStore) GetEntryByID(id string) (Entry, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var entry, found = m.idMap[id]
	if !found {
		return Entry{}, errors.ErrNotFound
	}
	return entry.Value.(Entry), nil
}

func (m *MemoryStore) GetAllEntries() ([]Entry, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var values = m.tree.Values()
	var entries = make([]Entry, len(values))
	for idx, val := range values {
		entries[idx] = val.(Entry)
	}
	return entries, nil
}

func (m *MemoryStore) GetEntryRange(start_index int, end_index int) ([]Entry, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var start_key = -m.counter + start_index + 1
	var entry = m.tree.GetNode(start_key)
	if entry == nil {
		return []Entry{}, errors.ErrNotFound
	}

	var entries []Entry
	var it = m.tree.IteratorAt(entry)
	for i := start_index; i <= end_index; i++ {
		var entry = it.Value().(Entry)
		entries = append(entries, entry)

		if !it.Next() {
			break
		}
	}

	if len(entries) < end_index-start_index+1 {
		return entries, fmt.Errorf("end index out of bounds: %w", errors.ErrOutOfBounds)
	}
	return entries, nil
}

func (m *MemoryStore) Clear() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.tree = *redblacktree.NewWith(m.tree.Comparator)
	m.idMap = make(map[string]*redblacktree.Node)
	m.first = nil
	m.counter = 0
	return nil
}

func (m *MemoryStore) Size() (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.tree.Size(), nil
}

func (m *MemoryStore) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return nil
}
