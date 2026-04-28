package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func TestProfileStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewProfileStore(dir)

	p := snapshot.Profile{
		Name:     "production",
		Prefixes: []string{"APP_", "DB_"},
		Exclude:  []string{"APP_DEBUG"},
		Redact:   []string{"DB_PASSWORD"},
	}

	if err := store.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load("production")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got.Name != p.Name {
		t.Errorf("Name: got %q, want %q", got.Name, p.Name)
	}
	if len(got.Prefixes) != 2 {
		t.Errorf("Prefixes len: got %d, want 2", len(got.Prefixes))
	}
}

func TestProfileStore_Load_Missing(t *testing.T) {
	store := snapshot.NewProfileStore(t.TempDir())
	_, err := store.Load("ghost")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestProfileStore_InvalidName(t *testing.T) {
	store := snapshot.NewProfileStore(t.TempDir())
	err := store.Save(snapshot.Profile{Name: "bad name!"})
	if err == nil {
		t.Fatal("expected error for invalid name")
	}
}

func TestProfileStore_List(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewProfileStore(dir)

	for _, name := range []string{"alpha", "beta", "gamma"} {
		if err := store.Save(snapshot.Profile{Name: name}); err != nil {
			t.Fatalf("Save %q: %v", name, err)
		}
	}
	// add a non-json file to ensure it's skipped
	_ = os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("hi"), 0o644)

	names, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("List len: got %d, want 3", len(names))
	}
}

func TestProfileStore_Delete(t *testing.T) {
	store := snapshot.NewProfileStore(t.TempDir())
	_ = store.Save(snapshot.Profile{Name: "temp"})

	if err := store.Delete("temp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := store.Load("temp")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestProfileStore_Delete_Missing(t *testing.T) {
	store := snapshot.NewProfileStore(t.TempDir())
	if err := store.Delete("nope"); err == nil {
		t.Fatal("expected error deleting non-existent profile")
	}
}
