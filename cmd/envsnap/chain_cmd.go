package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envsnap/internal/snapshot"
)

// runChain resolves a priority-ordered chain of snapshot files and prints the
// merged result. Usage: envsnap chain [--format=text|json|env] file1 file2 ...
func runChain(args []string, format string) error {
	if len(args) < 1 {
		return fmt.Errorf("chain requires at least one snapshot file")
	}

	var layers []snapshot.Snapshot
	for _, path := range args {
		snap, err := snapshot.Load(path)
		if err != nil {
			return fmt.Errorf("loading %s: %w", path, err)
		}
		layers = append(layers, snap)
	}

	chain := snapshot.NewChain(layers...)
	if err := chain.Validate(); err != nil {
		return err
	}
	resolved := chain.Resolve()

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(resolved)
	case "env":
		for k, v := range resolved {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
	case "explain":
		origin := chain.Explain()
		for k, v := range resolved {
			fmt.Fprintf(os.Stdout, "%s=%s  (layer %d: %s)\n", k, v, origin[k], args[origin[k]])
		}
	default:
		for k, v := range resolved {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
	}
	return nil
}
