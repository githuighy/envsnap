package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

// runScore loads a snapshot file and prints its health score.
// Usage: envsnap score <snapshot-file> [--required KEY,...] [--forbidden KEY,...] [--format text|json]
func runScore(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap score <snapshot-file> [--required KEY,...] [--forbidden KEY,...] [--format text|json]")
	}

	file := args[0]
	var required, forbidden []string
	format := "text"

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--required":
			if i+1 < len(args) {
				i++
				required = strings.Split(args[i], ",")
			}
		case "--forbidden":
			if i+1 < len(args) {
				i++
				forbidden = strings.Split(args[i], ",")
			}
		case "--format":
			if i+1 < len(args) {
				i++
				format = args[i]
			}
		}
	}

	snap, err := snapshot.Load(file)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	lo := snapshot.DefaultLintOptions()
	res := snapshot.Score(snap, snapshot.ScoreOptions{
		RequiredKeys:  required,
		ForbiddenKeys: forbidden,
		LintOptions:   &lo,
	})

	switch format {
	case "json":
		return printScoreJSON(res)
	default:
		return printScoreText(res)
	}
}

func printScoreText(res snapshot.ScoreResult) error {
	fmt.Fprintf(os.Stdout, "Score: %d/100\n", res.Score)
	if len(res.Deductions) == 0 {
		fmt.Fprintln(os.Stdout, "No issues found.")
		return nil
	}
	fmt.Fprintln(os.Stdout, "Deductions:")
	for _, d := range res.Deductions {
		fmt.Fprintf(os.Stdout, "  - %s\n", d)
	}
	return nil
}

func printScoreJSON(res snapshot.ScoreResult) error {
	out := struct {
		Score      int      `json:"score"`
		Deductions []string `json:"deductions"`
	}{
		Score:      res.Score,
		Deductions: res.Deductions,
	}
	if out.Deductions == nil {
		out.Deductions = []string{}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
