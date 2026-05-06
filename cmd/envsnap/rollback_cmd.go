package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"envsnap/internal/snapshot"
)

func runRollback(args []string, flags map[string]string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envsnap rollback <current> <baseline> [--output file] [--keys K1,K2] [--prefix P] [--dry-run] [--format text|json]")
	}

	current, err := snapshot.Load(args[0])
	if err != nil {
		return fmt.Errorf("load current: %w", err)
	}
	baseline, err := snapshot.Load(args[1])
	if err != nil {
		return fmt.Errorf("load baseline: %w", err)
	}

	opts := snapshot.RollbackOptions{}
	if k, ok := flags["keys"]; ok && k != "" {
		for _, part := range strings.Split(k, ",") {
			if p := strings.TrimSpace(part); p != "" {
				opts.Keys = append(opts.Keys, p)
			}
		}
	}
	if p, ok := flags["prefix"]; ok && p != "" {
		for _, part := range strings.Split(p, ",") {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				opts.Prefixes = append(opts.Prefixes, trimmed)
			}
		}
	}
	if _, ok := flags["dry-run"]; ok {
		opts.DryRun = true
	}

	res, err := snapshot.Rollback(current, baseline, opts)
	if err != nil {
		return err
	}

	format := flags["format"]
	if format == "json" {
		return printRollbackJSON(res)
	}
	printRollbackText(res, opts.DryRun)

	if !opts.DryRun {
		outPath := flags["output"]
		if outPath == "" {
			outPath = args[0]
		}
		if err := snapshot.Save(res.Snapshot, outPath); err != nil {
			return fmt.Errorf("save: %w", err)
		}
	}
	return nil
}

func printRollbackText(res snapshot.RollbackResult, dryRun bool) {
	prefix := ""
	if dryRun {
		prefix = "[dry-run] "
	}
	keys := make([]string, 0, len(res.Restored))
	for k := range res.Restored {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(os.Stdout, "%srestored %s = %s\n", prefix, k, res.Restored[k])
	}
	sort.Strings(res.Dropped)
	for _, k := range res.Dropped {
		fmt.Fprintf(os.Stdout, "%sdropped  %s\n", prefix, k)
	}
	if len(res.Restored) == 0 && len(res.Dropped) == 0 {
		fmt.Println("nothing to roll back")
	}
}

func printRollbackJSON(res snapshot.RollbackResult) error {
	return json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
		"restored":  res.Restored,
		"dropped":   res.Dropped,
		"unchanged": res.Unchanged,
		"snapshot":  res.Snapshot.Vars,
	})
}
