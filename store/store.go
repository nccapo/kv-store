// Package store provides simple KV Storage functionalities
// package uses RWMutex for concurrent access
package store

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// item is data representation
type item struct {
	value    string
	expireAt time.Time
	cancel   context.CancelFunc
}

// Store uses a thread safe RWMutex for concurrent access
type Store struct {
	data sync.Map
}

// Set method received (key, value) parameters and set it in-memory using RWMutex for concurrent access
func (kv *Store) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	kCtx, cancel := context.WithTimeout(ctx, expiration)
	kv.data.Store(key, &item{
		value:    value,
		expireAt: time.Now().Add(expiration),
		cancel:   cancel,
	})

	go func() {
		<-kCtx.Done()
		err := kv.Delete(ctx, key)
		if err != nil {
			return
		}
	}()

	return nil
}

// Get method received key string and search data using key, it returns data and bool value
func (kv *Store) Get(ctx context.Context, key string) (string, error) {
	val, ok := kv.data.Load(key)
	if !ok {
		return "", fmt.Errorf("key does not exist")
	}

	it := val.(*item)

	if time.Now().After(it.expireAt) {
		err := kv.Delete(ctx, key)
		if err != nil {
			return "", err
		}

		return "", fmt.Errorf("key expired")
	}

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("context canceled")
	default:
		return it.value, nil
	}
}

// Delete method received key string and returns error if data with provided key doesn't exist
func (kv *Store) Delete(ctx context.Context, key string) error {
	if val, ok := kv.data.LoadAndDelete(key); ok {
		it := val.(*item)
		it.cancel() // Clean up the context
	}
	return nil
}

// Exists method received key string and search key in map, it returns bool value depending on existed data
func (kv *Store) Exists(key string) bool {
	_, ok := kv.data.Load(key)
	if ok {
		return true
	}

	return false
}
