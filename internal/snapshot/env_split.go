package snapshot

import (
	"fmt"
	"strings"
)

// SplitOptions controls how a snapshot is split into multiple named buckets.
type SplitOptions struct {
	// Buckets maps a bucket name to a list of key prefixes that belong to it.
	Buckets map[string][]string
	// Remainder is the name of the bucket that receives unmatched keys.
	// If empty, unmatched keys are discarded.
	Remainder string
}

// SplitResult holds the resulting snapshots keyed by bucket name.
type SplitResult map[string]Snapshot

// Split partitions a snapshot into multiple named snapshots according to
// prefix-based bucket rules. Keys that match no bucket are placed in the
// Remainder bucket (if configured).
func Split(snap Snapshot, opts SplitOptions) (SplitResult, error) {
	if len(opts.Buckets) == 0 {
		return nil, fmt.Errorf("envsnap split: at least one bucket must be defined")
	}

	result := make(SplitResult)
	for name := range opts.Buckets {
		result[name] = Snapshot{}
	}
	if opts.Remainder != "" {
		result[opts.Remainder] = Snapshot{}
	}

	for k, v := range snap {
		placed := false
		for bucket, prefixes := range opts.Buckets {
			for _, p := range prefixes {
				if strings.HasPrefix(k, p) {
					result[bucket][k] = v
					placed = true
					break
				}
			}
			if placed {
				break
			}
		}
		if !placed && opts.Remainder != "" {
			result[opts.Remainder][k] = v
		}
	}

	return result, nil
}
