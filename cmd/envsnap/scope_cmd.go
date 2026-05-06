package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/user/envsnap/internal/snapshot"
)

func runScope(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envsnap scope <snapshot.json> <scope> [--sep SEP] [--strip] [--unscope] [--format text|json] [--out FILE]")
	}

	snapPath := args[0]
	scopeName := args[1]

	sep := "_"
	strip := false
	unscope := false
	format := "text"
	outPath := ""

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "--sep":
			i++
			if i < len(args) {
				sep = args[i]
			}
		case "--strip":
			strip = true
		case "--unscope":
			unscope = true
		case "--format":
			i++
			if i < len(args) {
				format = args[i]
			}
		case "--out":
			i++
			if i < len(args) {
				outPath = args[i]
			}
		}
	}

	snap, err := snapshot.Load(snapPath)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	var result map[string]string
	if unscope {
		result = snapshot.Unscope(snap, scopeName, sep)
	} else {
		res, err := snapshot.Scope(snap, snapshot.ScopeOptions{
			Scope:           scopeName,
			PrefixSeparator: sep,
			StripExisting:   strip,
		})
		if err != nil {
			return fmt.Errorf("scoping snapshot: %w", err)
		}
		result = res.Snapshot
	}

	if outPath != "" {
		return snapshot.Save(result, outPath)
	}

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	default:
		keys := make([]string, 0, len(result))
		for k := range result {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("%s=%s\n", k, result[k])
		}
	}
	return nil
}
