package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const groupStoreFile = "groups.json"

// SaveGroups persists all groups in the store to a JSON file in dir.
func (s *GroupStore) SaveGroups(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	groups := make([]EnvGroup, 0, len(s.groups))
	for _, g := range s.groups {
		groups = append(groups, g)
	}
	data, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal groups: %w", err)
	}
	dest := filepath.Join(dir, groupStoreFile)
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return fmt.Errorf("write groups file: %w", err)
	}
	return nil
}

// LoadGroups reads groups from a JSON file in dir and populates the store.
func (s *GroupStore) LoadGroups(dir string) error {
	dest := filepath.Join(dir, groupStoreFile)
	data, err := os.ReadFile(dest)
	if os.IsNotExist(err) {
		return nil // no groups saved yet
	}
	if err != nil {
		return fmt.Errorf("read groups file: %w", err)
	}
	var groups []EnvGroup
	if err := json.Unmarshal(data, &groups); err != nil {
		return fmt.Errorf("unmarshal groups: %w", err)
	}
	for _, g := range groups {
		s.groups[g.Name] = g
	}
	return nil
}
