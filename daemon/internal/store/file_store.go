package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/normanchenn/clipd/daemon/internal/config"
	"github.com/normanchenn/clipd/daemon/internal/errors"
	"github.com/normanchenn/clipd/daemon/internal/logging"
)

type FileStore struct {
	mutex      sync.Mutex
	dataFile   *os.File
	offsetFile *os.File
	idFile     *os.File
	offset     int
	logger     logging.Logger
}

func NewFileStore(config *config.Config, logger logging.Logger) (*FileStore, error) {
	filepaths := map[string]string{
		"data":   filepath.Join(config.StoreDir, "data.log"),
		"offset": filepath.Join(config.StoreDir, "offset.log"),
		"id":     filepath.Join(config.StoreDir, "id.log"),
	}

	files := make(map[string]*os.File)
	for key, path := range filepaths {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		files[key] = file
	}

	store := &FileStore{
		dataFile:   files["data"],
		offsetFile: files["offset"],
		idFile:     files["id"],
		offset:     0,
		logger:     logger,
	}

	// TODO: load offset from previous files

	return store, nil
}

func (f *FileStore) serializeEntry(entry Entry) (chunk []byte, err error) {
	data, err := json.Marshal(entry)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (f *FileStore) deserializeEntry(chunk []byte) (entry Entry, err error) {
	err = json.Unmarshal(chunk, &entry)
	if err != nil {
		return Entry{}, err
	}
	return entry, nil
}

func (f *FileStore) serializeOffset(offset int) (line string) {
	// make sure last entry is on last line
	line = fmt.Sprintf("\n%d", offset)
	return line
}

func (f *FileStore) deserializeOffset(line string) (offset int, err error) {
	// TODO: check if this format string is correct
	_, err = fmt.Sscanf(line, "\n%d", &offset)
	if err != nil {
		return -1, err
	}
	return offset, nil
}

func (f *FileStore) serializeIDToOffset(id string, offset int) (line string) {
	line = fmt.Sprintf("\n%s:%d", id, offset)
	return line
}

func (f *FileStore) deserializeIDToOffset(line string) (id string, offset int, err error) {
	// TODO: check if this format string is correct
	_, err = fmt.Sscanf(line, "\n%s:%d", &id, &offset)
	if err != nil {
		return "", 0, err
	}
	return id, offset, nil

}

func (f *FileStore) seekPreviousLine(currentPosition int64, file *os.File) (int64, error) {
	return 0, nil
}

func (f *FileStore) readLine(currentPosition int64, file *os.File) (string, error) {
	return "", nil
}

func (f *FileStore) getOffset(currentPosition int64, file *os.File) (int, error) {
	line, err := f.readLine(currentPosition, f.offsetFile)
	if err != nil {

	}
	offset, err := f.deserializeOffset(line)
	if err != nil {

	}
	return offset, nil
}

func (f *FileStore) findOffsetsByIndex(index int, file *os.File) (start_offset int, end_offset int, err error) {
	fileInfo, err := file.Stat()
	if err != nil {

	}
	currentPosition := fileInfo.Size()

	for i := 0; i < index; i++ {
		currentPosition, err = f.seekPreviousLine(currentPosition, file)
		if err != nil {

		}
	}

	if currentPosition == fileInfo.Size() {
		end_offset = int(currentPosition)
	} else {
		end_offset, err = f.getOffset(currentPosition, file)
		if err != nil {

		}
	}

	currentPosition, err = f.seekPreviousLine(currentPosition, file)
	if err != nil {

	}

	start_offset, err = f.getOffset(currentPosition, file)
	if err != nil {

	}
	return start_offset, end_offset, nil
}

func (f *FileStore) getEntryByOffsets(start_offset int, end_offset int, file *os.File) (Entry, error) {
	_, err := file.Seek(int64(start_offset), io.SeekStart)
	if err != nil {

	}

	size := end_offset - start_offset
	data := make([]byte, size)
	_, err = file.Read(data)
	if err != nil && err != io.EOF {

	}

	entry, err := f.deserializeEntry(data)
	if err != nil {

	}
	return entry, nil
}

func (f *FileStore) SaveEntry(entry Entry) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	data, err := f.serializeEntry(entry)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrSerializationFailed, err)
	}

	dataWriter := bufio.NewWriter(f.dataFile)
	defer dataWriter.Flush()
	n, err := dataWriter.Write(data)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrWriteFailed, err)
	}

	offsetWriter := bufio.NewWriter(f.offsetFile)
	defer offsetWriter.Flush()
	_, err = offsetWriter.WriteString(f.serializeOffset(f.offset))
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrWriteFailed, err)
	}

	idWriter := bufio.NewWriter(f.idFile)
	defer idWriter.Flush()
	_, err = idWriter.WriteString(f.serializeIDToOffset(entry.ID, f.offset))
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrWriteFailed, err)
	}

	f.offset += n

	return nil
}

func (f *FileStore) GetEntryByIndex(index int) (Entry, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	start_offset, end_offset, err := f.findOffsetsByIndex(index, f.offsetFile)
	if err != nil {
		return Entry{}, err
	}

	entry, err := f.getEntryByOffsets(start_offset, end_offset, f.dataFile)
	if err != nil {
		return Entry{}, err
	}
	return entry, nil
}

func (f *FileStore) GetEntryByID(id string) (Entry, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return Entry{}, nil
}

func (f *FileStore) GetAllEntries() ([]Entry, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return []Entry{}, nil
}

func (f *FileStore) GetEntryRange(start_index int, end_index int) ([]Entry, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return []Entry{}, nil
}

func (f *FileStore) Clear() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return nil
}

func (f *FileStore) Size() (int, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return 0, nil
}

func (f *FileStore) Close() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return nil
}
