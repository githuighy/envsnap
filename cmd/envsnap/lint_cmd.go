package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envsnap/internal/snapshot"
)

// runLint loads a snapshot file and runs lint checks against it.
// Exits with code 2 if any lint issues are found.
func runLint(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: envsnap lint <snapshot-file> [--format text|json] [--no-lowercase] [--no-empty] [--no-whitespace]")
		os.Exit(1)
	}

	path := args[0]
	format := "text"
	opts := snapshot.DefaultLintOptions()

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--format":
			if i+1 < len(args) {
				i++
				format = args[i]
			} else {
				fmt.Fprintln(os.Stderr, "error: --format requires an argument (text or json)")
				os.Exit(1)
			}
		case "--no-lowercase":
			opts.WarnLowercase = false
		case "--no-empty":
			opts.WarnEmpty = false
		case "--no-whitespace":
			opts.WarnWhitespace = false
		}
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading snapshot: %v\n", err)
		os.Exit(1)
	}

	issues := snapshot.Lint(snap, opts)

	switch format {
	case "json":
		printLintJSON(issues)
	case "text":
		printLintText(issues)
	default:
		fmt.Fprintf(os.Stderr, "error: unknown format %q, expected text or json\n", format)
		os.Exit(1)
	}

	if len(issues) > 0 {
		os.Exit(2)
	}
}

func printLintText(issues []snapshot.LintIssue) {
	if len(issues) == 0 {
		fmt.Println("No lint issues found.")
		return
	}
	fmt.Printf("%d issue(s) found:\n", len(issues))
	for _, issue := range issues {
		fmt.Printf("  [%s] %s\n", issue.Rule, issue.Message)
	}
}

func printLintJSON(issues []snapshot.LintIssue) {
	if issues == nil {
		issues = []snapshot.LintIssue{}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(issues); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}
