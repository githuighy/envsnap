// envsnap is a CLI tool for snapshotting and diffing environment variables
// across deployments. It supports capturing the current environment, saving
// snapshots to disk, and comparing two snapshots to detect changes.
package main

import (
	"fmt"
	"os"

	"github.com/user/envsnap/internal/snapshot"
)

const usage = `envsnap - snapshot and diff environment variables

Usage:
  envsnap capture <output-file>   Capture current env vars and save to file
  envsnap diff <file-a> <file-b>  Diff two snapshot files
  envsnap help                    Show this help message

Examples:
  envsnap capture ./snapshots/prod.json
  envsnap diff ./snapshots/staging.json ./snapshots/prod.json
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "capture":
		runCapture(os.Args[2:])
	case "diff":
		runDiff(os.Args[2:])
	case "help", "--help", "-h":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %q\n\n", os.Args[1])
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
}

// runCapture handles the "capture" subcommand. It takes a single argument
// specifying the output file path where the snapshot will be saved as JSON.
func runCapture(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: envsnap capture <output-file>")
		os.Exit(1)
	}

	outputPath := args[0]

	snap := snapshot.Capture()
	if err := snapshot.Save(snap, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "error saving snapshot: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("snapshot saved to %s (%d variables)\n", outputPath, len(snap))
}

// runDiff handles the "diff" subcommand. It loads two snapshot files and
// prints a human-readable summary of added, removed, and changed variables.
func runDiff(args []string) {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: envsnap diff <file-a> <file-b>")
		os.Exit(1)
	}

	fileA, fileB := args[0], args[1]

	snapA, err := snapshot.Load(fileA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %s: %v\n", fileA, err)
		os.Exit(1)
	}

	snapB, err := snapshot.Load(fileB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %s: %v\n", fileB, err)
		os.Exit(1)
	}

	result := snapshot.Diff(snapA, snapB)

	if len(result.Added) == 0 && len(result.Removed) == 0 && len(result.Changed) == 0 {
		fmt.Println("no differences found")
		return
	}

	for _, key := range result.Added {
		fmt.Printf("+ %s=%s\n", key, snapB[key])
	}
	for _, key := range result.Removed {
		fmt.Printf("- %s=%s\n", key, snapA[key])
	}
	for _, key := range result.Changed {
		fmt.Printf("~ %s: %q -> %q\n", key, snapA[key], snapB[key])
	}

	fmt.Printf("\nsummary: %d added, %d removed, %d changed\n",
		len(result.Added), len(result.Removed), len(result.Changed))
}
