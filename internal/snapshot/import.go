package snapshot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// ImportOptions controls how a snapshot is imported.
type ImportOptions struct {
	Format ImportFormat
}

// ImportFormat defines supported import formats.
type ImportFormat string

const (
	ImportFormatEnv  ImportFormat = "env"
	ImportFormatJSON ImportFormat = "json"
)

// Import reads a file and returns a Snapshot.
func Import(path string, opts ImportOptions) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("import: read file: %w", err)
	}

	fmt := opts.Format
	if fmt == "" {
		fmt = detectFormat(path)
	}

	var vars map[string]string
	switch fmt {
	case ImportFormatJSON:
		if err := json.Unmarshal(data, &vars); err != nil {
			return Snapshot{}, fmt.Errorf("import: parse JSON: %w", err)
		}
	case ImportFormatEnv, "":
		vars, err = parseEnvFile(string(data))
		if err != nil {
			return Snapshot{}, fmt.Errorf("import: parse env: %w", err)
		}
	default:
		return Snapshot{}, fmt.Errorf("import: unsupported format: %q", fmt)
	}

	return Snapshot{
		Vars:      vars,
		CapturedAt: time.Now(),
	}, nil
}

func parseEnvFile(content string) (map[string]string, error) {
	vars := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		vars[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return vars, scanner.Err()
}

func detectFormat(path string) ImportFormat {
	if strings.HasSuffix(path, ".json") {
		return ImportFormatJSON
	}
	return ImportFormatEnv
}
