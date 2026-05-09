package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsnap/internal/snapshot"
)

func runReport(args []string) error {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	format := fs.String("format", "text", "output format: text|json")
	byPrefix := fs.Bool("by-prefix", false, "group key counts by prefix")

	if err := fs.Parse(args); err != nil {
		return err
	}

	positional := fs.Args()
	if len(positional) < 2 {
		return fmt.Errorf("usage: envsnap report <snapshot-a> <snapshot-b> [flags]")
	}

	snapA, err := snapshot.Load(positional[0])
	if err != nil {
		return fmt.Errorf("loading snapshot A: %w", err)
	}

	snapB, err := snapshot.Load(positional[1])
	if err != nil {
		return fmt.Errorf("loading snapshot B: %w", err)
	}

	diffs := snapshot.Diff(*snapA, *snapB)

	// Convert snapshot.Diff results to DiffEntry slice expected by report.
	entries := make([]snapshot.DiffEntry, 0, len(diffs))
	for _, d := range diffs {
		entries = append(entries, snapshot.DiffEntry{
			Key:      d.Key,
			Status:   d.Status,
			OldValue: d.OldValue,
			NewValue: d.NewValue,
		})
	}

	report := snapshot.BuildCompareReport(entries, *byPrefix)

	out, err := snapshot.RenderCompareReport(report, *format)
	if err != nil {
		return fmt.Errorf("rendering report: %w", err)
	}

	fmt.Fprint(os.Stdout, out)
	return nil
}
