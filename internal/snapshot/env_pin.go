package snapshot

import (
	"fmt"
	"sort"
)

// PinOptions controls how pinning behaves.
type PinOptions struct {
	// Keys is the explicit list of keys to pin. If empty, all keys are pinned.
	Keys []string
	// AllowMissing skips keys not present in the snapshot instead of returning an error.
	AllowMissing bool
}

// PinResult holds the pinned key-value pairs and any keys that were skipped.
type PinResult struct {
	// Pinned maps each key to its locked value.
	Pinned map[string]string
	// Skipped contains keys requested but absent from the snapshot.
	Skipped []string
}

// Pin records the current values of the specified keys (or all keys) from snap
// so they can later be asserted against via Lock. It returns a PinResult
// containing the pinned values and any skipped keys.
func Pin(snap Snapshot, opts PinOptions) (PinResult, error) {
	result := PinResult{
		Pinned: make(map[string]string),
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range snap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	for _, k := range keys {
		v, ok := snap[k]
		if !ok {
			if opts.AllowMissing {
				result.Skipped = append(result.Skipped, k)
				continue
			}
			return PinResult{}, fmt.Errorf("pin: key %q not found in snapshot", k)
		}
		result.Pinned[k] = v
	}

	return result, nil
}

// PinToSnapshot converts a PinResult into a Snapshot so it can be saved to
// disk and later used with Lock for drift detection.
func PinToSnapshot(result PinResult) Snapshot {
	snap := make(Snapshot, len(result.Pinned))
	for k, v := range result.Pinned {
		snap[k] = v
	}
	return snap
}
