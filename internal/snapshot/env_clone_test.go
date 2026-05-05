package snapshot_test

import (
	"testing"

	"github.com/your-org/envsnap/internal/snapshot"
)

var baseCloneSnap = map[string]string{
	"APP_HOST":   "localhost",
	"APP_PORT":   "8080",
	"DB_HOST":    "db.local",
	"DB_PASS":    "secret",
	"LOG_LEVEL":  "info",
}

func TestClone_NoOptions_CopiesAll(t *testing.T) {
	out, err := snapshot.Clone(baseCloneSnap, snapshot.CloneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(baseCloneSnap) {
		t.Errorf("expected %d keys, got %d", len(baseCloneSnap), len(out))
	}
}

func TestClone_FilterByPrefix(t *testing.T) {
	out, err := snapshot.Clone(baseCloneSnap, snapshot.CloneOptions{
		Prefixes: []string{"APP_"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should be excluded")
	}
}

func TestClone_ExcludePrefix(t *testing.T) {
	out, err := snapshot.Clone(baseCloneSnap, snapshot.CloneOptions{
		Exclude: []string{"DB_"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("DB_PASS should be excluded")
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("APP_HOST should be present")
	}
}

func TestClone_KeyMap_RenamesKey(t *testing.T) {
	out, err := snapshot.Clone(baseCloneSnap, snapshot.CloneOptions{
		KeyMap: map[string]string{"APP_HOST": "SERVICE_HOST"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Error("original key APP_HOST should be absent after rename")
	}
	if out["SERVICE_HOST"] != "localhost" {
		t.Errorf("expected SERVICE_HOST=localhost, got %q", out["SERVICE_HOST"])
	}
}

func TestClone_OverrideValues(t *testing.T) {
	out, err := snapshot.Clone(baseCloneSnap, snapshot.CloneOptions{
		OverrideValues: map[string]string{"APP_PORT": "9090", "NEW_KEY": "new"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_PORT"] != "9090" {
		t.Errorf("expected APP_PORT=9090, got %q", out["APP_PORT"])
	}
	if out["NEW_KEY"] != "new" {
		t.Errorf("expected NEW_KEY=new, got %q", out["NEW_KEY"])
	}
}

func TestClone_NilSource_ReturnsError(t *testing.T) {
	_, err := snapshot.Clone(nil, snapshot.CloneOptions{})
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestClone_EmptyKeyMapDest_ReturnsError(t *testing.T) {
	_, err := snapshot.Clone(baseCloneSnap, snapshot.CloneOptions{
		KeyMap: map[string]string{"APP_HOST": ""},
	})
	if err == nil {
		t.Error("expected error for empty destination key in KeyMap")
	}
}
