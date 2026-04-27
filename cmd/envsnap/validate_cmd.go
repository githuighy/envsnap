package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

// runValidate loads a snapshot file and validates it against the provided
// required/forbidden key lists and optional max key length.
func runValidate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap validate <snapshot> [--require K1,K2] [--forbid K3,K4] [--max-key-len N] [--json]")
	}

	snapshotFile := args[0]
	var required, forbidden []string
	maxKeyLen := 0
	jsonOutput := false

	for i := 1; i < len(args); i++ {
		switch {
		case args[i] == "--require" && i+1 < len(args):
			i++
			required = strings.Split(args[i], ",")
		case args[i] == "--forbid" && i+1 < len(args):
			i++
			forbidden = strings.Split(args[i], ",")
		case args[i] == "--max-key-len" && i+1 < len(args):
			i++
			fmt.Sscanf(args[i], "%d", &maxKeyLen)
		case args[i] == "--json":
			jsonOutput = true
		}
	}

	snap, err := snapshot.Load(snapshotFile)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	validateErr := snapshot.Validate(snap, snapshot.ValidateOptions{
		RequiredKeys:  required,
		ForbiddenKeys: forbidden,
		MaxKeyLength:  maxKeyLen,
	})

	if jsonOutput {
		return printValidateJSON(validateErr)
	}

	if validateErr == nil {
		fmt.Println("✓ Snapshot is valid")
		return nil
	}
	ve := validateErr.(*snapshot.ValidationError)
	fmt.Fprintln(os.Stderr, "✗ Validation failed:")
	for _, issue := range ve.Issues {
		fmt.Fprintf(os.Stderr, "  - %s\n", issue)
	}
	return validateErr
}

func printValidateJSON(validateErr error) error {
	type result struct {
		Valid  bool     `json:"valid"`
		Issues []string `json:"issues,omitempty"`
	}
	r := result{Valid: validateErr == nil}
	if validateErr != nil {
		r.Issues = validateErr.(*snapshot.ValidationError).Issues
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}
	if !r.Valid {
		return validateErr
	}
	return nil
}
