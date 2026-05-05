package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envsnap/internal/snapshot"
)

// runResolve loads a snapshot file, expands all ${VAR} / $VAR references
// within its values, and prints or saves the result.
//
// Usage:
//
//	envsnap resolve <snapshot.json> [--allow-missing] [--format env|json] [--out <file>]
func runResolve(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap resolve <snapshot.json> [--allow-missing] [--format env|json] [--out <file>]")
	}

	snapPath := args[0]
	allowMissing := false
	format := "env"
	outPath := ""

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--allow-missing":
			allowMissing = true
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

	snap, err := snapshot.Load(snapPath)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	opts := snapshot.DefaultResolveOptions()
	opts.AllowMissing = allowMissing

	resolved, err := snapshot.Resolve(snap, opts)
	if err != nil {
		return fmt.Errorf("resolving references: %w", err)
	}

	var output []byte
	switch format {
	case "json":
		output, err = json.MarshalIndent(resolved, "", "  ")
		if err != nil {
			return fmt.Errorf("marshalling JSON: %w", err)
		}
		output = append(output, '\n')
	case "env":
		var buf []byte
		for k, v := range resolved {
			buf = append(buf, fmt.Sprintf("%s=%s\n", k, v)...)
		}
		output = buf
	default:
		return fmt.Errorf("unknown format %q: use env or json", format)
	}

	if outPath != "" {
		if err := os.WriteFile(outPath, output, 0o644); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
		fmt.Fprintf(os.Stderr, "resolved snapshot written to %s\n", outPath)
		return nil
	}

	_, err = os.Stdout.Write(output)
	return err
}
