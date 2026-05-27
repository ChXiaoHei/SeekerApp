package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	goruntime "runtime"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go-everything-toolbar/internal/config"
	"go-everything-toolbar/internal/everything"
	"go-everything-toolbar/internal/history"
	"go-everything-toolbar/internal/hotkey"
)

// App is the main application structure that binds to the frontend.
type App struct {
	ctx         context.Context
	sdk         *everything.SDK
	config      *config.Manager
	history     *history.Manager
	hotkey      *hotkey.Manager
	isVisible   bool
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize configuration
	cfg, err := config.NewManager()
	if err != nil {
		runtime.LogFatalf(ctx, "Failed to initialize config: %v", err)
		return
	}
	a.config = cfg

	// Initialize search history
	hist, err := history.NewManager(a.config.Get().MaxHistory)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to initialize history: %v", err)
	} else {
		a.history = hist
	}

	// Initialize Everything SDK
	dllPath := a.config.Get().DLLPath
	sdk, err := everything.NewSDK(dllPath)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to initialize Everything SDK: %v", err)
	} else {
		a.sdk = sdk
	}

	// Initialize hotkey manager
	a.hotkey = hotkey.NewManager()

	// Register global hotkey
	cfgVal := a.config.Get()
	err = a.hotkey.RegisterFromString(cfgVal.Hotkey, func() {
		a.ToggleWindow()
	})
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to register hotkey: %v", err)
	}

	runtime.LogInfo(ctx, "GoEverythingToolbar started successfully")
}

// shutdown is called when the app is closing.
func (a *App) shutdown(ctx context.Context) {
	if a.sdk != nil {
		a.sdk.CleanUp()
	}
	if a.hotkey != nil {
		_ = a.hotkey.Unregister()
	}
	runtime.LogInfo(ctx, "GoEverythingToolbar shutdown complete")
}

// Search performs a file search using Everything.
func (a *App) Search(query string, opts everything.SearchOptions) ([]everything.SearchResult, error) {
	if a.sdk == nil {
		return nil, fmt.Errorf("Everything SDK not initialized. Please ensure Everything is running")
	}

	// Record search history
	if a.history != nil && query != "" {
		a.history.Add(query)
	}

	return a.sdk.Search(query, opts)
}

// OpenFile opens a file with its default application.
func (a *App) OpenFile(path string) error {
	switch goruntime.GOOS {
	case "windows":
		return exec.Command("cmd", "/c", "start", "", path).Run()
	case "darwin":
		return exec.Command("open", path).Run()
	default:
		return exec.Command("xdg-open", path).Run()
	}
}

// OpenFolder opens the containing folder and selects the file.
func (a *App) OpenFolder(path string) error {
	dir := filepath.Dir(path)

	switch goruntime.GOOS {
	case "windows":
		// Use explorer with /select to highlight the file
		return exec.Command("explorer", "/select,", path).Run()
	case "darwin":
		return exec.Command("open", "-R", path).Run()
	default:
		return exec.Command("xdg-open", dir).Run()
	}
}

// CopyPathToClipboard copies the full path to the system clipboard.
func (a *App) CopyPathToClipboard(path string) error {
	return runtime.ClipboardSetText(a.ctx, path)
}

// GetHistory returns recent search queries.
func (a *App) GetHistory() []string {
	if a.history == nil {
		return nil
	}
	return a.history.GetQueries()
}

// ClearHistory clears all search history.
func (a *App) ClearHistory() error {
	if a.history == nil {
		return nil
	}
	return a.history.Clear()
}

// GetConfig returns the current configuration.
func (a *App) GetConfig() config.Config {
	if a.config == nil {
		return config.DefaultConfig()
	}
	return a.config.Get()
}

// SaveConfig updates and persists configuration.
func (a *App) SaveConfig(c config.Config) error {
	if a.config == nil {
		return fmt.Errorf("config manager not initialized")
	}

	// Re-register hotkey if changed
	if c.Hotkey != a.config.Get().Hotkey {
		_ = a.hotkey.Unregister()
		err := a.hotkey.RegisterFromString(c.Hotkey, func() {
			a.ToggleWindow()
		})
		if err != nil {
			return fmt.Errorf("failed to register new hotkey: %w", err)
		}
	}

	// Update history max if changed
	if a.history != nil && c.MaxHistory != a.config.Get().MaxHistory {
		a.history.SetMaxEntries(c.MaxHistory)
	}

	return a.config.Set(c)
}

// ShowWindow shows the main search window.
func (a *App) ShowWindow() {
	a.isVisible = true
	runtime.WindowShow(a.ctx)
	runtime.WindowSetAlwaysOnTop(a.ctx, true)
	runtime.WindowSetAlwaysOnTop(a.ctx, false)
}

// HideWindow hides the main search window.
func (a *App) HideWindow() {
	a.isVisible = false
	runtime.WindowHide(a.ctx)
}

// ToggleWindow toggles the visibility of the main window.
func (a *App) ToggleWindow() {
	if a.isVisible {
		a.HideWindow()
	} else {
		a.ShowWindow()
	}
}

// IsWindowVisible returns whether the window is currently visible.
func (a *App) IsWindowVisible() bool {
	return a.isVisible
}

// GetSearchOptions returns the default search options.
func (a *App) GetSearchOptions() everything.SearchOptions {
	return everything.DefaultSearchOptions()
}

