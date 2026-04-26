package snapshot

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

// OutputFormat defines the format for rendering diffs and snapshots.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
	FormatEnv  OutputFormat = "env"
)

// RenderDiff writes a human-readable or structured diff to the given writer.
func RenderDiff(w io.Writer, diffs []DiffEntry, format OutputFormat) error {
	switch format {
	case FormatJSON:
		return renderDiffJSON(w, diffs)
	case FormatEnv:
		return renderDiffEnv(w, diffs)
	default:
		return renderDiffText(w, diffs)
	}
}

func renderDiffText(w io.Writer, diffs []DiffEntry) error {
	if len(diffs) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "STATUS\tKEY\tOLD VALUE\tNEW VALUE")
	fmt.Fprintln(tw, "------\t---\t---------\t---------")
	sort.Slice(diffs, func(i, j int) bool { return diffs[i].Key < diffs[j].Key })
	for _, d := range diffs {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", d.Status, d.Key, d.OldValue, d.NewValue)
	}
	return tw.Flush()
}

func renderDiffJSON(w io.Writer, diffs []DiffEntry) error {
	sort.Slice(diffs, func(i, j int) bool { return diffs[i].Key < diffs[j].Key })
	fmt.Fprintln(w, "[")
	for i, d := range diffs {
		comma := ","
		if i == len(diffs)-1 {
			comma = ""
		}
		fmt.Fprintf(w, "  {\"status\":%q,\"key\":%q,\"old\":%q,\"new\":%q}%s\n",
			d.Status, d.Key, d.OldValue, d.NewValue, comma)
	}
	fmt.Fprintln(w, "]")
	return nil
}

func renderDiffEnv(w io.Writer, diffs []DiffEntry) error {
	sort.Slice(diffs, func(i, j int) bool { return diffs[i].Key < diffs[j].Key })
	for _, d := range diffs {
		prefix := "#"
		switch d.Status {
		case StatusAdded:
			prefix = "+"
		case StatusRemoved:
			prefix = "-"
		case StatusChanged:
			prefix = "~"
		}
		val := strings.ReplaceAll(d.NewValue, "\n", "\\n")
		fmt.Fprintf(w, "%s %s=%s\n", prefix, d.Key, val)
	}
	return nil
}
