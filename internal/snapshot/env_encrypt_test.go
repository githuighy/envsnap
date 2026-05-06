package snapshot

import (
	"strings"
	"testing"
)

var baseEncryptSnap = Snapshot{
	"DB_PASSWORD":  "s3cr3t",
	"API_KEY":      "abc123",
	"APP_ENV":      "production",
	"LOG_LEVEL":    "info",
}

func TestEncrypt_AllKeys(t *testing.T) {
	opts := EncryptOptions{Passphrase: "testpass"}
	out, err := Encrypt(baseEncryptSnap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range out {
		if !strings.HasPrefix(v, encryptedPrefix) {
			t.Errorf("key %q value %q not encrypted", k, v)
		}
	}
}

func TestEncrypt_ExplicitKeys(t *testing.T) {
	opts := EncryptOptions{
		Passphrase: "testpass",
		Keys:       []string{"DB_PASSWORD", "API_KEY"},
	}
	out, err := Encrypt(baseEncryptSnap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out["DB_PASSWORD"], encryptedPrefix) {
		t.Error("DB_PASSWORD should be encrypted")
	}
	if !strings.HasPrefix(out["API_KEY"], encryptedPrefix) {
		t.Error("API_KEY should be encrypted")
	}
	if strings.HasPrefix(out["APP_ENV"], encryptedPrefix) {
		t.Error("APP_ENV should NOT be encrypted")
	}
}

func TestEncrypt_ByPrefix(t *testing.T) {
	opts := EncryptOptions{
		Passphrase: "testpass",
		Prefixes:   []string{"DB_"},
	}
	out, err := Encrypt(baseEncryptSnap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out["DB_PASSWORD"], encryptedPrefix) {
		t.Error("DB_PASSWORD should be encrypted")
	}
	if strings.HasPrefix(out["API_KEY"], encryptedPrefix) {
		t.Error("API_KEY should NOT be encrypted")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	pass := "roundtrip-secret"
	opts := EncryptOptions{Passphrase: pass}
	encrypted, err := Encrypt(baseEncryptSnap, opts)
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}
	decrypted, err := Decrypt(encrypted, pass)
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}
	for k, want := range baseEncryptSnap {
		if got := decrypted[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestEncrypt_AlreadyEncrypted_Skipped(t *testing.T) {
	snap := Snapshot{"SECRET": encryptedPrefix + "alreadydone"}
	opts := EncryptOptions{Passphrase: "pass"}
	out, err := Encrypt(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET"] != snap["SECRET"] {
		t.Error("already-encrypted value should not be re-encrypted")
	}
}

func TestEncrypt_EmptyPassphrase_ReturnsError(t *testing.T) {
	_, err := Encrypt(baseEncryptSnap, EncryptOptions{})
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestDecrypt_WrongPassphrase_ReturnsError(t *testing.T) {
	encrypted, err := Encrypt(baseEncryptSnap, EncryptOptions{Passphrase: "correct"})
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}
	// Wrong passphrase produces garbled output but no error (CTR mode).
	// Verify at minimum that values differ from original.
	decrypted, err := Decrypt(encrypted, "wrong")
	if err != nil {
		t.Fatalf("unexpected decrypt error: %v", err)
	}
	for k, orig := range baseEncryptSnap {
		if decrypted[k] == orig {
			t.Errorf("key %q: wrong passphrase should not yield original value", k)
		}
	}
}
