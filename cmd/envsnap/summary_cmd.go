package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/user/envsnap/internal/snapshot"
)

// runSummary loads two snapshot files, diffs them, and prints a summary.
func runSummary(args []string, format string, byPrefix bool) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envsnap summary <before> <after> [--format text|json] [--by-prefix]")
	}

	before, err := snapshot.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading before snapshot: %w", err)
	}
	after, err := snapshot.Load(args[1])
	if err != nil {
		return fmt.Errorf("loading after snapshot: %w", err)
	}

	diffs := snapshot.Diff(before, after)

	if byPrefix {
		return printSummaryByPrefix(diffs, format)
	}

	summary := snapshot.SummariseDiff(diffs)

	if format == "json" {
		return printSummaryJSON(summary)
	}
	return printSummaryText(summary)
}

func printSummaryText(s snapshot.DiffSummary) error {
	if !s.HasChanges() {
		fmt.Fprintln(os.Stdout, "No changes detected.")
		return nil
	}
	fmt.Fprintln(os.Stdout, s.String())
	return nil
}

func printSummaryJSON(s snapshot.DiffSummary) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

func printSummaryByPrefix(diffs []snapshot.DiffEntry, format string) error {
	byPrefix := snapshot.DiffSummaryByPrefix(diffs)

	keys := make([]string, 0, len(byPrefix))
	for k := range byPrefix {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(byPrefix)
	}

	for _, prefix := range keys {
		s := byPrefix[prefix]
		fmt.Fprintf(os.Stdout, "[%s] %s\n", prefix, s.String())
	}
	return nil
}
