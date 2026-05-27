//go:build windows

package tray

import (
	"fmt"

	"github.com/lxn/walk"
)

// MenuItem represents a menu item in the system tray.
type MenuItem struct {
	Text    string
	Handler func()
}

// Callbacks are functions called from tray menu actions.
type Callbacks struct {
	OnShowWindow   func()
	OnSettings     func()
	OnExit         func()
}

// Manager manages the system tray icon and menu.
type Manager struct {
	notifyIcon *walk.NotifyIcon
	mainWindow *walk.MainWindow
	callbacks  Callbacks
}

// NewManager creates a new system tray manager.
func NewManager(callbacks Callbacks) (*Manager, error) {
	// We need a hidden main window for the notify icon
	mw, err := walk.NewMainWindow()
	if err != nil {
		return nil, fmt.Errorf("failed to create main window: %w", err)
	}

	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		mw.Close()
		return nil, fmt.Errorf("failed to create notify icon: %w", err)
	}

	m := &Manager{
		notifyIcon: ni,
		mainWindow: mw,
		callbacks:  callbacks,
	}

	return m, nil
}

// SetIcon sets the tray icon from an icon file path.
func (m *Manager) SetIcon(iconPath string) error {
	icon, err := walk.NewIconFromFile(iconPath)
	if err != nil {
		// Try loading from embedded resources or use default
		return fmt.Errorf("failed to load icon: %w", err)
	}
	return m.notifyIcon.SetIcon(icon)
}

// SetIconFromResource sets the tray icon from embedded resources.
func (m *Manager) SetIconFromResource(resourceID int) error {
	icon, err := walk.NewIconFromResource(resourceID)
	if err != nil {
		return fmt.Errorf("failed to load icon from resource: %w", err)
	}
	return m.notifyIcon.SetIcon(icon)
}

// Show displays the tray icon.
func (m *Manager) Show() error {
	return m.notifyIcon.SetVisible(true)
}

// Hide hides the tray icon.
func (m *Manager) Hide() error {
	return m.notifyIcon.SetVisible(false)
}

// SetTooltip sets the tray icon tooltip text.
func (m *Manager) SetTooltip(text string) error {
	return m.notifyIcon.SetToolTip(text)
}

// SetupMenu creates the context menu for the tray icon.
func (m *Manager) SetupMenu() error {
	menu := m.notifyIcon.ContextMenu()

	// Show search window
	showAction := walk.NewAction()
	showAction.SetText("Show Search")
	showAction.Triggered().Attach(func() {
		if m.callbacks.OnShowWindow != nil {
			m.callbacks.OnShowWindow()
		}
	})
	menu.Actions().Add(showAction)

	// Separator
	sep := walk.NewSeparatorAction()
	menu.Actions().Add(sep)

	// Settings
	settingsAction := walk.NewAction()
	settingsAction.SetText("Settings")
	settingsAction.Triggered().Attach(func() {
		if m.callbacks.OnSettings != nil {
			m.callbacks.OnSettings()
		}
	})
	menu.Actions().Add(settingsAction)

	// Separator
	menu.Actions().Add(walk.NewSeparatorAction())

	// Exit
	exitAction := walk.NewAction()
	exitAction.SetText("Exit")
	exitAction.Triggered().Attach(func() {
		if m.callbacks.OnExit != nil {
			m.callbacks.OnExit()
		}
	})
	menu.Actions().Add(exitAction)

	return nil
}

// Dispose releases tray resources.
func (m *Manager) Dispose() {
	if m.notifyIcon != nil {
		m.notifyIcon.Dispose()
	}
	if m.mainWindow != nil {
		m.mainWindow.Close()
	}
}

// RunMessageLoop runs the Windows message loop (blocking).
func (m *Manager) RunMessageLoop() {
	m.mainWindow.Run()
}
