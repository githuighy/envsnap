package snapshot

import (
	"fmt"
	"strings"
)

// PatchOp represents a single patch operation.
type PatchOp struct {
	Op    string // "set", "delete", "rename"
	Key   string
	Value string // used by "set"
	To    string // used by "rename"
}

// PatchResult holds the outcome of applying a patch.
type PatchResult struct {
	Applied  []string
	Skipped  []string
	Snapshot Snapshot
}

// Patch applies a sequence of patch operations to a snapshot.
// Supported ops: "set", "delete", "rename".
// Unknown ops or invalid arguments return an error.
func Patch(snap Snapshot, ops []PatchOp) (PatchResult, error) {
	out := make(Snapshot, len(snap))
	for k, v := range snap {
		out[k] = v
	}

	var applied, skipped []string

	for _, op := range ops {
		if strings.TrimSpace(op.Key) == "" {
			return PatchResult{}, fmt.Errorf("patch op %q has empty key", op.Op)
		}

		switch op.Op {
		case "set":
			out[op.Key] = op.Value
			applied = append(applied, op.Key)

		case "delete":
			if _, exists := out[op.Key]; exists {
				delete(out, op.Key)
				applied = append(applied, op.Key)
			} else {
				skipped = append(skipped, op.Key)
			}

		case "rename":
			if strings.TrimSpace(op.To) == "" {
				return PatchResult{}, fmt.Errorf("rename op for key %q has empty 'to' field", op.Key)
			}
			if val, exists := out[op.Key]; exists {
				out[op.To] = val
				delete(out, op.Key)
				applied = append(applied, op.Key+"→"+op.To)
			} else {
				skipped = append(skipped, op.Key)
			}

		default:
			return PatchResult{}, fmt.Errorf("unknown patch op: %q", op.Op)
		}
	}

	return PatchResult{
		Applied:  applied,
		Skipped:  skipped,
		Snapshot: out,
	}, nil
}
