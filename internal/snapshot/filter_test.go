package snapshot

import (
	"testing"
)

func baseSnap() Snapshot {
	return Snapshot{
		"APP_ENV":     "production",
		"APP_PORT":    "8080",
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"HOME":        "/root",
		"PATH":        "/usr/bin",
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	snap := baseSnap()
	got := Filter(snap, FilterOptions{Prefixes: []string{"APP_"}})

	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if _, ok := got["APP_ENV"]; !ok {
		t.Error("expected APP_ENV in result")
	}
	if _, ok := got["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in result")
	}
}

func TestFilter_Exclude(t *testing.T) {
	snap := baseSnap()
	got := Filter(snap, FilterOptions{Exclude: []string{"DB_"}})

	if _, ok := got["DB_HOST"]; ok {
		t.Error("DB_HOST should have been excluded")
	}
	if _, ok := got["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been excluded")
	}
	if len(got) != 4 {
		t.Fatalf("expected 4 entries after exclusion, got %d", len(got))
	}
}

func TestFilter_PrefixAndExclude(t *testing.T) {
	snap := baseSnap()
	got := Filter(snap, FilterOptions{
		Prefixes: []string{"APP_", "DB_"},
		Exclude:  []string{"DB_PASSWORD"},
	})

	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	if _, ok := got["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been excluded")
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	snap := baseSnap()
	got := Filter(snap, FilterOptions{})

	if len(got) != len(snap) {
		t.Fatalf("expected %d entries, got %d", len(snap), len(got))
	}
}

func TestFilter_EmptySnapshot(t *testing.T) {
	got := Filter(Snapshot{}, FilterOptions{Prefixes: []string{"APP_"}})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(got))
	}
}

func TestFilter_ExcludeExactKey(t *testing.T) {
	snap := baseSnap()
	got := Filter(snap, FilterOptions{Exclude: []string{"HOME"}})

	if _, ok := got["HOME"]; ok {
		t.Error("HOME should have been excluded")
	}
	if len(got) != len(snap)-1 {
		t.Fatalf("expected %d entries, got %d", len(snap)-1, len(got))
	}
}
