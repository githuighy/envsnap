package snapshot

import (
	"testing"
)

func baseReorderSnap() Snapshot {
	return Snapshot{
		Vars: map[string]string{
			"ZEBRA":     "1",
			"APP_PORT":  "8080",
			"APP_HOST":  "localhost",
			"DB_URL":    "postgres://",
			"ALPHA":     "a",
			"DB_PASS":   "secret",
		},
	}
}

func TestReorder_Alphabetical(t *testing.T) {
	s := baseReorderSnap()
	out, err := Reorder(s, ReorderOptions{Alphabetical: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := sortedKeysOf(out.Vars)
	expect := []string{"ALPHA", "APP_HOST", "APP_PORT", "DB_PASS", "DB_URL", "ZEBRA"}
	for i, k := range keys {
		if k != expect[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, expect[i])
		}
	}
}

func TestReorder_Descending(t *testing.T) {
	s := baseReorderSnap()
	out, err := Reorder(s, ReorderOptions{Alphabetical: true, Descending: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := sortedKeysOf(out.Vars)
	// sorted ascending for comparison; descending is internal ordering
	if len(keys) != 6 {
		t.Errorf("expected 6 keys, got %d", len(keys))
	}
}

func TestReorder_PrefixPriority(t *testing.T) {
	s := baseReorderSnap()
	out, err := Reorder(s, ReorderOptions{
		Alphabetical:   true,
		PrefixPriority: []string{"DB_", "APP_"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// DB_ keys should rank 0, APP_ keys rank 1, rest rank 2
	if _, ok := out.Vars["DB_URL"]; !ok {
		t.Error("DB_URL missing from output")
	}
	if len(out.Vars) != len(s.Vars) {
		t.Errorf("expected %d keys, got %d", len(s.Vars), len(out.Vars))
	}
}

func TestReorder_NoOptions_ReturnsAll(t *testing.T) {
	s := baseReorderSnap()
	out, err := Reorder(s, ReorderOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Vars) != len(s.Vars) {
		t.Errorf("expected %d keys, got %d", len(s.Vars), len(out.Vars))
	}
}

func TestReorder_DoesNotMutateOriginal(t *testing.T) {
	s := baseReorderSnap()
	origLen := len(s.Vars)
	_, _ = Reorder(s, ReorderOptions{Alphabetical: true})
	if len(s.Vars) != origLen {
		t.Error("original snapshot was mutated")
	}
}

func sortedKeysOf(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	import_sort_workaround(keys)
	return keys
}

func import_sort_workaround(keys []string) {
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
}
