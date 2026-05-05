package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

// runTransform applies a built-in transform operation to a snapshot file.
// Usage: envsnap transform <file> <op> [--prefix P] [--keys k1,k2] [--skip-errors] [--format text|json]
func runTransform(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envsnap transform <file> <op> [--prefix P] [--keys k1,k2] [--skip-errors] [--format text|json]")
	}

	filePath := args[0]
	op := args[1]

	var prefix string
	var keys []string
	skipErrors := false
	format := "text"

	for i := 2; i < len(args); i++ {
		switch {
		case args[i] == "--skip-errors":
			skipErrors = true
		case strings.HasPrefix(args[i], "--prefix="):
			prefix = strings.TrimPrefix(args[i], "--prefix=")
		case strings.HasPrefix(args[i], "--keys="):
			raw := strings.TrimPrefix(args[i], "--keys=")
			keys = strings.Split(raw, ",")
		case strings.HasPrefix(args[i], "--format="):
			format = strings.TrimPrefix(args[i], "--format=")
		}
	}

	snap, err := snapshot.Load(filePath)
	if err != nil {
		return fmt.Errorf("transform: load %q: %w", filePath, err)
	}

	fn, err := resolveTransformOp(op)
	if err != nil {
		return err
	}

	opts := snapshot.TransformOptions{
		Keys:       keys,
		Prefix:     prefix,
		SkipErrors: skipErrors,
	}

	result, err := snapshot.Transform(snap, fn, opts)
	if err != nil {
		return fmt.Errorf("transform: %w", err)
	}

	switch format {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(result)
	default:
		for k, v := range result {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
	}
	return nil
}

func resolveTransformOp(op string) (snapshot.TransformFunc, error) {
	switch op {
	case "uppercase":
		return func(k, v string) (string, error) { return strings.ToUpper(v), nil }, nil
	case "lowercase":
		return func(k, v string) (string, error) { return strings.ToLower(v), nil }, nil
	case "trim":
		return func(k, v string) (string, error) { return strings.TrimSpace(v), nil }, nil
	default:
		return nil, fmt.Errorf("transform: unknown op %q (supported: uppercase, lowercase, trim)", op)
	}
}
