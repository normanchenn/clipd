package store_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/normanchenn/clipd/daemon/internal/config"
	"github.com/normanchenn/clipd/daemon/internal/store"
)

func saveEntry(t *testing.T, storage *store.MemoryStore, entry store.Entry) {
	t.Helper()

	err := storage.SaveEntry(entry)
	if err != nil {
		t.Fatalf("failed to save entry: %v", err)
	}
	t.Logf("saved entry %v", entry)
}

func assertSize(t *testing.T, entries []store.Entry, expected int) {
	t.Helper()

	if len(entries) != expected {
		t.Fatalf("expected size %d, got %d", expected, len(entries))
	}
	t.Logf("expected size matches entry length: %d", expected)
}

func assertEntry(t *testing.T, entry store.Entry, expected string) {
	t.Helper()

	if entry.Data != expected {
		t.Fatalf("expected data %s, got %s", expected, entry.Data)
	}
	t.Logf("expected matches data: %s", expected)
}

func assertEntries(t *testing.T, entries []store.Entry, expected []string) {
	t.Helper()

	if len(entries) != len(expected) {
		t.Fatalf("expected %d entries, received %d", len(expected), len(entries))
	}

	for i, entry := range entries {
		if entry.Data != expected[i] {
			t.Fatalf("expected data %s, got %s", expected[i], entry.Data)
		}
	}
	t.Logf("expected data matches entries of size %d", len(entries))
}

func assertError(t *testing.T, err error, expectError bool) {
	t.Helper()

	if expectError {
		if err == nil {
			t.Fatal("expected an error but didn't get one")
		} else {
			t.Logf("expected an error and got an error: %v", err)
		}
	}
	if !expectError {
		if err != nil {
			t.Fatalf("did not expect an error but got an error: %v", err)
		} else {
			t.Log("did not expect an error and did not get an error")
		}
	}
}

func TestMemoryStore_GetEntryByIndex_Basic(t *testing.T) {
	config := &config.Config{CacheSize: 5}
	storage := store.NewMemoryStore(config)
	defer storage.Close()

	entries := []store.Entry{
		store.NewEntry("third", []string{}),
		store.NewEntry("second", []string{}),
		store.NewEntry("first", []string{}),
	}

	tests := []struct {
		index       int
		expected    string
		expectError bool
	}{
		{0, "first", false},
		{1, "second", false},
		{2, "third", false},
		{3, "", true},
	}

	for _, entry := range entries {
		saveEntry(t, storage, entry)
	}

	for _, test := range tests {
		t.Logf("getting entry by index for index %d, expecting %s", test.index, test.expected)
		entry, err := storage.GetEntryByIndex(test.index)
		assertError(t, err, test.expectError)
		if !test.expectError {
			assertEntry(t, entry, test.expected)
		}
	}
}

func TestMemoryStore_GetEntryByIndex_CacheLimit(t *testing.T) {
	config := &config.Config{CacheSize: 2}
	storage := store.NewMemoryStore(config)
	defer storage.Close()

	entries := []store.Entry{
		store.NewEntry("third", []string{}),
		store.NewEntry("second", []string{}),
		store.NewEntry("first", []string{}),
	}

	tests := []struct {
		index       int
		expected    string
		expectError bool
	}{
		{0, "first", false},
		{1, "second", false},
		{2, "", true},
	}

	for _, entry := range entries {
		saveEntry(t, storage, entry)
	}

	for _, test := range tests {
		t.Logf("getting entry by index for index %d, expecting %s", test.index, test.expected)
		entry, err := storage.GetEntryByIndex(test.index)
		assertError(t, err, test.expectError)
		if !test.expectError {
			assertEntry(t, entry, test.expected)
		}
	}
}

func TestMemoryStore_GetEntryByIndex_EmptyCache(t *testing.T) {
	config := &config.Config{CacheSize: 2}
	storage := store.NewMemoryStore(config)
	defer storage.Close()

	tests := []struct {
		index       int
		expected    string
		expectError bool
	}{
		{0, "first", true},
		{1, "first", true},
	}

	for _, test := range tests {
		t.Logf("getting entry by index for index %d, expecting %s", test.index, test.expected)
		entry, err := storage.GetEntryByIndex(test.index)
		assertError(t, err, test.expectError)
		if !test.expectError {
			assertEntry(t, entry, test.expected)
		}
	}
}

func TestMemoryStore_GetEntryByID_Basic(t *testing.T) {
	config := &config.Config{CacheSize: 2}
	storage := store.NewMemoryStore(config)
	defer storage.Close()

	entry := store.NewEntry("first", []string{})

	tests := []struct {
		id          string
		expected    string
		expectError bool
	}{
		{entry.ID, "first", false},
		{uuid.New().String(), "", true},
	}

	saveEntry(t, storage, entry)

	for _, test := range tests {
		t.Logf("getting entry by id for id %s, expecting %s", test.id, test.expected)
		entry, err := storage.GetEntryByID(test.id)
		assertError(t, err, test.expectError)
		if !test.expectError {
			assertEntry(t, entry, test.expected)
		}
	}
}

func TestMemoryStore_GetAllEntries_Basic(t *testing.T) {
	config := &config.Config{CacheSize: 5}
	storage := store.NewMemoryStore(config)
	defer storage.Close()

	entries := []store.Entry{
		store.NewEntry("third", []string{}),
		store.NewEntry("second", []string{}),
		store.NewEntry("first", []string{}),
	}

	expected := []string{"first", "second", "third"}

	for _, entry := range entries {
		saveEntry(t, storage, entry)
	}

	allEntries, err := storage.GetAllEntries()
	if err != nil {
		t.Fatalf("failed to get all entries: %d", err)
	}

	assertEntries(t, allEntries, expected)
}

func TestMemoryStore_GetAllEntries_CacheLimit(t *testing.T) {
	config := &config.Config{CacheSize: 2}
	storage := store.NewMemoryStore(config)
	defer storage.Close()

	entries := []store.Entry{
		store.NewEntry("third", []string{}),
		store.NewEntry("second", []string{}),
		store.NewEntry("first", []string{}),
	}

	expected := []string{"first", "second"}

	for _, entry := range entries {
		saveEntry(t, storage, entry)
	}

	allEntries, err := storage.GetAllEntries()
	if err != nil {
		t.Fatalf("failed to get all entries: %v", err)
	}

	assertEntries(t, allEntries, expected)
}

func TestMemoryStore_GetEntryRange_Basic(t *testing.T) {
	config := &config.Config{CacheSize: 5}
	storage := store.NewMemoryStore(config)
	defer storage.Close()

	entries := []store.Entry{
		store.NewEntry("third", []string{}),
		store.NewEntry("second", []string{}),
		store.NewEntry("first", []string{}),
	}

	tests := []struct {
		start_idx   int
		end_idx     int
		expected    []string
		expectError bool
	}{
		{0, 0, []string{"first"}, false},
		{0, 1, []string{"first", "second"}, false},
		{0, 2, []string{"first", "second", "third"}, false},
		{1, 1, []string{"second"}, false},
		{1, 2, []string{"second", "third"}, false},
		{2, 2, []string{"third"}, false},
		{3, 3, []string{}, true},
		{1, 3, []string{}, true},
		{-1, 1, []string{}, true},
	}

	for _, entry := range entries {
		saveEntry(t, storage, entry)
	}

	for _, test := range tests {
		entries, err := storage.GetEntryRange(test.start_idx, test.end_idx)
		assertError(t, err, test.expectError)
		if !test.expectError {
			assertEntries(t, entries, test.expected)
		}
	}
}

func TestMemoryStore_Clear_Basic(t *testing.T) {
	config := &config.Config{CacheSize: 10}
	storage := store.NewMemoryStore(config)
	defer storage.Close()

	entries := []store.Entry{
		store.NewEntry("third", []string{}),
		store.NewEntry("second", []string{}),
		store.NewEntry("first", []string{}),
	}

	for _, entry := range entries {
		saveEntry(t, storage, entry)
	}

	storage.Clear()

	entries2 := []store.Entry{
		store.NewEntry("six", []string{}),
		store.NewEntry("five", []string{}),
		store.NewEntry("four", []string{}),
	}

	expected := []string{"four", "five", "six"}

	for _, entry := range entries2 {
		saveEntry(t, storage, entry)
	}

	allEntries, err := storage.GetAllEntries()
	if err != nil {
		t.Fatalf("failed to get all entries: %d", err)
	}

	assertEntries(t, allEntries, expected)
}

func TestMemoryStore_Size_Basic(t *testing.T) {
	config := &config.Config{CacheSize: 10}
	storage := store.NewMemoryStore(config)
	defer storage.Close()

	entries := []store.Entry{
		store.NewEntry("third", []string{}),
		store.NewEntry("second", []string{}),
		store.NewEntry("first", []string{}),
	}

	expected := []int{1, 2, 3}

	allEntries, err := storage.GetAllEntries()
	if err != nil {
		t.Fatalf("failed to get all entries: %d", err)
	}
	assertSize(t, allEntries, 0)

	for i, entry := range entries {
		saveEntry(t, storage, entry)
		allEntries, err = storage.GetAllEntries()
		if err != nil {
			t.Fatalf("failed to get all entries: %d", err)
		}
		assertSize(t, allEntries, expected[i])
	}

}
