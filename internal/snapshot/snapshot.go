package snapshot

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

// Snapshot represents a captured state of environment variables.
type Snapshot struct {
	Name      string            `json:"name"`
	Timestamp time.Time         `json:"timestamp"`
	Vars      map[string]string `json:"vars"`
}

// Capture reads the current process environment and returns a Snapshot.
func Capture(name string) *Snapshot {
	vars := make(map[string]string)
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		} else {
			vars[parts[0]] = ""
		}
	}
	return &Snapshot{
		Name:      name,
		Timestamp: time.Now().UTC(),
		Vars:      vars,
	}
}

// Save serialises the snapshot to a JSON file at the given path.
func Save(s *Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load deserialises a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}
