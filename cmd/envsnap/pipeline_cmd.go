package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

// runPipeline loads a snapshot file and applies a sequence of built-in
// pipeline steps specified as comma-separated names via --steps.
//
// Supported step names:
//   filter:<PREFIX>   – keep only keys with the given prefix
//   redact            – redact sensitive keys with default settings
//   lint              – abort pipeline if lint issues are found
func runPipeline(args []string, steps string, format string, outputPath string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap pipeline <snapshot-file> --steps <step1,step2,...>")
	}

	snap, err := snapshot.Load(args[0])
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	p := snapshot.NewPipeline()
	for _, raw := range strings.Split(steps, ",") {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		switch {
		case strings.HasPrefix(raw, "filter:"):
			prefix := strings.TrimPrefix(raw, "filter:")
			p.AddStep(raw, func(s snapshot.Snapshot) (snapshot.Snapshot, error) {
				return snapshot.Filter(s, snapshot.FilterOptions{Prefixes: []string{prefix}}), nil
			})
		case raw == "redact":
			p.AddStep("redact", func(s snapshot.Snapshot) (snapshot.Snapshot, error) {
				return snapshot.Redact(s, snapshot.RedactOptions{}), nil
			})
		case raw == "lint":
			p.AddStep("lint", func(s snapshot.Snapshot) (snapshot.Snapshot, error) {
				issues := snapshot.Lint(s, snapshot.DefaultLintOptions())
				if len(issues) > 0 {
					return nil, fmt.Errorf("lint failed: %d issue(s) found", len(issues))
				}
				return s, nil
			})
		default:
			return fmt.Errorf("unknown pipeline step: %q", raw)
		}
	}

	out, err := p.Final(snap)
	if err != nil {
		return err
	}

	if format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	if outputPath != "" {
		return snapshot.Save(out, outputPath)
	}

	for k, v := range out {
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}
