//go:build !windows

package hotkey

import (
	"fmt"
)

// Handler is called when the hotkey is triggered.
type Handler func()

// Manager manages global hotkey registration.
type Manager struct{}

// NewManager creates a new hotkey manager.
func NewManager() *Manager {
	return &Manager{}
}

// Register registers a global hotkey with the given modifiers and key.
func (m *Manager) Register(mods interface{}, key interface{}, handler Handler) error {
	return fmt.Errorf("hotkey is only supported on Windows")
}

// RegisterFromString parses a hotkey string and registers it.
func (m *Manager) RegisterFromString(s string, handler Handler) error {
	return fmt.Errorf("hotkey is only supported on Windows")
}

// Unregister removes the current hotkey registration.
func (m *Manager) Unregister() error {
	return nil
}
