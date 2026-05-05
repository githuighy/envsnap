package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

// runRename loads a snapshot, applies key renaming, and prints or saves the result.
//
// Usage:
//
//	envsnap rename <snapshot> [--map OLD=NEW,...] [--strip-prefix PREFIX]
//	              [--add-prefix PREFIX] [--format env|json] [--out FILE]
//	              [--fail-on-conflict]
func runRename(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap rename <snapshot> [options]")
	}

	snapPath := args[0]
	snap, err := snapshot.Load(snapPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	opts := snapshot.RenameOptions{
		Map: make(map[string]string),
	}

	format := "env"
	outPath := ""

	for i := 1; i < len(args); i++ {
		switch {
		case args[i] == "--map" && i+1 < len(args):
			i++
			for _, pair := range strings.Split(args[i], ",") {
				parts := strings.SplitN(pair, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid --map entry %q, expected OLD=NEW", pair)
				}
				opts.Map[parts[0]] = parts[1]
			}
		case args[i] == "--strip-prefix" && i+1 < len(args):
			i++
			opts.StripPrefix = args[i]
		case args[i] == "--add-prefix" && i+1 < len(args):
			i++
			opts.AddPrefix = args[i]
		case args[i] == "--fail-on-conflict":
			opts.FailOnConflict = true
		case args[i] == "--format" && i+1 < len(args):
			i++
			format = args[i]
		case args[i] == "--out" && i+1 < len(args):
			i++
			outPath = args[i]
		}
	}

	res, err := snapshot.Rename(snap, opts)
	if err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	if len(res.Conflict) > 0 {
		fmt.Fprintf(os.Stderr, "warn: skipped %d conflict(s): %s\n",
			len(res.Conflict), strings.Join(res.Conflict, ", "))
	}

	if outPath != "" {
		return snapshot.Save(res.Snap, outPath)
	}

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(res.Snap)
	case "env":
		for _, k := range sortedKeys(res.Snap) {
			fmt.Printf("%s=%s\n", k, res.Snap[k])
		}
		return nil
	default:
		return fmt.Errorf("unknown format %q, use env or json", format)
	}
}
