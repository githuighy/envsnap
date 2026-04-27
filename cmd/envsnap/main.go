package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/user/envsnap/internal/snapshot"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envsnap <capture|diff|merge> [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "capture":
		runCapture(os.Args[2:])
	case "diff":
		runDiff(os.Args[2:])
	case "merge":
		runMerge(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runCapture(args []string) {
	fs := flag.NewFlagSet("capture", flag.ExitOnError)
	out := fs.String("out", "snapshot.json", "output file")
	_ = fs.Parse(args)

	snap := snapshot.Capture()
	if err := snapshot.Save(snap, *out); err != nil {
		fmt.Fprintf(os.Stderr, "error saving snapshot: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("snapshot saved to %s (%d vars)\n", *out, len(snap))
}

func runDiff(args []string) {
	fs := flag.NewFlagSet("diff", flag.ExitOnError)
	format := fs.String("format", "text", "output format: text|json|env")
	_ = fs.Parse(args)

	if fs.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "usage: envsnap diff <base> <override>")
		os.Exit(1)
	}

	base, err := snapshot.Load(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading base: %v\n", err)
		os.Exit(1)
	}
	over, err := snapshot.Load(fs.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading override: %v\n", err)
		os.Exit(1)
	}

	diffs := snapshot.Diff(base, over)
	out, err := snapshot.RenderDiff(diffs, *format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(out)
}

func runMerge(args []string) {
	fs := flag.NewFlagSet("merge", flag.ExitOnError)
	prefer := fs.String("prefer", "override", "conflict resolution: base|override")
	out := fs.String("out", "", "write merged snapshot to file (prints JSON if omitted)")
	_ = fs.Parse(args)

	if fs.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "usage: envsnap merge <base> <override> [options]")
		os.Exit(1)
	}

	base, err := snapshot.Load(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading base: %v\n", err)
		os.Exit(1)
	}
	over, err := snapshot.Load(fs.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading override: %v\n", err)
		os.Exit(1)
	}

	merged := snapshot.Merge(base, over, snapshot.MergeOptions{Prefer: *prefer})

	if *out != "" {
		if err := snapshot.Save(merged, *out); err != nil {
			fmt.Fprintf(os.Stderr, "error saving merged snapshot: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("merged snapshot saved to %s (%d vars)\n", *out, len(merged))
		return
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(merged); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding output: %v\n", err)
		os.Exit(1)
	}
}
