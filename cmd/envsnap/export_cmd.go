package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envsnap/internal/snapshot"
)

func runExport(args []string) error {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	format := fs.String("format", "env", "output format: env, json, shell")
	outFile := fs.String("out", "", "output file path (default: stdout)")
	inFile := fs.String("in", "", "snapshot file to export (required)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *inFile == "" {
		return fmt.Errorf("export: --in is required")
	}

	snap, err := snapshot.Load(*inFile)
	if err != nil {
		return fmt.Errorf("export: load snapshot: %w", err)
	}

	opts := snapshot.ExportOptions{
		Format:  snapshot.ExportFormat(*format),
		OutFile: *outFile,
	}

	if err := snapshot.Export(snap, opts); err != nil {
		return fmt.Errorf("export: %w", err)
	}
	return nil
}

func runImport(args []string) error {
	fs := flag.NewFlagSet("import", flag.ContinueOnError)
	format := fs.String("format", "", "input format: env, json (auto-detected if empty)")
	inFile := fs.String("in", "", "file to import (required)")
	outFile := fs.String("out", "", "snapshot output path (required)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *inFile == "" {
		return fmt.Errorf("import: --in is required")
	}
	if *outFile == "" {
		return fmt.Errorf("import: --out is required")
	}

	opts := snapshot.ImportOptions{
		Format: snapshot.ImportFormat(*format),
	}

	snap, err := snapshot.Import(*inFile, opts)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}

	if err := snapshot.Save(snap, *outFile); err != nil {
		return fmt.Errorf("import: save snapshot: %w", err)
	}

	fmt.Fprintf(os.Stderr, "imported %d variables into %s\n", len(snap.Vars), *outFile)
	return nil
}
