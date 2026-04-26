package snapshot

import (
	"testing"
)

func baseRedactSnap() map[string]string {
	return map[string]string{
		"HOME":           "/home/user",
		"DB_PASSWORD":    "s3cr3t",
		"API_KEY":        "abc123",
		"GITHUB_TOKEN":   "ghp_xxxx",
		"AWS_ACCESS_KEY": "AKIA...",
		"APP_ENV":        "production",
		"MY_SECRET":      "topsecret",
	}
}

func TestRedact_DefaultKeys(t *testing.T) {
	snap := baseRedactSnap()
	result := Redact(snap, RedactOptions{})

	sensitiveVars := []string{"DB_PASSWORD", "API_KEY", "GITHUB_TOKEN", "AWS_ACCESS_KEY", "MY_SECRET"}
	for _, key := range sensitiveVars {
		if result[key] != "[REDACTED]" {
			t.Errorf("expected %s to be redacted, got %q", key, result[key])
		}
	}
}

func TestRedact_SafeKeysUnchanged(t *testing.T) {
	snap := baseRedactSnap()
	result := Redact(snap, RedactOptions{})

	if result["HOME"] != "/home/user" {
		t.Errorf("expected HOME to be unchanged, got %q", result["HOME"])
	}
	if result["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV to be unchanged, got %q", result["APP_ENV"])
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	snap := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"HOME":        "/home/user",
	}
	result := Redact(snap, RedactOptions{Placeholder: "***"})

	if result["DB_PASSWORD"] != "***" {
		t.Errorf("expected *** placeholder, got %q", result["DB_PASSWORD"])
	}
	if result["HOME"] != "/home/user" {
		t.Errorf("expected HOME unchanged, got %q", result["HOME"])
	}
}

func TestRedact_CustomSensitiveKeys(t *testing.T) {
	snap := map[string]string{
		"MY_CUSTOM_CERT": "certdata",
		"NORMAL_VAR":     "value",
	}
	result := Redact(snap, RedactOptions{
		SensitiveKeys: []string{"CERT"},
	})

	if result["MY_CUSTOM_CERT"] != "[REDACTED]" {
		t.Errorf("expected MY_CUSTOM_CERT to be redacted, got %q", result["MY_CUSTOM_CERT"])
	}
	if result["NORMAL_VAR"] != "value" {
		t.Errorf("expected NORMAL_VAR unchanged, got %q", result["NORMAL_VAR"])
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	snap := map[string]string{
		"DB_PASSWORD": "s3cr3t",
	}
	_ = Redact(snap, RedactOptions{})

	if snap["DB_PASSWORD"] != "s3cr3t" {
		t.Error("original snapshot was mutated by Redact")
	}
}
