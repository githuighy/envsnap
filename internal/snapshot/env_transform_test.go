package snapshot

import (
	"fmt"
	"strings"
	"testing"
)

func baseTransformSnap() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://localhost/db",
		"LOG_LEVEL": "debug",
	}
}

func TestTransform_AllKeys(t *testing.T) {
	snap := baseTransformSnap()
	upper := func(k, v string) (string, error) { return strings.ToUpper(v), nil }
	got, err := Transform(snap, upper, TransformOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range got {
		if v != strings.ToUpper(snap[k]) {
			t.Errorf("key %q: want %q, got %q", k, strings.ToUpper(snap[k]), v)
		}
	}
}

func TestTransform_ByPrefix(t *testing.T) {
	snap := baseTransformSnap()
	upper := func(k, v string) (string, error) { return strings.ToUpper(v), nil }
	got, err := Transform(snap, upper, TransformOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["APP_HOST"] != "LOCALHOST" {
		t.Errorf("APP_HOST: want LOCALHOST, got %q", got["APP_HOST"])
	}
	if got["DB_URL"] != snap["DB_URL"] {
		t.Errorf("DB_URL should be unchanged")
	}
}

func TestTransform_ExplicitKeys(t *testing.T) {
	snap := baseTransformSnap()
	upper := func(k, v string) (string, error) { return strings.ToUpper(v), nil }
	got, err := Transform(snap, upper, TransformOptions{Keys: []string{"LOG_LEVEL"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["LOG_LEVEL"] != "DEBUG" {
		t.Errorf("LOG_LEVEL: want DEBUG, got %q", got["LOG_LEVEL"])
	}
	if got["APP_HOST"] != snap["APP_HOST"] {
		t.Errorf("APP_HOST should be unchanged")
	}
}

func TestTransform_ErrorPropagates(t *testing.T) {
	snap := baseTransformSnap()
	failFn := func(k, v string) (string, error) { return "", fmt.Errorf("boom") }
	_, err := Transform(snap, failFn, TransformOptions{Keys: []string{"APP_HOST"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTransform_SkipErrors(t *testing.T) {
	snap := baseTransformSnap()
	failFn := func(k, v string) (string, error) { return "", fmt.Errorf("boom") }
	got, err := Transform(snap, failFn, TransformOptions{Keys: []string{"APP_HOST"}, SkipErrors: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["APP_HOST"] != snap["APP_HOST"] {
		t.Errorf("APP_HOST should retain original value on skip")
	}
}

func TestTransform_NilFn(t *testing.T) {
	snap := baseTransformSnap()
	_, err := Transform(snap, nil, TransformOptions{})
	if err == nil {
		t.Fatal("expected error for nil fn")
	}
}
