package everything

import "time"

// SearchOptions configures search behavior for Everything queries.
type SearchOptions struct {
	MatchCase      bool `json:"matchCase"`
	MatchWholeWord bool `json:"matchWholeWord"`
	MatchPath      bool `json:"matchPath"`
	UseRegex       bool `json:"useRegex"`
	MaxResults     int  `json:"maxResults"`
}

// DefaultSearchOptions returns sensible defaults for search.
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		MaxResults: 100,
	}
}

// SearchResult represents a single file/folder match from Everything.
type SearchResult struct {
	FileName     string    `json:"fileName"`
	FullPath     string    `json:"fullPath"`
	IsFolder     bool      `json:"isFolder"`
	Size         int64     `json:"size"`
	DateModified time.Time `json:"dateModified"`
}