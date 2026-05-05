package snapshot

import (
	"testing"
)

var basePinSnap = Snapshot{
	"APP_ENV":    "production",
	"APP_PORT":   "8080",
	"DB_HOST":    "db.example.com",
	"DB_PASS":    "secret",
	"FEATURE_X":  "true",
}

func TestPin_AllKeys(t *testing.T) {
	result, err := Pin(basePinSnap, PinOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Pinned) != len(basePinSnap) {
		t.Errorf("expected %d pinned keys, got %d", len(basePinSnap), len(result.Pinned))
	}
	for k, v := range basePinSnap {
		if result.Pinned[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, result.Pinned[k])
		}
	}
}

func TestPin_ExplicitKeys(t *testing.T) {
	result, err := Pin(basePinSnap, PinOptions{Keys: []string{"APP_ENV", "APP_PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Pinned) != 2 {
		t.Errorf("expected 2 pinned keys, got %d", len(result.Pinned))
	}
	if result.Pinned["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: expected production, got %q", result.Pinned["APP_ENV"])
	}
}

func TestPin_MissingKey_ReturnsError(t *testing.T) {
	_, err := Pin(basePinSnap, PinOptions{Keys: []string{"DOES_NOT_EXIST"}})
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestPin_AllowMissing_SkipsAbsentKeys(t *testing.T) {
	result, err := Pin(basePinSnap, PinOptions{
		Keys:         []string{"APP_ENV", "MISSING_KEY"},
		AllowMissing: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Pinned) != 1 {
		t.Errorf("expected 1 pinned key, got %d", len(result.Pinned))
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY in skipped, got %v", result.Skipped)
	}
}

func TestPinToSnapshot_RoundTrip(t *testing.T) {
	result, err := Pin(basePinSnap, PinOptions{Keys: []string{"DB_HOST", "APP_PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snap := PinToSnapshot(result)
	if snap["DB_HOST"] != "db.example.com" {
		t.Errorf("DB_HOST: expected db.example.com, got %q", snap["DB_HOST"])
	}
	if snap["APP_PORT"] != "8080" {
		t.Errorf("APP_PORT: expected 8080, got %q", snap["APP_PORT"])
	}
}
