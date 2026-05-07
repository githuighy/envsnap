package snapshot

import (
	"testing"
)

func baseSplitSnap() Snapshot {
	return Snapshot{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"CACHE_HOST":  "redis",
		"CACHE_TTL":   "300",
		"APP_VERSION": "1.0.0",
		"LOG_LEVEL":   "info",
	}
}

func TestSplit_BasicBuckets(t *testing.T) {
	result, err := Split(baseSplitSnap(), SplitOptions{
		Buckets: map[string][]string{
			"database": {"DB_"},
			"cache":    {"CACHE_"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["database"]["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST in database bucket")
	}
	if result["cache"]["CACHE_TTL"] != "300" {
		t.Errorf("expected CACHE_TTL in cache bucket")
	}
}

func TestSplit_RemainderBucket(t *testing.T) {
	result, err := Split(baseSplitSnap(), SplitOptions{
		Buckets: map[string][]string{
			"database": {"DB_"},
		},
		Remainder: "other",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["other"]["APP_VERSION"]; !ok {
		t.Errorf("expected APP_VERSION in remainder bucket")
	}
	if _, ok := result["other"]["DB_HOST"]; ok {
		t.Errorf("DB_HOST should not appear in remainder bucket")
	}
}

func TestSplit_NoBuckets_ReturnsError(t *testing.T) {
	_, err := Split(baseSplitSnap(), SplitOptions{})
	if err == nil {
		t.Fatal("expected error when no buckets defined")
	}
}

func TestSplit_UnmatchedKeysDiscardedWithoutRemainder(t *testing.T) {
	result, err := Split(baseSplitSnap(), SplitOptions{
		Buckets: map[string][]string{
			"database": {"DB_"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["other"]; ok {
		t.Errorf("no remainder bucket should exist")
	}
	if len(result["database"]) != 2 {
		t.Errorf("expected 2 keys in database bucket, got %d", len(result["database"]))
	}
}

func TestSplit_MultiplePrefixesPerBucket(t *testing.T) {
	result, err := Split(baseSplitSnap(), SplitOptions{
		Buckets: map[string][]string{
			"infra": {"DB_", "CACHE_"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result["infra"]) != 4 {
		t.Errorf("expected 4 keys in infra bucket, got %d", len(result["infra"]))
	}
}
