package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

func runPromote(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envsnap promote <src> <dst> [--overwrite] [--dry-run] [--prefix=X] [--exclude=X] [--format=text|json]")
	}

	srcPath := args[0]
	dstPath := args[1]

	var prefixes, excludes []string
	overwrite := false
	dryRun := false
	format := "text"

	for _, arg := range args[2:] {
		switch {
		case arg == "--overwrite":
			overwrite = true
		case arg == "--dry-run":
			dryRun = true
		case strings.HasPrefix(arg, "--prefix="):
			prefixes = append(prefixes, strings.TrimPrefix(arg, "--prefix="))
		case strings.HasPrefix(arg, "--exclude="):
			excludes = append(excludes, strings.TrimPrefix(arg, "--exclude="))
		case strings.HasPrefix(arg, "--format="):
			format = strings.TrimPrefix(arg, "--format=")
		}
	}

	src, err := snapshot.Load(srcPath)
	if err != nil {
		return fmt.Errorf("loading src: %w", err)
	}
	dst, err := snapshot.Load(dstPath)
	if err != nil {
		return fmt.Errorf("loading dst: %w", err)
	}

	results, out, err := snapshot.Promote(src, dst, snapshot.PromoteOptions{
		Prefixes:  prefixes,
		Exclude:   excludes,
		Overwrite: overwrite,
		DryRun:    dryRun,
	})
	if err != nil {
		return err
	}

	if format == "json" {
		return printPromoteJSON(results, dryRun)
	}

	for _, r := range results {
		switch r.Action {
		case "promoted":
			fmt.Printf("  + %s\n", r.Key)
		case "skipped_exists":
			fmt.Printf("  ~ %s (skipped, exists)\n", r.Key)
		}
	}
	fmt.Println(snapshot.PromoteSummary(results))

	if !dryRun {
		if err := snapshot.Save(dstPath, out); err != nil {
			return fmt.Errorf("saving dst: %w", err)
		}
	}
	return nil
}

func printPromoteJSON(results []snapshot.PromoteResult, dryRun bool) error {
	payload := map[string]interface{}{
		"dry_run": dryRun,
		"results": results,
		"summary": snapshot.PromoteSummary(results),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
