package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envsnap/internal/snapshot"
)

func runEncrypt(args []string, mode string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap %s <snapshot> [--pass <passphrase>] [--keys k1,k2] [--prefix p]", mode)
	}

	file := args[0]
	passphrase := flagString(args, "--pass", "")
	if passphrase == "" {
		passphrase = os.Getenv("ENVSNAP_PASSPHRASE")
	}
	if passphrase == "" {
		return fmt.Errorf("%s: passphrase required (--pass or ENVSNAP_PASSPHRASE)", mode)
	}

	snap, err := snapshot.Load(file)
	if err != nil {
		return fmt.Errorf("%s: load %q: %w", mode, file, err)
	}

	outFile := flagString(args, "--out", file)
	format := flagString(args, "--format", "text")

	switch mode {
	case "encrypt":
		opts := snapshot.EncryptOptions{
			Passphrase: passphrase,
			Keys:       splitCSV(flagString(args, "--keys", "")),
			Prefixes:   splitCSV(flagString(args, "--prefix", "")),
		}
		result, err := snapshot.Encrypt(snap, opts)
		if err != nil {
			return fmt.Errorf("encrypt: %w", err)
		}
		if err := snapshot.Save(result, outFile); err != nil {
			return fmt.Errorf("encrypt: save: %w", err)
		}
		printEncryptResult(result, format, "encrypted")

	case "decrypt":
		result, err := snapshot.Decrypt(snap, passphrase)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
		if err := snapshot.Save(result, outFile); err != nil {
			return fmt.Errorf("decrypt: save: %w", err)
		}
		printEncryptResult(result, format, "decrypted")
	}
	return nil
}

func printEncryptResult(snap snapshot.Snapshot, format, action string) {
	switch format {
	case "json":
		payload := map[string]interface{}{
			"action": action,
			"count":  len(snap),
			"keys":   sortedKeys(snap),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(payload)
	default:
		fmt.Printf("%s: %d keys written\n", action, len(snap))
	}
}

// sortedKeys returns snapshot keys in sorted order.
func sortedKeys(snap snapshot.Snapshot) []string {
	keys := make([]string, 0, len(snap))
	for k := range snap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
