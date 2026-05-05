package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/your-org/envsnap/internal/snapshot"
)

// runClone implements: envsnap clone <src> <dst> [flags]
//
//	--prefix        only include keys with this prefix (repeatable)
//	--exclude       exclude keys with this prefix (repeatable)
//	--rename        old=new key rename (repeatable)
//	--set           key=value override applied after cloning (repeatable)
//	--format        output format for confirmation: text|json (default: text)
func runClone(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envsnap clone <src> <dst> [--prefix P] [--exclude E] [--rename OLD=NEW] [--set K=V]")
	}
	srcPath, dstPath := args[0], args[1]
	rest := args[2:]

	var prefixes, excludes, renames, sets []string
	format := "text"

	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case "--prefix":
			i++; prefixes = append(prefixes, rest[i])
		case "--exclude":
			i++; excludes = append(excludes, rest[i])
		case "--rename":
			i++; renames = append(renames, rest[i])
		case "--set":
			i++; sets = append(sets, rest[i])
		case "--format":
			i++; format = rest[i]
		}
	}

	src, err := snapshot.Load(srcPath)
	if err != nil {
		return fmt.Errorf("clone: load source: %w", err)
	}

	keyMap := make(map[string]string)
	for _, r := range renames {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("clone: invalid --rename %q (want OLD=NEW)", r)
		}
		keyMap[parts[0]] = parts[1]
	}

	overrides := make(map[string]string)
	for _, s := range sets {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("clone: invalid --set %q (want KEY=VALUE)", s)
		}
		overrides[parts[0]] = parts[1]
	}

	out, err := snapshot.Clone(src, snapshot.CloneOptions{
		Prefixes:       prefixes,
		Exclude:        excludes,
		KeyMap:         keyMap,
		OverrideValues: overrides,
	})
	if err != nil {
		return fmt.Errorf("clone: %w", err)
	}

	if err := snapshot.Save(dstPath, out); err != nil {
		return fmt.Errorf("clone: save: %w", err)
	}

	switch format {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"src": srcPath, "dst": dstPath, "keys": len(out),
		})
	default:
		fmt.Printf("cloned %d key(s) from %s → %s\n", len(out), srcPath, dstPath)
	}
	return nil
}
