package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

func templateStore() (*snapshot.TemplateStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return snapshot.NewTemplateStore(filepath.Join(home, ".envsnap", "templates"))
}

func runTemplateSave(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap template save <name> [--required KEY,...] [--default KEY=VALUE,...]")
	}
	name := args[0]
	tmpl := snapshot.Template{Name: name, Defaults: map[string]string{}, Required: []string{}}

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--required":
			i++
			if i < len(args) {
				tmpl.Required = strings.Split(args[i], ",")
			}
		case "--default":
			i++
			if i < len(args) {
				for _, pair := range strings.Split(args[i], ",") {
					parts := strings.SplitN(pair, "=", 2)
					if len(parts) == 2 {
						tmpl.Defaults[parts[0]] = parts[1]
					}
				}
			}
		}
	}

	store, err := templateStore()
	if err != nil {
		return err
	}
	if err := store.Save(tmpl); err != nil {
		return err
	}
	fmt.Printf("template %q saved\n", name)
	return nil
}

func runTemplateApply(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envsnap template apply <template-name> <snapshot-file>")
	}
	store, err := templateStore()
	if err != nil {
		return err
	}
	tmpl, err := store.Load(args[0])
	if err != nil {
		return err
	}
	snap, err := snapshot.Load(args[1])
	if err != nil {
		return err
	}
	result, err := snapshot.ApplyTemplate(snap, tmpl)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}

func runTemplateList(args []string) error {
	store, err := templateStore()
	if err != nil {
		return err
	}
	names, err := store.List()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Println("no templates saved")
		return nil
	}
	for _, n := range names {
		fmt.Println(n)
	}
	return nil
}
