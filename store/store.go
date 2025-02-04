package store

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"

	"github.com/nccapo/kv-store/internal/gen"
)

// KeyValueStore uses a thread safe RWMutex for concurrent access
type KeyValueStore struct {
	data map[string]string
	mu   sync.RWMutex // for concurrent access
}

func New() *KeyValueStore {
	return &KeyValueStore{
		data: make(map[string]string),
	}
}

func (kv *KeyValueStore) Set(key, value string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[key] = value
}

func (kv *KeyValueStore) Get(key string) (string, bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	val, ok := kv.data[key]
	return val, ok
}

func (kv *KeyValueStore) Delete(key string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	_, ok := kv.data[key]
	if !ok {
		return fmt.Errorf("key '%s' doesn't exist", key)
	}

	delete(kv.data, key)

	return nil
}

func (kv *KeyValueStore) Exists(key string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	_, ok := kv.data[key]
	if ok {
		return true
	}

	return false
}

func (kv *KeyValueStore) SaveSnapshot() error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	s := gen.RandStringRunes(10)

	temp, err := os.CreateTemp("", fmt.Sprintf("sanpshot-%s.tmp", s))
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())

	encoder := gob.NewEncoder(temp)
	if err := encoder.Encode(kv.data); err != nil {
		return err
	}

	if err := os.Rename(temp.Name(), "snapshot.gob"); err != nil {
		return err
	}

	return nil
}
