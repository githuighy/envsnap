package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGroupStore_SaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := NewGroupStore()
	_ = s.Add("aws", []string{"AWS_", "S3_"})
	_ = s.Add("db", []string{"DB_", "DATABASE_"})

	if err := s.SaveGroups(dir); err != nil {
		t.Fatalf("SaveGroups: %v", err)
	}

	s2 := NewGroupStore()
	if err := s2.LoadGroups(dir); err != nil {
		t.Fatalf("LoadGroups: %v", err)
	}

	names := s2.List()
	if len(names) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(names))
	}
	g, ok := s2.Get("aws")
	if !ok {
		t.Fatal("expected aws group")
	}
	if len(g.Patterns) != 2 {
		t.Errorf("expected 2 patterns, got %d", len(g.Patterns))
	}
}

func TestGroupStore_LoadGroups_MissingFile(t *testing.T) {
	dir := t.TempDir()
	s := NewGroupStore()
	// Should not error when file doesn't exist
	if err := s.LoadGroups(dir); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(s.List()) != 0 {
		t.Error("expected empty store")
	}
}

func TestGroupStore_SaveGroups_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "groups")
	s := NewGroupStore()
	_ = s.Add("test", []string{"TEST_"})
	if err := s.SaveGroups(dir); err != nil {
		t.Fatalf("SaveGroups: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "groups.json")); err != nil {
		t.Errorf("expected groups.json to exist: %v", err)
	}
}
