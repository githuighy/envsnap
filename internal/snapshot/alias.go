package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var validAliasName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// Alias maps a short name to a snapshot file path.
type Alias struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// AliasStore manages named aliases for snapshot files.
type AliasStore struct {
	dir string
}

// NewAliasStore returns a new AliasStore backed by the given directory.
func NewAliasStore(dir string) *AliasStore {
	return &AliasStore{dir: dir}
}

func (s *AliasStore) indexPath() string {
	return filepath.Join(s.dir, "aliases.json")
}

func (s *AliasStore) load() (map[string]string, error) {
	data, err := os.ReadFile(s.indexPath())
	if errors.Is(err, os.ErrNotExist) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *AliasStore) save(m map[string]string) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.indexPath(), data, 0o644)
}

// Set creates or updates an alias mapping name -> path.
func (s *AliasStore) Set(name, path string) error {
	if !validAliasName.MatchString(name) {
		return fmt.Errorf("invalid alias name %q: use only letters, digits, hyphens, underscores", name)
	}
	m, err := s.load()
	if err != nil {
		return err
	}
	m[name] = path
	return s.save(m)
}

// Resolve returns the file path for the given alias name.
func (s *AliasStore) Resolve(name string) (string, error) {
	m, err := s.load()
	if err != nil {
		return "", err
	}
	p, ok := m[name]
	if !ok {
		return "", fmt.Errorf("alias %q not found", name)
	}
	return p, nil
}

// Delete removes an alias by name.
func (s *AliasStore) Delete(name string) error {
	m, err := s.load()
	if err != nil {
		return err
	}
	if _, ok := m[name]; !ok {
		return fmt.Errorf("alias %q not found", name)
	}
	delete(m, name)
	return s.save(m)
}

// List returns all aliases sorted by name.
func (s *AliasStore) List() ([]Alias, error) {
	m, err := s.load()
	if err != nil {
		return nil, err
	}
	aliases := make([]Alias, 0, len(m))
	for k, v := range m {
		aliases = append(aliases, Alias{Name: k, Path: v})
	}
	sort.Slice(aliases, func(i, j int) bool { return aliases[i].Name < aliases[j].Name })
	return aliases, nil
}
