package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

// runReorder implements the `reorder` sub-command.
//
// Usage:
//
//	envsnap reorder <input> [--alpha] [--desc] [--priority DB_,APP_] [--out <file>] [--format json|env]
func runReorder(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: envsnap reorder <snapshot> [--alpha] [--desc] [--priority PREFIX,...] [--out FILE] [--format json|env]")
	}

	var (
		input    = args[0]
		alpha    bool
		desc     bool
		priority []string
		outFile  string
		format   = "env"
	)

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--alpha":
			alpha = true
		case "--desc":
			desc = true
		case "--priority":
			i++
			if i >= len(args) {
				return fmt.Errorf("--priority requires a value")
			}
			for _, p := range strings.Split(args[i], ",") {
				if p = strings.TrimSpace(p); p != "" {
					priority = append(priority, p)
				}
			}
		case "--out":
			i++
			if i >= len(args) {
				return fmt.Errorf("--out requires a value")
			}
			outFile = args[i]
		case "--format":
			i++
			if i >= len(args) {
				return fmt.Errorf("--format requires a value")
			}
			format = args[i]
		}
	}

	snap, err := snapshot.Load(input)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	out, err := snapshot.Reorder(snap, snapshot.ReorderOptions{
		Alphabetical:   alpha,
		Descending:     desc,
		PrefixPriority: priority,
	})
	if err != nil {
		return fmt.Errorf("reorder: %w", err)
	}

	if outFile != "" {
		return snapshot.Save(out, outFile)
	}

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out.Vars)
	default:
		for k, v := range out.Vars {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
	}
	return nil
}
