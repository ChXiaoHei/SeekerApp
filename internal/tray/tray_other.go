//go:build !windows

package tray

import (
	"fmt"
)

// MenuItem represents a menu item in the system tray.
type MenuItem struct {
	Text    string
	Handler func()
}

// Callbacks are functions called from tray menu actions.
type Callbacks struct {
	OnShowWindow func()
	OnSettings   func()
	OnExit       func()
}

// Manager manages the system tray icon and menu.
type Manager struct {
	callbacks Callbacks
}

// NewManager creates a new system tray manager.
func NewManager(callbacks Callbacks) (*Manager, error) {
	return &Manager{callbacks: callbacks}, nil
}

// SetIcon sets the tray icon from an icon file path.
func (m *Manager) SetIcon(iconPath string) error {
	return fmt.Errorf("system tray is only supported on Windows")
}

// SetIconFromResource sets the tray icon from embedded resources.
func (m *Manager) SetIconFromResource(resourceID int) error {
	return fmt.Errorf("system tray is only supported on Windows")
}

// Show displays the tray icon.
func (m *Manager) Show() error {
	return nil
}

// Hide hides the tray icon.
func (m *Manager) Hide() error {
	return nil
}

// SetTooltip sets the tray icon tooltip text.
func (m *Manager) SetTooltip(text string) error {
	return nil
}

// SetupMenu creates the context menu for the tray icon.
func (m *Manager) SetupMenu() error {
	return nil
}

// Dispose releases tray resources.
func (m *Manager) Dispose() {}

// RunMessageLoop runs the Windows message loop (blocking).
func (m *Manager) RunMessageLoop() {}
