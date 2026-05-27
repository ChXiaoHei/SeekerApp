//go:build !windows

package everything

import (
	"fmt"
)

// SDK is a stub for non-Windows platforms.
type SDK struct{}

// NewSDK returns an error on non-Windows platforms.
func NewSDK(dllPath string) (*SDK, error) {
	return nil, fmt.Errorf("Everything SDK is only supported on Windows")
}

// Search returns an error on non-Windows platforms.
func (s *SDK) Search(query string, opts SearchOptions) ([]SearchResult, error) {
	return nil, fmt.Errorf("Everything SDK is only supported on Windows")
}

// CleanUp does nothing on non-Windows platforms.
func (s *SDK) CleanUp() {}
