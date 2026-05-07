package snapshot

import (
	"sort"
	"testing"
)

func sortedStrings(s []string) []string {
	out := make([]string, len(s))
	copy(out, s)
	sort.Strings(out)
	return out
}

var leftSnap = map[string]string{
	"APP_HOST": "localhost",
	"APP_PORT": "8080",
	"DB_URL":   "postgres://left",
	"SECRET":   "abc",
}

var rightSnap = map[string]string{
	"APP_HOST": "prod.example.com",
	"APP_PORT": "443",
	"DB_URL":   "postgres://right",
	"NEW_KEY":  "only-in-right",
}

func TestIntersect_OnlyCommonKeys(t *testing.T) {
	res := Intersect(leftSnap, rightSnap, IntersectOptions{})
	if _, ok := res.Snapshot["SECRET"]; ok {
		t.Error("expected SECRET to be absent (not in right)")
	}
	if _, ok := res.Snapshot["NEW_KEY"]; ok {
		t.Error("expected NEW_KEY to be absent (not in left)")
	}
	if len(res.Snapshot) != 3 {
		t.Errorf("expected 3 common keys, got %d", len(res.Snapshot))
	}
}

func TestIntersect_PreferRight_DefaultBehaviour(t *testing.T) {
	res := Intersect(leftSnap, rightSnap, IntersectOptions{})
	if got := res.Snapshot["APP_HOST"]; got != "prod.example.com" {
		t.Errorf("expected right value, got %q", got)
	}
}

func TestIntersect_PreferLeft(t *testing.T) {
	res := Intersect(leftSnap, rightSnap, IntersectOptions{PreferLeft: true})
	if got := res.Snapshot["APP_HOST"]; got != "localhost" {
		t.Errorf("expected left value, got %q", got)
	}
}

func TestIntersect_PrefixFilter(t *testing.T) {
	res := Intersect(leftSnap, rightSnap, IntersectOptions{Prefixes: []string{"APP_"}})
	keys := sortedStrings(res.CommonKeys)
	if len(keys) != 2 || keys[0] != "APP_HOST" || keys[1] != "APP_PORT" {
		t.Errorf("unexpected keys: %v", keys)
	}
	if _, ok := res.Snapshot["DB_URL"]; ok {
		t.Error("DB_URL should be excluded by prefix filter")
	}
}

func TestIntersect_ExcludeKeys(t *testing.T) {
	res := Intersect(leftSnap, rightSnap, IntersectOptions{ExcludeKeys: []string{"DB_URL"}})
	if _, ok := res.Snapshot["DB_URL"]; ok {
		t.Error("DB_URL should be excluded")
	}
	if len(res.Snapshot) != 2 {
		t.Errorf("expected 2 keys after exclusion, got %d", len(res.Snapshot))
	}
}

func TestIntersect_EmptySnapshots(t *testing.T) {
	res := Intersect(map[string]string{}, map[string]string{}, IntersectOptions{})
	if len(res.Snapshot) != 0 {
		t.Error("expected empty result for empty inputs")
	}
	if len(res.CommonKeys) != 0 {
		t.Error("expected no common keys")
	}
}
