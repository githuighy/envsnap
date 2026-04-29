package snapshot

import (
	"fmt"
	"sort"
	"strings"
)

// EnvGroup represents a named collection of environment variable key patterns.
type EnvGroup struct {
	Name     string   `json:"name"`
	Patterns []string `json:"patterns"`
}

// GroupStore manages named environment variable groups.
type GroupStore struct {
	groups map[string]EnvGroup
}

// NewGroupStore creates an empty GroupStore.
func NewGroupStore() *GroupStore {
	return &GroupStore{groups: make(map[string]EnvGroup)}
}

// Add adds or replaces a group by name.
func (s *GroupStore) Add(name string, patterns []string) error {
	if name == "" {
		return fmt.Errorf("group name must not be empty")
	}
	if strings.ContainsAny(name, " /\\.") {
		return fmt.Errorf("group name %q contains invalid characters", name)
	}
	if len(patterns) == 0 {
		return fmt.Errorf("group %q must have at least one pattern", name)
	}
	s.groups[name] = EnvGroup{Name: name, Patterns: patterns}
	return nil
}

// Get retrieves a group by name.
func (s *GroupStore) Get(name string) (EnvGroup, bool) {
	g, ok := s.groups[name]
	return g, ok
}

// Delete removes a group by name.
func (s *GroupStore) Delete(name string) error {
	if _, ok := s.groups[name]; !ok {
		return fmt.Errorf("group %q not found", name)
	}
	delete(s.groups, name)
	return nil
}

// List returns all group names in sorted order.
func (s *GroupStore) List() []string {
	names := make([]string, 0, len(s.groups))
	for name := range s.groups {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ExtractGroup returns a new Snapshot containing only keys that match
// any pattern in the named group. Patterns are prefix-matched.
func (s *GroupStore) ExtractGroup(name string, snap Snapshot) (Snapshot, error) {
	g, ok := s.groups[name]
	if !ok {
		return nil, fmt.Errorf("group %q not found", name)
	}
	result := make(Snapshot)
	for k, v := range snap {
		for _, pat := range g.Patterns {
			if strings.HasPrefix(k, pat) {
				result[k] = v
				break
			}
		}
	}
	return result, nil
}
