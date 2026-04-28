package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestAliasStore(t *testing.T) *AliasStore {
	t.Helper()
	dir := t.TempDir()
	return NewAliasStore(dir)
}

func TestAliasStore_SetAndResolve(t *testing.T) {
	s := newTestAliasStore(t)
	if err := s.Set("prod", "/snaps/prod.json"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	path, err := s.Resolve("prod")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if path != "/snaps/prod.json" {
		t.Errorf("expected /snaps/prod.json, got %q", path)
	}
}

func TestAliasStore_Resolve_Missing(t *testing.T) {
	s := newTestAliasStore(t)
	_, err := s.Resolve("nope")
	if err == nil {
		t.Fatal("expected error for missing alias")
	}
}

func TestAliasStore_InvalidName(t *testing.T) {
	s := newTestAliasStore(t)
	err := s.Set("bad name!", "/path")
	if err == nil {
		t.Fatal("expected error for invalid alias name")
	}
}

func TestAliasStore_Delete(t *testing.T) {
	s := newTestAliasStore(t)
	_ = s.Set("staging", "/snaps/staging.json")
	if err := s.Delete("staging"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Resolve("staging")
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestAliasStore_Delete_Missing(t *testing.T) {
	s := newTestAliasStore(t)
	if err := s.Delete("ghost"); err == nil {
		t.Fatal("expected error deleting non-existent alias")
	}
}

func TestAliasStore_List(t *testing.T) {
	s := newTestAliasStore(t)
	_ = s.Set("beta", "/b")
	_ = s.Set("alpha", "/a")
	_ = s.Set("gamma", "/g")

	aliases, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(aliases) != 3 {
		t.Fatalf("expected 3 aliases, got %d", len(aliases))
	}
	if aliases[0].Name != "alpha" || aliases[1].Name != "beta" || aliases[2].Name != "gamma" {
		t.Errorf("unexpected order: %v", aliases)
	}
}

func TestAliasStore_List_Empty(t *testing.T) {
	s := newTestAliasStore(t)
	aliases, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(aliases) != 0 {
		t.Errorf("expected empty list, got %v", aliases)
	}
}

func TestAliasStore_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "aliases")
	s := NewAliasStore(dir)
	if err := s.Set("x", "/x"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected dir to be created: %v", err)
	}
}
