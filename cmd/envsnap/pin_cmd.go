package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

// runPin captures the current values of specified keys from a snapshot file
// and writes a pinned snapshot to an output file for later drift detection.
//
// Usage: envsnap pin <input> <output> [--keys KEY1,KEY2] [--allow-missing]
func runPin(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envsnap pin <input> <output> [--keys k1,k2] [--allow-missing]")
	}

	inputPath := args[0]
	outputPath := args[1]

	var explicitKeys []string
	allowMissing := false

	for i := 2; i < len(args); i++ {
		switch {
		case strings.HasPrefix(args[i], "--keys="):
			val := strings.TrimPrefix(args[i], "--keys=")
			for _, k := range strings.Split(val, ",") {
				if k = strings.TrimSpace(k); k != "" {
					explicitKeys = append(explicitKeys, k)
				}
			}
		case args[i] == "--allow-missing":
			allowMissing = true
		}
	}

	snap, err := snapshot.Load(inputPath)
	if err != nil {
		return fmt.Errorf("pin: load %q: %w", inputPath, err)
	}

	result, err := snapshot.Pin(snap, snapshot.PinOptions{
		Keys:         explicitKeys,
		AllowMissing: allowMissing,
	})
	if err != nil {
		return fmt.Errorf("pin: %w", err)
	}

	if len(result.Skipped) > 0 {
		fmt.Fprintf(os.Stderr, "pin: skipped missing keys: %s\n", strings.Join(result.Skipped, ", "))
	}

	pinSnap := snapshot.PinToSnapshot(result)

	if err := snapshot.Save(pinSnap, outputPath); err != nil {
		return fmt.Errorf("pin: save %q: %w", outputPath, err)
	}

	summary := map[string]interface{}{
		"pinned":  len(result.Pinned),
		"skipped": result.Skipped,
		"output":  outputPath,
	}
	out, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Println(string(out))
	return nil
}
