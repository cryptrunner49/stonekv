package stone

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
)

// Store represents the StoneKVR key/value store with on-disk persistence.
type Store struct {
	file  *os.File          // File handle for the database
	index map[string]uint64 // In-memory index mapping keys to value offsets
	mu    sync.RWMutex      // Mutex for concurrent access
}

// NewStore initializes or opens a StoneKVR store at the given file path.
// It creates the file if it doesn't exist and builds the in-memory index.
func NewStore(path string) (*Store, error) {
	// Open file with read/write, create if not exists, and append mode
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	// Initialize the store
	store := &Store{
		file:  file,
		index: make(map[string]uint64),
	}

	// Build the index from the file
	err = store.buildIndex()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to build index: %v", err)
	}

	return store, nil
}

// buildIndex reads the file and constructs the in-memory index.
// It processes set and delete records to determine the latest state of each key.
func (s *Store) buildIndex() error {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	for {
		// Get the starting offset of the current record
		startOffset, err := s.file.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}

		// Read record type (0 for set, 1 for delete)
		var typeByte byte
		err = binary.Read(s.file, binary.LittleEndian, &typeByte)
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			return err
		}

		// Read key length
		var keyLen uint32
		err = binary.Read(s.file, binary.LittleEndian, &keyLen)
		if err != nil {
			return err
		}

		// Read key bytes
		keyBytes := make([]byte, keyLen)
		_, err = s.file.Read(keyBytes)
		if err != nil {
			return err
		}
		keyStr := string(keyBytes)

		if typeByte == 0 { // Set record
			// Calculate offset where value length starts
			valLenOffset := uint64(startOffset) + 1 + 4 + uint64(keyLen)
			s.index[keyStr] = valLenOffset

			// Read and skip the value
			var valLen uint32
			err = binary.Read(s.file, binary.LittleEndian, &valLen)
			if err != nil {
				return err
			}
			_, err = s.file.Seek(int64(valLen), io.SeekCurrent)
			if err != nil {
				return err
			}
		} else if typeByte == 1 { // Delete record
			delete(s.index, keyStr)
		} else {
			return fmt.Errorf("invalid record type: %d", typeByte)
		}
	}
	return nil
}

// Set stores a key/value pair in the database.
// It appends the record to the file and updates the index.
func (s *Store) Set(key, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Construct the set record: [0][key_len][key][val_len][val]
	record := make([]byte, 1+4+len(key)+4+len(value))
	record[0] = 0 // Type: set
	binary.LittleEndian.PutUint32(record[1:5], uint32(len(key)))
	copy(record[5:5+len(key)], key)
	binary.LittleEndian.PutUint32(record[5+len(key):9+len(key)], uint32(len(value)))
	copy(record[9+len(key):], value)

	// Write the record to the file
	_, err := s.file.Write(record)
	if err != nil {
		return fmt.Errorf("failed to write record: %v", err)
	}

	// Calculate the offset where val_len starts
	stat, err := s.file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stat: %v", err)
	}
	startOffset := stat.Size() - int64(len(record))
	valLenOffset := uint64(startOffset) + 1 + 4 + uint64(len(key))

	// Update the index
	s.index[string(key)] = valLenOffset
	return nil
}

// Get retrieves the value associated with a key.
// Returns an error if the key does not exist.
func (s *Store) Get(key []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Look up the value offset in the index
	offset, ok := s.index[string(key)]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	// Seek to the value length position
	_, err := s.file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to seek: %v", err)
	}

	// Read value length
	var valLen uint32
	err = binary.Read(s.file, binary.LittleEndian, &valLen)
	if err != nil {
		return nil, fmt.Errorf("failed to read value length: %v", err)
	}

	// Read the value
	value := make([]byte, valLen)
	_, err = s.file.Read(value)
	if err != nil {
		return nil, fmt.Errorf("failed to read value: %v", err)
	}

	return value, nil
}

// Delete removes a key from the database.
// It appends a delete record and removes the key from the index.
func (s *Store) Delete(key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Construct the delete record: [1][key_len][key]
	record := make([]byte, 1+4+len(key))
	record[0] = 1 // Type: delete
	binary.LittleEndian.PutUint32(record[1:5], uint32(len(key)))
	copy(record[5:], key)

	// Write the record to the file
	_, err := s.file.Write(record)
	if err != nil {
		return fmt.Errorf("failed to write delete record: %v", err)
	}

	// Remove the key from the index
	delete(s.index, string(key))
	return nil
}

// Close closes the store and releases resources.
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %v", err)
	}
	return nil
}