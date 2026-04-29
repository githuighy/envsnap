package snapshot

import (
	"testing"
)

var baseLockSnap = Snapshot{
	"APP_ENV":    "production",
	"LOG_LEVEL":  "info",
	"DB_HOST":    "db.example.com",
	"API_KEY":    "secret123",
}

func TestLock_AllMatch(t *testing.T) {
	opts := LockOptions{
		LockedKeys: map[string]string{
			"APP_ENV":   "production",
			"LOG_LEVEL": "info",
		},
	}
	results, err := Lock(baseLockSnap, opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Locked {
			t.Errorf("expected key %q to be locked", r.Key)
		}
	}
}

func TestLock_ValueMismatch(t *testing.T) {
	opts := LockOptions{
		LockedKeys: map[string]string{
			"APP_ENV": "staging",
		},
	}
	results, err := Lock(baseLockSnap, opts)
	if err == nil {
		t.Fatal("expected error for value mismatch")
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Locked {
		t.Error("expected result to be unlocked")
	}
	if results[0].Actual != "production" {
		t.Errorf("unexpected actual value: %q", results[0].Actual)
	}
}

func TestLock_MissingKey(t *testing.T) {
	opts := LockOptions{
		LockedKeys: map[string]string{
			"MISSING_KEY": "value",
		},
	}
	results, err := Lock(baseLockSnap, opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Actual != "" {
		t.Errorf("expected empty actual for missing key")
	}
}

func TestLock_NoLockedKeys(t *testing.T) {
	opts := LockOptions{LockedKeys: map[string]string{}}
	results, err := Lock(baseLockSnap, opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestLock_MultipleViolations(t *testing.T) {
	opts := LockOptions{
		LockedKeys: map[string]string{
			"APP_ENV":   "staging",
			"LOG_LEVEL": "debug",
		},
	}
	results, err := Lock(baseLockSnap, opts)
	if err == nil {
		t.Fatal("expected error for multiple violations")
	}
	locked := 0
	for _, r := range results {
		if r.Locked {
			locked++
		}
	}
	if locked != 0 {
		t.Errorf("expected 0 locked keys, got %d", locked)
	}
}
