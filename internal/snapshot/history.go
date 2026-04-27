package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// HistoryEntry represents a saved snapshot with metadata.
type HistoryEntry struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Label     string            `json:"label,omitempty"`
	Vars      map[string]string `json:"vars"`
}

// HistoryStore manages a collection of historical snapshots.
type HistoryStore struct {
	Dir string
}

// NewHistoryStore creates a HistoryStore rooted at dir.
func NewHistoryStore(dir string) *HistoryStore {
	return &HistoryStore{Dir: dir}
}

// Add saves a snapshot to the history store with an optional label.
func (h *HistoryStore) Add(snap map[string]string, label string) (string, error) {
	if err := os.MkdirAll(h.Dir, 0755); err != nil {
		return "", fmt.Errorf("history: create dir: %w", err)
	}
	now := time.Now().UTC()
	id := now.Format("20060102T150405Z")
	entry := HistoryEntry{
		ID:        id,
		Timestamp: now,
		Label:     label,
		Vars:      snap,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return "", fmt.Errorf("history: marshal: %w", err)
	}
	path := filepath.Join(h.Dir, id+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("history: write: %w", err)
	}
	return id, nil
}

// List returns all history entries sorted by timestamp ascending.
func (h *HistoryStore) List() ([]HistoryEntry, error) {
	glob := filepath.Join(h.Dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("history: glob: %w", err)
	}
	var entries []HistoryEntry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("history: read %s: %w", m, err)
		}
		var e HistoryEntry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("history: parse %s: %w", m, err)
		}
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})
	return entries, nil
}

// Get retrieves a single entry by ID.
func (h *HistoryStore) Get(id string) (*HistoryEntry, error) {
	path := filepath.Join(h.Dir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("history: get %s: %w", id, err)
	}
	var e HistoryEntry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("history: parse %s: %w", id, err)
	}
	return &e, nil
}
