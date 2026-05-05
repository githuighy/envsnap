package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envsnap/internal/snapshot"
)

// runNormalize loads a snapshot file, applies normalization, and prints
// or saves the result. Flags mirror the NormalizeOptions fields.
func runNormalize(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap normalize <snapshot> [--no-uppercase] [--no-trim] [--remove-empty] [--no-sanitize] [--format text|json] [--out <file>]")
	}

	path := args[0]
	opts := snapshot.DefaultNormalizeOptions()
	format := "text"
	outPath := ""

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--no-uppercase":
			opts.UppercaseKeys = false
		case "--no-trim":
			opts.TrimValues = false
		case "--remove-empty":
			opts.RemoveEmptyValues = true
		case "--no-sanitize":
			opts.SanitizeKeys = false
		case "--format":
			if i+1 >= len(args) {
				return fmt.Errorf("--format requires a value")
			}
			i++
			format = args[i]
		case "--out":
			if i+1 >= len(args) {
				return fmt.Errorf("--out requires a value")
			}
			i++
			outPath = args[i]
		}
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	normalized := snapshot.Normalize(snap, opts)

	if outPath != "" {
		if err := snapshot.Save(normalized, outPath); err != nil {
			return fmt.Errorf("save normalized snapshot: %w", err)
		}
		fmt.Fprintf(os.Stdout, "normalized snapshot saved to %s\n", outPath)
		return nil
	}

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(normalized)
	default:
		for k, v := range normalized {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
	}
	return nil
}
