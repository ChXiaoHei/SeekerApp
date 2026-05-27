package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry represents a single search history entry.
type Entry struct {
	Query     string    `json:"query"`
	Timestamp time.Time `json:"timestamp"`
}

// Manager handles search history storage and retrieval.
type Manager struct {
	mu         sync.RWMutex
	history    []Entry
	maxEntries int
	filePath   string
}

// NewManager creates a new history manager.
func NewManager(maxEntries int) (*Manager, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(configDir, "GoEverythingToolbar")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	m := &Manager{
		history:    make([]Entry, 0),
		maxEntries: maxEntries,
		filePath:   filepath.Join(appDir, "history.json"),
	}

	if err := m.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return m, nil
}

// Add adds a new search query to the history.
func (m *Manager) Add(query string) {
	if query == "" {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove existing entry with same query (for reordering)
	for i, e := range m.history {
		if e.Query == query {
			m.history = append(m.history[:i], m.history[i+1:]...)
			break
		}
	}

	// Add new entry at the beginning
	m.history = append([]Entry{{Query: query, Timestamp: time.Now()}}, m.history...)

	// Trim to max size
	if len(m.history) > m.maxEntries {
		m.history = m.history[:m.maxEntries]
	}

	// Persist
	_ = m.saveUnlocked()
}

// GetAll returns all history entries (most recent first).
func (m *Manager) GetAll() []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Entry, len(m.history))
	copy(result, m.history)
	return result
}

// GetQueries returns just the query strings (most recent first).
func (m *Manager) GetQueries() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	queries := make([]string, len(m.history))
	for i, e := range m.history {
		queries[i] = e.Query
	}
	return queries
}

// Clear removes all history entries.
func (m *Manager) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.history = make([]Entry, 0)
	return m.saveUnlocked()
}

// Load reads history from disk.
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &m.history)
}

// saveUnlocked writes history to disk (must be called with lock held).
func (m *Manager) saveUnlocked() error {
	data, err := json.MarshalIndent(m.history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.filePath, data, 0644)
}

// SetMaxEntries updates the maximum number of entries.
func (m *Manager) SetMaxEntries(max int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.maxEntries = max

	// Trim if needed
	if len(m.history) > m.maxEntries {
		m.history = m.history[:m.maxEntries]
	}
}
