package snapshot

import (
	"testing"
)

func baseSnaps() (base, override map[string]string) {
	base = map[string]string{
		"APP_ENV":  "production",
		"LOG_LEVEL": "warn",
		"DB_HOST":  "db.prod.internal",
	}
	override = map[string]string{
		"APP_ENV":   "staging",
		"LOG_LEVEL": "debug",
		"NEW_KEY":   "hello",
	}
	return
}

func TestMerge_OverrideWins(t *testing.T) {
	base, override := baseSnaps()
	result := Merge(base, override, MergeOptions{Prefer: "override"})

	if result["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %s", result["APP_ENV"])
	}
	if result["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %s", result["LOG_LEVEL"])
	}
}

func TestMerge_BaseWins(t *testing.T) {
	base, override := baseSnaps()
	result := Merge(base, override, MergeOptions{Prefer: "base"})

	if result["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %s", result["APP_ENV"])
	}
	if result["LOG_LEVEL"] != "warn" {
		t.Errorf("expected LOG_LEVEL=warn, got %s", result["LOG_LEVEL"])
	}
}

func TestMerge_UniqueKeysAlwaysIncluded(t *testing.T) {
	base, override := baseSnaps()
	result := Merge(base, override, MergeOptions{})

	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST from base to be present")
	}
	if _, ok := result["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY from override to be present")
	}
}

func TestMerge_DefaultPrefersOverride(t *testing.T) {
	base, override := baseSnaps()
	result := Merge(base, override, MergeOptions{})

	if result["APP_ENV"] != "staging" {
		t.Errorf("default should prefer override; got APP_ENV=%s", result["APP_ENV"])
	}
}

func TestMerge_EmptyBase(t *testing.T) {
	override := map[string]string{"ONLY_KEY": "value"}
	result := Merge(map[string]string{}, override, MergeOptions{})

	if result["ONLY_KEY"] != "value" {
		t.Errorf("expected ONLY_KEY=value, got %s", result["ONLY_KEY"])
	}
}

func TestMerge_EmptyOverride(t *testing.T) {
	base := map[string]string{"BASE_KEY": "base_val"}
	result := Merge(base, map[string]string{}, MergeOptions{})

	if result["BASE_KEY"] != "base_val" {
		t.Errorf("expected BASE_KEY=base_val, got %s", result["BASE_KEY"])
	}
}
