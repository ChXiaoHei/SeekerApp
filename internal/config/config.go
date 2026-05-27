package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all user-configurable settings.
type Config struct {
	Hotkey       string `json:"hotkey"`        // Default: "win+alt+s"
	MaxHistory   int    `json:"maxHistory"`    // Default: 50
	DLLPath      string `json:"dllPath"`       // Path to Everything64.dll, empty = default
	Theme        string `json:"theme"`         // "light", "dark", "system"
	WindowWidth  int    `json:"windowWidth"`   // Default: 600
	WindowHeight int    `json:"windowHeight"`  // Default: 400
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Hotkey:       "win+alt+s",
		MaxHistory:   50,
		DLLPath:      "",
		Theme:        "system",
		WindowWidth:  600,
		WindowHeight: 400,
	}
}

// Manager handles loading and saving configuration.
type Manager struct {
	configPath string
	config     Config
}

// NewManager creates a new configuration manager.
func NewManager() (*Manager, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(configDir, "GoEverythingToolbar")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	m := &Manager{
		configPath: filepath.Join(appDir, "config.json"),
		config:     DefaultConfig(),
	}

	if err := m.Load(); err != nil {
		// If file doesn't exist, use defaults and save
		if os.IsNotExist(err) {
			_ = m.Save()
			return m, nil
		}
		return nil, err
	}

	return m, nil
}

// Load reads configuration from disk.
func (m *Manager) Load() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &m.config)
}

// Save writes configuration to disk.
func (m *Manager) Save() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0644)
}

// Get returns the current configuration.
func (m *Manager) Get() Config {
	return m.config
}

// Set updates the configuration and saves to disk.
func (m *Manager) Set(c Config) error {
	m.config = c
	return m.Save()
}

// SetHotkey updates the hotkey setting.
func (m *Manager) SetHotkey(hotkey string) error {
	m.config.Hotkey = hotkey
	return m.Save()
}

// SetMaxHistory updates the max history setting.
func (m *Manager) SetMaxHistory(max int) error {
	m.config.MaxHistory = max
	return m.Save()
}

// SetTheme updates the theme setting.
func (m *Manager) SetTheme(theme string) error {
	m.config.Theme = theme
	return m.Save()
}
