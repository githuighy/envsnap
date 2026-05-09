package snapshot

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// CompareReport holds a structured summary of a comparison between two snapshots.
type CompareReport struct {
	TotalKeys   int            `json:"total_keys"`
	Added       int            `json:"added"`
	Removed     int            `json:"removed"`
	Changed     int            `json:"changed"`
	Unchanged   int            `json:"unchanged"`
	Similarity  float64        `json:"similarity_pct"`
	ByPrefix    map[string]int `json:"by_prefix,omitempty"`
}

// BuildCompareReport generates a CompareReport from a slice of DiffEntry values.
func BuildCompareReport(diffs []DiffEntry, groupByPrefix bool) CompareReport {
	r := CompareReport{
		ByPrefix: make(map[string]int),
	}

	for _, d := range diffs {
		r.TotalKeys++
		switch d.Status {
		case "added":
			r.Added++
		case "removed":
			r.Removed++
		case "changed":
			r.Changed++
		case "unchanged":
			r.Unchanged++
		}
		if groupByPrefix {
			pfx := keyPrefix(d.Key)
			r.ByPrefix[pfx]++
		}
	}

	if r.TotalKeys > 0 {
		r.Similarity = float64(r.Unchanged) / float64(r.TotalKeys) * 100.0
	}
	if !groupByPrefix {
		r.ByPrefix = nil
	}
	return r
}

// RenderCompareReport formats a CompareReport as text or JSON.
func RenderCompareReport(rep CompareReport, format string) (string, error) {
	switch format {
	case "json":
		b, err := json.MarshalIndent(rep, "", "  ")
		if err != nil {
			return "", err
		}
		return string(b), nil
	default:
		var sb strings.Builder
		fmt.Fprintf(&sb, "Total keys : %d\n", rep.TotalKeys)
		fmt.Fprintf(&sb, "Added      : %d\n", rep.Added)
		fmt.Fprintf(&sb, "Removed    : %d\n", rep.Removed)
		fmt.Fprintf(&sb, "Changed    : %d\n", rep.Changed)
		fmt.Fprintf(&sb, "Unchanged  : %d\n", rep.Unchanged)
		fmt.Fprintf(&sb, "Similarity : %.1f%%\n", rep.Similarity)
		if len(rep.ByPrefix) > 0 {
			prefixes := make([]string, 0, len(rep.ByPrefix))
			for p := range rep.ByPrefix {
				prefixes = append(prefixes, p)
			}
			sort.Strings(prefixes)
			sb.WriteString("\nBy prefix:\n")
			for _, p := range prefixes {
				fmt.Fprintf(&sb, "  %-20s %d\n", p, rep.ByPrefix[p])
			}
		}
		return sb.String(), nil
	}
}
