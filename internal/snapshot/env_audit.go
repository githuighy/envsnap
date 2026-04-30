package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// AuditEventType describes the kind of change recorded in an audit entry.
type AuditEventType string

const (
	AuditAdded   AuditEventType = "added"
	AuditRemoved AuditEventType = "removed"
	AuditChanged AuditEventType = "changed"
)

// AuditEntry represents a single recorded change to an environment variable.
type AuditEntry struct {
	Timestamp time.Time      `json:"timestamp"`
	Key       string         `json:"key"`
	Event     AuditEventType `json:"event"`
	OldValue  string         `json:"old_value,omitempty"`
	NewValue  string         `json:"new_value,omitempty"`
	Source    string         `json:"source,omitempty"`
}

// AuditLog holds a sequence of audit entries.
type AuditLog struct {
	Entries []AuditEntry `json:"entries"`
}

// Audit compares two snapshots and returns an AuditLog describing all changes.
// source is an optional label (e.g. deployment name) attached to each entry.
func Audit(before, after Snapshot, source string) AuditLog {
	log := AuditLog{}
	now := time.Now().UTC()

	keys := unionAuditKeys(before, after)
	sort.Strings(keys)

	for _, k := range keys {
		oldVal, inBefore := before[k]
		newVal, inAfter := after[k]

		switch {
		case inBefore && !inAfter:
			log.Entries = append(log.Entries, AuditEntry{
				Timestamp: now,
				Key:       k,
				Event:     AuditRemoved,
				OldValue:  oldVal,
				Source:    source,
			})
		case !inBefore && inAfter:
			log.Entries = append(log.Entries, AuditEntry{
				Timestamp: now,
				Key:       k,
				Event:     AuditAdded,
				NewValue:  newVal,
				Source:    source,
			})
		case oldVal != newVal:
			log.Entries = append(log.Entries, AuditEntry{
				Timestamp: now,
				Key:       k,
				Event:     AuditChanged,
				OldValue:  oldVal,
				NewValue:  newVal,
				Source:    source,
			})
		}
	}

	return log
}

// Summary returns a human-readable summary line for the audit log.
func (a AuditLog) Summary() string {
	var added, removed, changed int
	for _, e := range a.Entries {
		switch e.Event {
		case AuditAdded:
			added++
		case AuditRemoved:
			removed++
		case AuditChanged:
			changed++
		}
	}
	return fmt.Sprintf("audit: %d added, %d removed, %d changed", added, removed, changed)
}

func unionAuditKeys(a, b Snapshot) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
