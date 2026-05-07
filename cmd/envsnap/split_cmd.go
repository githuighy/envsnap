package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

// runSplit splits a snapshot file into multiple output files based on
// prefix bucket rules supplied as "bucket=PREFIX1,PREFIX2" arguments.
//
// Usage:
//
//	envsnap split <input.json> --bucket db=DB_ --bucket cache=CACHE_ [--remainder other] [--format env|json]
func runSplit(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap split <snapshot> --bucket name=PREFIX [--remainder name] [--format json|env]")
	}

	inputPath := args[0]
	buckets := map[string][]string{}
	remainder := ""
	format := "json"

	for i := 1; i < len(args); i++ {
		switch {
		case args[i] == "--bucket" && i+1 < len(args):
			i++
			parts := strings.SplitN(args[i], "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid --bucket format, expected name=PREFIX")
			}
			prefixes := strings.Split(parts[1], ",")
			buckets[parts[0]] = prefixes
		case args[i] == "--remainder" && i+1 < len(args):
			i++
			remainder = args[i]
		case args[i] == "--format" && i+1 < len(args):
			i++
			format = args[i]
		}
	}

	snap, err := snapshot.Load(inputPath)
	if err != nil {
		return fmt.Errorf("split: load %s: %w", inputPath, err)
	}

	result, err := snapshot.Split(snap, snapshot.SplitOptions{
		Buckets:   buckets,
		Remainder: remainder,
	})
	if err != nil {
		return err
	}

	if format == "json" {
		return printSplitJSON(result)
	}

	for name, s := range result {
		out := name + ".env"
		if err := snapshot.Export(s, "env", out); err != nil {
			return fmt.Errorf("split: write %s: %w", out, err)
		}
		fmt.Fprintf(os.Stdout, "wrote %d keys to %s\n", len(s), out)
	}
	return nil
}

func printSplitJSON(result snapshot.SplitResult) error {
	out := make(map[string]map[string]string, len(result))
	for name, s := range result {
		out[name] = map[string]string(s)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
