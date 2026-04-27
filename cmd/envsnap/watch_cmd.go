package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/user/envsnap/internal/snapshot"
)

func runWatch(args []string) {
	fs := flag.NewFlagSet("watch", flag.ExitOnError)
	interval := fs.Duration("interval", 5*time.Second, "polling interval (e.g. 2s, 500ms)")
	prefixList := fs.String("prefix", "", "comma-separated list of key prefixes to watch")
	excludeList := fs.String("exclude", "", "comma-separated list of key prefixes to exclude")
	format := fs.String("format", "text", "output format: text|json|env")
	fs.Parse(args)

	var prefixes, excludes []string
	if *prefixList != "" {
		prefixes = strings.Split(*prefixList, ",")
	}
	if *excludeList != "" {
		excludes = strings.Split(*excludeList, ",")
	}

	fmt.Fprintf(os.Stderr, "Watching environment every %s (press Ctrl+C to stop)...\n", *interval)

	wr := snapshot.Watch(snapshot.WatchOptions{
		Interval: *interval,
		Prefixes: prefixes,
		Exclude:  excludes,
		OnChange: func(diffs []snapshot.DiffEntry) {
			timestamp := time.Now().Format(time.RFC3339)
			fmt.Fprintf(os.Stderr, "[%s] %d change(s) detected:\n", timestamp, len(diffs))
			out, err := snapshot.RenderDiff(diffs, *format)
			if err != nil {
				fmt.Fprintf(os.Stderr, "render error: %v\n", err)
				return
			}
			fmt.Println(out)
		},
		OnError: func(err error) {
			fmt.Fprintf(os.Stderr, "watch error: %v\n", err)
		},
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	wr.Stop()
	fmt.Fprintln(os.Stderr, "Watch stopped.")
}
