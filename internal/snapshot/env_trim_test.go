package snapshot

import (
	"testing"
)

func baseTrimSnap() map[string]string {
	return map[string]string{
		"APP_HOST":    "  localhost  ",
		"APP_PORT":    "  8080  ",
		"DB_URL":      "postgres://localhost/mydb",
		"DB_PASSWORD": "  secret  ",
		"PLAIN":       "  value  ",
	}
}

func TestTrim_TrimSpace_AllKeys(t *testing.T) {
	snap := baseTrimSnap()
	out, err := Trim(snap, TrimOptions{TrimSpace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", out["APP_HOST"])
	}
	if out["APP_PORT"] != "8080" {
		t.Errorf("expected '8080', got %q", out["APP_PORT"])
	}
	if out["PLAIN"] != "value" {
		t.Errorf("expected 'value', got %q", out["PLAIN"])
	}
}

func TestTrim_TrimSpace_ByPrefix(t *testing.T) {
	snap := baseTrimSnap()
	out, err := Trim(snap, TrimOptions{
		Prefixes: []string{"APP_"},
		TrimSpace: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", out["APP_HOST"])
	}
	// DB_PASSWORD should be untouched
	if out["DB_PASSWORD"] != "  secret  " {
		t.Errorf("expected '  secret  ', got %q", out["DB_PASSWORD"])
	}
}

func TestTrim_TrimPrefix_ExplicitKeys(t *testing.T) {
	snap := map[string]string{
		"DB_URL": "postgres://localhost/mydb",
		"OTHER":  "postgres://other",
	}
	out, err := Trim(snap, TrimOptions{
		Keys:       []string{"DB_URL"},
		TrimPrefix: "postgres://",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_URL"] != "localhost/mydb" {
		t.Errorf("expected 'localhost/mydb', got %q", out["DB_URL"])
	}
	if out["OTHER"] != "postgres://other" {
		t.Errorf("OTHER should be unchanged, got %q", out["OTHER"])
	}
}

func TestTrim_TrimSuffix(t *testing.T) {
	snap := map[string]string{
		"APP_NAME": "myapp_v1",
		"APP_ENV":  "production_v1",
	}
	out, err := Trim(snap, TrimOptions{
		TrimSuffix: "_v1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected 'myapp', got %q", out["APP_NAME"])
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected 'production', got %q", out["APP_ENV"])
	}
}

func TestTrim_DoesNotMutateOriginal(t *testing.T) {
	snap := map[string]string{"KEY": "  val  "}
	_, err := Trim(snap, TrimOptions{TrimSpace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap["KEY"] != "  val  " {
		t.Errorf("original snapshot was mutated")
	}
}
