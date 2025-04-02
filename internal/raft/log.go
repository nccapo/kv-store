package raft

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// Log represents the Raft log
type Log struct {
	mu sync.RWMutex

	entries []LogEntry
	// The index of the first entry in the log
	startIndex int64
	// The index of the last entry in the log
	lastIndex int64

	// Path to the log file
	logPath string
}

// NewLog creates a new Raft log
func NewLog(logPath string) (*Log, error) {
	l := &Log{
		entries:    make([]LogEntry, 0),
		startIndex: 0,
		lastIndex:  0,
		logPath:    logPath,
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, err
	}

	// Load existing log entries if any
	if err := l.load(); err != nil {
		return nil, err
	}

	return l, nil
}

// Append adds new entries to the log
func (l *Log) Append(entries ...LogEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries = append(l.entries, entries...)
	l.lastIndex = l.startIndex + int64(len(l.entries)) - 1

	return l.persist()
}

// Get returns the entry at the given index
func (l *Log) Get(index int64) (LogEntry, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if index < l.startIndex || index > l.lastIndex {
		return LogEntry{}, ErrLogIndexOutOfRange
	}

	return l.entries[index-l.startIndex], nil
}

// GetRange returns entries in the given range
func (l *Log) GetRange(start, end int64) ([]LogEntry, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if start < l.startIndex || end > l.lastIndex {
		return nil, ErrLogIndexOutOfRange
	}

	return l.entries[start-l.startIndex : end-l.startIndex+1], nil
}

// Truncate removes all entries after the given index
func (l *Log) Truncate(index int64) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if index < l.startIndex || index > l.lastIndex {
		return ErrLogIndexOutOfRange
	}

	l.entries = l.entries[:index-l.startIndex+1]
	l.lastIndex = index

	return l.persist()
}

// GetLastIndex returns the index of the last entry
func (l *Log) GetLastIndex() int64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.lastIndex
}

// GetLastTerm returns the term of the last entry
func (l *Log) GetLastTerm() int64 {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(l.entries) == 0 {
		return 0
	}

	return l.entries[len(l.entries)-1].Term
}

// load loads the log entries from disk
func (l *Log) load() error {
	data, err := os.ReadFile(l.logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &l.entries)
}

// persist saves the log entries to disk
func (l *Log) persist() error {
	data, err := json.Marshal(l.entries)
	if err != nil {
		return err
	}

	return os.WriteFile(l.logPath, data, 0644)
}

// Errors
var (
	ErrLogIndexOutOfRange = errors.New("log index out of range")
)
