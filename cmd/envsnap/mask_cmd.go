package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

// runMask loads a snapshot file and prints it with selected values masked.
// Usage: envsnap mask <snapshot> [--keys A,B] [--prefixes P,Q]
//
//	[--visible-chars N] [--show-length] [--placeholder STR]
//	[--format text|json]
func runMask(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap mask <snapshot> [options]")
	}

	snapshotPath := args[0]
	rest := args[1:]

	var (
		rawKeys       string
		rawPrefixes   string
		placeholder   string
		visibleChars  int
		showLength    bool
		format        = "text"
	)

	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case "--keys":
			i++
			if i < len(rest) {
				rawKeys = rest[i]
			}
		case "--prefixes":
			i++
			if i < len(rest) {
				rawPrefixes = rest[i]
			}
		case "--placeholder":
			i++
			if i < len(rest) {
				placeholder = rest[i]
			}
		case "--visible-chars":
			i++
			if i < len(rest) {
				fmt.Sscanf(rest[i], "%d", &visibleChars)
			}
		case "--show-length":
			showLength = true
		case "--format":
			i++
			if i < len(rest) {
				format = rest[i]
			}
		}
	}

	snap, err := snapshot.Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	opts := snapshot.MaskOptions{
		Placeholder:  placeholder,
		ShowLength:   showLength,
		VisibleChars: visibleChars,
	}
	if rawKeys != "" {
		opts.Keys = strings.Split(rawKeys, ",")
	}
	if rawPrefixes != "" {
		opts.Prefixes = strings.Split(rawPrefixes, ",")
	}

	masked := snapshot.Mask(snap, opts)

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(masked)
	default:
		for k, v := range masked {
			fmt.Printf("%s=%s\n", k, v)
		}
	}
	return nil
}
