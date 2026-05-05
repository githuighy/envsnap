package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExportFormat defines supported export formats.
type ExportFormat string

const (
	ExportFormatEnv   ExportFormat = "env"
	ExportFormatJSON  ExportFormat = "json"
	ExportFormatShell ExportFormat = "shell"
)

// ExportOptions controls how a snapshot is exported.
type ExportOptions struct {
	Format  ExportFormat
	OutFile string // empty means stdout
}

// Export writes a snapshot to a file or stdout in the requested format.
func Export(snap Snapshot, opts ExportOptions) error {
	var content string
	var err error

	switch opts.Format {
	case ExportFormatJSON:
		content, err = exportJSON(snap)
	case ExportFormatShell:
		content, err = exportShell(snap), nil
	case ExportFormatEnv, "":
		content, err = exportEnv(snap), nil
	default:
		return fmt.Errorf("unsupported export format: %q", opts.Format)
	}

	if err != nil {
		return fmt.Errorf("export: %w", err)
	}

	if opts.OutFile == "" {
		_, err = fmt.Fprint(os.Stdout, content)
		return err
	}

	if err := os.MkdirAll(filepath.Dir(opts.OutFile), 0o755); err != nil {
		return fmt.Errorf("export: mkdir: %w", err)
	}
	if err := os.WriteFile(opts.OutFile, []byte(content), 0o644); err != nil {
		return fmt.Errorf("export: write file %q: %w", opts.OutFile, err)
	}
	return nil
}

func exportEnv(snap Snapshot) string {
	var sb strings.Builder
	for k, v := range snap.Vars {
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}
	return sb.String()
}

func exportShell(snap Snapshot) string {
	var sb strings.Builder
	for k, v := range snap.Vars {
		fmt.Fprintf(&sb, "export %s=%q\n", k, v)
	}
	return sb.String()
}

func exportJSON(snap Snapshot) (string, error) {
	b, err := json.MarshalIndent(snap.Vars, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}
