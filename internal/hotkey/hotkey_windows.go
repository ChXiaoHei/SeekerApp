//go:build windows

package hotkey

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"golang.design/x/hotkey"
)

// Handler is called when the hotkey is triggered.
type Handler func()

// Manager manages global hotkey registration.
type Manager struct {
	mu       sync.Mutex
	hotkey   *hotkey.Hotkey
	handler  Handler
	running  bool
	stopCh   chan struct{}
}

// NewManager creates a new hotkey manager.
func NewManager() *Manager {
	return &Manager{
		stopCh: make(chan struct{}),
	}
}

// Register registers a global hotkey with the given modifiers and key.
func (m *Manager) Register(mods []hotkey.Modifier, key hotkey.Key, handler Handler) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("hotkey already registered, unregister first")
	}

	hk := hotkey.New(mods, key)
	if err := hk.Register(); err != nil {
		return fmt.Errorf("failed to register hotkey: %w", err)
	}

	m.hotkey = hk
	m.handler = handler
	m.running = true

	go m.listen()

	return nil
}

// RegisterFromString parses a hotkey string like "win+alt+s" and registers it.
func (m *Manager) RegisterFromString(s string, handler Handler) error {
	mods, key, err := parseHotkeyString(s)
	if err != nil {
		return err
	}
	return m.Register(mods, key, handler)
}

// Unregister removes the current hotkey registration.
func (m *Manager) Unregister() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running || m.hotkey == nil {
		return nil
	}

	close(m.stopCh)
	err := m.hotkey.Unregister()
	m.running = false
	m.hotkey = nil
	m.stopCh = make(chan struct{})

	return err
}

// listen waits for hotkey events.
func (m *Manager) listen() {
	for {
		select {
		case <-m.stopCh:
			return
		case <-m.hotkey.Keydown():
			if m.handler != nil {
				m.handler()
			}
		}
	}
}

// parseHotkeyString parses a string like "win+alt+s" into modifiers and key.
// Supported modifiers: ctrl, alt, shift, win
// Supported keys: a-z, 0-9, F1-F24, and common special keys
func parseHotkeyString(s string) ([]hotkey.Modifier, hotkey.Key, error) {
	if s == "" {
		return nil, 0, fmt.Errorf("empty hotkey string")
	}

	s = strings.ToLower(strings.TrimSpace(s))
	parts := strings.Split(s, "+")

	if len(parts) < 2 {
		return nil, 0, fmt.Errorf("hotkey must have at least one modifier and one key, got: %s", s)
	}

	var mods []hotkey.Modifier
	var key hotkey.Key

	for i := 0; i < len(parts)-1; i++ {
		mod := strings.TrimSpace(parts[i])
		switch mod {
		case "ctrl", "control":
			mods = append(mods, hotkey.ModCtrl)
		case "alt":
			mods = append(mods, hotkey.ModAlt)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "win", "cmd", "super", "meta":
			mods = append(mods, hotkey.ModWin)
		default:
			return nil, 0, fmt.Errorf("unknown modifier: %s", mod)
		}
	}

	keyStr := strings.TrimSpace(parts[len(parts)-1])
	key, err := parseKey(keyStr)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid key: %w", err)
	}

	if len(mods) == 0 {
		return nil, 0, fmt.Errorf("no modifiers specified")
	}

	return mods, key, nil
}

// parseKey converts a key string to a hotkey.Key constant.
func parseKey(s string) (hotkey.Key, error) {
	if len(s) == 1 {
		c := s[0]
		if c >= 'a' && c <= 'z' {
			return hotkey.Key(c - 'a' + 'A'), nil
		}
		if c >= 'A' && c <= 'Z' {
			return hotkey.Key(c), nil
		}
		if c >= '0' && c <= '9' {
			return hotkey.Key(c), nil
		}
	}

	if strings.HasPrefix(s, "f") && len(s) >= 2 && len(s) <= 3 {
		numStr := s[1:]
		num, err := strconv.Atoi(numStr)
		if err == nil && num >= 1 && num <= 24 {
			return hotkey.KeyF1 + hotkey.Key(num-1), nil
		}
	}

	switch s {
	case "enter", "return":
		return hotkey.KeyReturn, nil
	case "space":
		return hotkey.KeySpace, nil
	case "tab":
		return hotkey.KeyTab, nil
	case "escape", "esc":
		return hotkey.KeyEscape, nil
	case "backspace", "back":
		// VK_BACK = 0x08
		return hotkey.Key(0x08), nil
	case "insert":
		// VK_INSERT = 0x2D
		return hotkey.Key(0x2D), nil
	case "delete", "del":
		return hotkey.KeyDelete, nil
	case "home":
		// VK_HOME = 0x24
		return hotkey.Key(0x24), nil
	case "end":
		// VK_END = 0x23
		return hotkey.Key(0x23), nil
	case "pageup", "pgup":
		// VK_PRIOR = 0x21
		return hotkey.Key(0x21), nil
	case "pagedown", "pgdn", "pgdown":
		// VK_NEXT = 0x22
		return hotkey.Key(0x22), nil
	case "up":
		return hotkey.KeyUp, nil
	case "down":
		return hotkey.KeyDown, nil
	case "left":
		return hotkey.KeyLeft, nil
	case "right":
		return hotkey.KeyRight, nil
	default:
		return 0, fmt.Errorf("unsupported key: %s", s)
	}
}
