package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// TemplateStore persists named templates to a directory.
type TemplateStore struct {
	dir string
}

// NewTemplateStore creates a TemplateStore rooted at dir.
func NewTemplateStore(dir string) (*TemplateStore, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create template store dir: %w", err)
	}
	return &TemplateStore{dir: dir}, nil
}

func (s *TemplateStore) path(name string) string {
	return filepath.Join(s.dir, name+".json")
}

// Save persists a template by its name.
func (s *TemplateStore) Save(tmpl Template) error {
	if !validTemplateName.MatchString(tmpl.Name) {
		return fmt.Errorf("invalid template name %q", tmpl.Name)
	}
	data, err := json.MarshalIndent(tmpl, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(tmpl.Name), data, 0644)
}

// Load retrieves a template by name.
func (s *TemplateStore) Load(name string) (Template, error) {
	data, err := os.ReadFile(s.path(name))
	if err != nil {
		if os.IsNotExist(err) {
			return Template{}, fmt.Errorf("template %q not found", name)
		}
		return Template{}, err
	}
	var tmpl Template
	if err := json.Unmarshal(data, &tmpl); err != nil {
		return Template{}, fmt.Errorf("parse template %q: %w", name, err)
	}
	return tmpl, nil
}

// List returns all saved template names in sorted order.
func (s *TemplateStore) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	sort.Strings(names)
	return names, nil
}

// Delete removes a template by name.
func (s *TemplateStore) Delete(name string) error {
	err := os.Remove(s.path(name))
	if os.IsNotExist(err) {
		return fmt.Errorf("template %q not found", name)
	}
	return err
}
