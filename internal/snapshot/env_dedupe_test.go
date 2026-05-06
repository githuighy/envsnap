package snapshot

import (
	"testing"
)

var baseDedupeSnap = map[string]string{
	"APP_HOST":    "localhost",
	"APP_ADDRESS": "localhost", // duplicate value of APP_HOST
	"DB_HOST":     "db.local",
	"DB_PORT":     "5432",
	"DB_REPLICA":  "5432", // duplicate value of DB_PORT
	"LOG_LEVEL":   "info",
}

func TestDedupe_ByValue_RemovesDuplicateValues(t *testing.T) {
	result := Dedupe(baseDedupeSnap, DedupeOptions{ByValue: true})

	// APP_ADDRESS duplicates APP_HOST; DB_REPLICA duplicates DB_PORT.
	for _, removed := range result.RemovedKeys {
		if removed != "APP_ADDRESS" && removed != "DB_REPLICA" {
			t.Errorf("unexpected removed key: %s", removed)
		}
	}
	if len(result.RemovedKeys) != 2 {
		t.Errorf("expected 2 removed keys, got %d: %v", len(result.RemovedKeys), result.RemovedKeys)
	}
	if _, exists := result.Snapshot["APP_HOST"]; !exists {
		t.Error("keeper key APP_HOST should remain")
	}
	if _, exists := result.Snapshot["DB_PORT"]; !exists {
		t.Error("keeper key DB_PORT should remain")
	}
}

func TestDedupe_ByPrefix_KeepsFirstKey(t *testing.T) {
	result := Dedupe(baseDedupeSnap, DedupeOptions{
		ByPrefix:  true,
		Prefixes:  []string{"APP_", "DB_"},
	})

	// For APP_ prefix: APP_ADDRESS and APP_HOST — only first (sorted) kept.
	// For DB_ prefix: DB_HOST, DB_PORT, DB_REPLICA — only DB_HOST kept.
	if len(result.RemovedKeys) != 3 {
		t.Errorf("expected 3 removed keys, got %d: %v", len(result.RemovedKeys), result.RemovedKeys)
	}
	if _, ok := result.Snapshot["LOG_LEVEL"]; !ok {
		t.Error("LOG_LEVEL should not be removed (no matching prefix)")
	}
}

func TestDedupe_ByValueAndPrefix_Combined(t *testing.T) {
	snap := map[string]string{
		"X_FOO": "same",
		"X_BAR": "same",
		"Y_BAZ": "other",
	}
	result := Dedupe(snap, DedupeOptions{
		ByValue:  true,
		ByPrefix: true,
		Prefixes: []string{"X_"},
	})
	// ByValue removes X_BAR (duplicate of X_FOO), ByPrefix would also remove it.
	if _, ok := result.Snapshot["X_FOO"]; !ok {
		t.Error("X_FOO should remain")
	}
	if _, ok := result.Snapshot["X_BAR"]; ok {
		t.Error("X_BAR should be removed")
	}
	if _, ok := result.Snapshot["Y_BAZ"]; !ok {
		t.Error("Y_BAZ should remain")
	}
}

func TestDedupe_NoOptions_ReturnsUnchanged(t *testing.T) {
	result := Dedupe(baseDedupeSnap, DedupeOptions{})
	if len(result.RemovedKeys) != 0 {
		t.Errorf("expected no removed keys, got %v", result.RemovedKeys)
	}
	if len(result.Snapshot) != len(baseDedupeSnap) {
		t.Errorf("snapshot size mismatch: got %d, want %d", len(result.Snapshot), len(baseDedupeSnap))
	}
}

func TestDedupe_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"A": "v", "B": "v"}
	origLen := len(input)
	Dedupe(input, DedupeOptions{ByValue: true})
	if len(input) != origLen {
		t.Error("Dedupe must not mutate the input snapshot")
	}
}
