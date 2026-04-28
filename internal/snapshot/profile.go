package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var validProfileName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// Profile represents a named configuration for capturing snapshots.
type Profile struct {
	Name     string   `json:"name"`
	Prefixes []string `json:"prefixes,omitempty"`
	Exclude  []string `json:"exclude,omitempty"`
	Redact   []string `json:"redact,omitempty"`
}

// ProfileStore manages named profiles on disk.
type ProfileStore struct {
	Dir string
}

// NewProfileStore returns a ProfileStore rooted at dir.
func NewProfileStore(dir string) *ProfileStore {
	return &ProfileStore{Dir: dir}
}

// Save writes a profile to disk.
func (ps *ProfileStore) Save(p Profile) error {
	if !validProfileName.MatchString(p.Name) {
		return fmt.Errorf("invalid profile name %q: use letters, digits, - or _", p.Name)
	}
	if err := os.MkdirAll(ps.Dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ps.profilePath(p.Name), data, 0o644)
}

// Load reads a profile by name from disk.
func (ps *ProfileStore) Load(name string) (Profile, error) {
	data, err := os.ReadFile(ps.profilePath(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Profile{}, fmt.Errorf("profile %q not found", name)
		}
		return Profile{}, err
	}
	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return Profile{}, err
	}
	return p, nil
}

// List returns all profile names in the store.
func (ps *ProfileStore) List() ([]string, error) {
	entries, err := os.ReadDir(ps.Dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

// Delete removes a profile by name.
func (ps *ProfileStore) Delete(name string) error {
	err := os.Remove(ps.profilePath(name))
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("profile %q not found", name)
	}
	return err
}

func (ps *ProfileStore) profilePath(name string) string {
	return filepath.Join(ps.Dir, name+".json")
}
