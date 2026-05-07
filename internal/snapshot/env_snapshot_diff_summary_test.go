package snapshot

import (
	"testing"
)

func makeSummaryDiffs() []DiffEntry {
	return []DiffEntry{
		{Key: "APP_HOST", Status: "added", NewValue: "localhost"},
		{Key: "APP_PORT", Status: "changed", OldValue: "8080", NewValue: "9090"},
		{Key: "DB_PASS", Status: "removed", OldValue: "secret"},
		{Key: "DB_HOST", Status: "changed", OldValue: "old", NewValue: "new"},
		{Key: "CACHE_URL", Status: "added", NewValue: "redis://localhost"},
	}
}

func TestSummariseDiff_Counts(t *testing.T) {
	diffs := makeSummaryDiffs()
	s := SummariseDiff(diffs)

	if s.Added != 2 {
		t.Errorf("expected Added=2, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", s.Removed)
	}
	if s.Changed != 2 {
		t.Errorf("expected Changed=2, got %d", s.Changed)
	}
	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
}

func TestSummariseDiff_HasChanges(t *testing.T) {
	if !SummariseDiff(makeSummaryDiffs()).HasChanges() {
		t.Error("expected HasChanges=true")
	}
	if SummariseDiff(nil).HasChanges() {
		t.Error("expected HasChanges=false for empty diffs")
	}
}

func TestSummariseDiff_String(t *testing.T) {
	s := DiffSummary{Added: 2, Removed: 1, Changed: 2, Total: 5}
	got := s.String()
	want := "+2 added, -1 removed, ~2 changed (5 total)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestDiffSummaryByPrefix(t *testing.T) {
	diffs := makeSummaryDiffs()
	byPrefix := DiffSummaryByPrefix(diffs)

	app, ok := byPrefix["APP"]
	if !ok {
		t.Fatal("expected APP prefix in summary")
	}
	if app.Added != 1 || app.Changed != 1 {
		t.Errorf("APP: got added=%d changed=%d, want added=1 changed=1", app.Added, app.Changed)
	}

	db, ok := byPrefix["DB"]
	if !ok {
		t.Fatal("expected DB prefix in summary")
	}
	if db.Removed != 1 || db.Changed != 1 {
		t.Errorf("DB: got removed=%d changed=%d, want removed=1 changed=1", db.Removed, db.Changed)
	}

	cache, ok := byPrefix["CACHE"]
	if !ok {
		t.Fatal("expected CACHE prefix in summary")
	}
	if cache.Added != 1 {
		t.Errorf("CACHE: got added=%d, want 1", cache.Added)
	}
}

func TestDiffSummaryByPrefix_NoUnderscore(t *testing.T) {
	diffs := []DiffEntry{
		{Key: "HOST", Status: "added", NewValue: "localhost"},
	}
	byPrefix := DiffSummaryByPrefix(diffs)
	if _, ok := byPrefix["HOST"]; !ok {
		t.Error("expected key with no underscore to use full key as prefix")
	}
}
