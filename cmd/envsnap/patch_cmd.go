package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

// runPatch applies a series of patch operations to a snapshot file.
//
// Usage:
//
//	envsnap patch <snapshot> --op set:KEY=VALUE --op delete:KEY --op rename:OLD=NEW [--format text|json] [--out file]
func runPatch(args []string, flags map[string]string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap patch <snapshot> [--op op:args...]")
	}

	snap, err := snapshot.Load(args[0])
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	rawOps, _ := flags["op"]
	var ops []snapshot.PatchOp
	for _, raw := range strings.Split(rawOps, ",") {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		op, err := parsePatchOp(raw)
		if err != nil {
			return fmt.Errorf("parse op %q: %w", raw, err)
		}
		ops = append(ops, op)
	}

	result, err := snapshot.Patch(snap, ops)
	if err != nil {
		return fmt.Errorf("patch: %w", err)
	}

	outFile := flags["out"]
	if outFile != "" {
		if err := snapshot.Save(result.Snapshot, outFile); err != nil {
			return fmt.Errorf("save: %w", err)
		}
	}

	format := flags["format"]
	if format == "json" {
		return printPatchJSON(result)
	}

	fmt.Printf("Applied: %d  Skipped: %d\n", len(result.Applied), len(result.Skipped))
	for _, k := range result.Applied {
		fmt.Printf("  + %s\n", k)
	}
	for _, k := range result.Skipped {
		fmt.Printf("  ~ %s (skipped)\n", k)
	}
	return nil
}

func parsePatchOp(raw string) (snapshot.PatchOp, error) {
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return snapshot.PatchOp{}, fmt.Errorf("expected format op:args, got %q", raw)
	}
	kind, payload := parts[0], parts[1]
	switch kind {
	case "set":
		kv := strings.SplitN(payload, "=", 2)
		if len(kv) != 2 {
			return snapshot.PatchOp{}, fmt.Errorf("set requires KEY=VALUE, got %q", payload)
		}
		return snapshot.PatchOp{Op: "set", Key: kv[0], Value: kv[1]}, nil
	case "delete":
		return snapshot.PatchOp{Op: "delete", Key: payload}, nil
	case "rename":
		kv := strings.SplitN(payload, "=", 2)
		if len(kv) != 2 {
			return snapshot.PatchOp{}, fmt.Errorf("rename requires OLD=NEW, got %q", payload)
		}
		return snapshot.PatchOp{Op: "rename", Key: kv[0], To: kv[1]}, nil
	default:
		return snapshot.PatchOp{}, fmt.Errorf("unknown op kind %q", kind)
	}
}

func printPatchJSON(result snapshot.PatchResult) error {
	out := map[string]interface{}{
		"applied": result.Applied,
		"skipped": result.Skipped,
		"snapshot": result.Snapshot,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
