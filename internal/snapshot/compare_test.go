package snapshot

import (
	"strings"
	"testing"
)

func baseCompareSnap() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "db.example.com",
		"LOG_LEVEL": "info",
	}
}

func TestCompare_AllUnchanged(t *testing.T) {
	snap := baseCompareSnap()
	result := Compare(snap, snap)

	if result.Added != 0 || result.Removed != 0 || result.Changed != 0 {
		t.Errorf("expected no changes, got %+v", result)
	}
	if result.Unchanged != 3 {
		t.Errorf("expected 3 unchanged, got %d", result.Unchanged)
	}
	if result.Score != 1.0 {
		t.Errorf("expected score 1.0, got %.2f", result.Score)
	}
}

func TestCompare_AddedAndRemoved(t *testing.T) {
	base := baseCompareSnap()
	other := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "db.example.com",
		"NEW_KEY":  "value",
	}

	result := Compare(base, other)

	if result.Added != 1 {
		t.Errorf("expected 1 added, got %d", result.Added)
	}
	if result.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", result.Removed)
	}
	if result.Changed != 0 {
		t.Errorf("expected 0 changed, got %d", result.Changed)
	}
}

func TestCompare_Changed(t *testing.T) {
	base := baseCompareSnap()
	other := map[string]string{
		"APP_ENV":   "staging",
		"DB_HOST":   "db.example.com",
		"LOG_LEVEL": "debug",
	}

	result := Compare(base, other)

	if result.Changed != 2 {
		t.Errorf("expected 2 changed, got %d", result.Changed)
	}
	if result.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", result.Unchanged)
	}
}

func TestCompare_EmptySnapshots(t *testing.T) {
	result := Compare(map[string]string{}, map[string]string{})

	if result.Score != 1.0 {
		t.Errorf("expected score 1.0 for empty snapshots, got %.2f", result.Score)
	}
}

func TestCompareResult_Summary(t *testing.T) {
	result := CompareResult{Added: 2, Removed: 1, Changed: 3, Unchanged: 5, Score: 0.5}
	summary := result.Summary()

	for _, want := range []string{"added=2", "removed=1", "changed=3", "unchanged=5", "50.0%"} {
		if !strings.Contains(summary, want) {
			t.Errorf("summary %q missing %q", summary, want)
		}
	}
}
