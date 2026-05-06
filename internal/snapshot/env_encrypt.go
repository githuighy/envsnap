package snapshot

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

// EncryptOptions controls which keys are encrypted and how.
type EncryptOptions struct {
	// Keys is an explicit list of keys to encrypt. If empty, all keys are encrypted.
	Keys []string
	// Prefixes restricts encryption to keys matching any of these prefixes.
	Prefixes []string
	// Passphrase is used to derive the AES-256 encryption key.
	Passphrase string
}

const encryptedPrefix = "enc:"

// Encrypt encrypts values in snap according to opts. Already-encrypted values
// (prefixed with "enc:") are left untouched.
func Encrypt(snap Snapshot, opts EncryptOptions) (Snapshot, error) {
	if opts.Passphrase == "" {
		return nil, errors.New("encrypt: passphrase must not be empty")
	}
	block, err := newAESBlock(opts.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}
	out := make(Snapshot, len(snap))
	for k, v := range snap {
		if !shouldEncrypt(k, opts) || strings.HasPrefix(v, encryptedPrefix) {
			out[k] = v
			continue
		}
		enc, err := aesCTREncrypt(block, v)
		if err != nil {
			return nil, fmt.Errorf("encrypt: key %q: %w", k, err)
		}
		out[k] = encryptedPrefix + enc
	}
	return out, nil
}

// Decrypt decrypts values in snap that carry the "enc:" prefix.
func Decrypt(snap Snapshot, passphrase string) (Snapshot, error) {
	if passphrase == "" {
		return nil, errors.New("decrypt: passphrase must not be empty")
	}
	block, err := newAESBlock(passphrase)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}
	out := make(Snapshot, len(snap))
	for k, v := range snap {
		if !strings.HasPrefix(v, encryptedPrefix) {
			out[k] = v
			continue
		}
		plain, err := aesCTRDecrypt(block, strings.TrimPrefix(v, encryptedPrefix))
		if err != nil {
			return nil, fmt.Errorf("decrypt: key %q: %w", k, err)
		}
		out[k] = plain
	}
	return out, nil
}

func shouldEncrypt(key string, opts EncryptOptions) bool {
	if len(opts.Keys) > 0 {
		for _, k := range opts.Keys {
			if k == key {
				return true
			}
		}
		return false
	}
	if len(opts.Prefixes) > 0 {
		return hasAnyPrefix(key, opts.Prefixes)
	}
	return true
}

func newAESBlock(passphrase string) (cipher.Block, error) {
	hash := sha256.Sum256([]byte(passphrase))
	return aes.NewCipher(hash[:])
}

func aesCTREncrypt(block cipher.Block, plaintext string) (string, error) {
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func aesCTRDecrypt(block cipher.Block, encoded string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])
	return string(plaintext), nil
}
