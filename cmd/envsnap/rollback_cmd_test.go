package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"envsnap/internal/snapshot"
)

func writeRollbackSnap(t *testing.T, dir, name string, vars map[string]string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	require.NoError(t, snapshot.Save(snapshot.Snapshot{Vars: vars}, p))
	return p
}

func TestRunRollback_Basic(t *testing.T) {
	dir := t.TempDir()
	cur := writeRollbackSnap(t, dir, "cur.json", map[string]string{"HOST": "prod", "PORT": "443"})
	base := writeRollbackSnap(t, dir, "base.json", map[string]string{"HOST": "staging", "PORT": "8443"})
	out := filepath.Join(dir, "out.json")

	err := runRollback([]string{cur, base}, map[string]string{"output": out})
	require.NoError(t, err)

	loaded, err := snapshot.Load(out)
	require.NoError(t, err)
	assert.Equal(t, "staging", loaded.Vars["HOST"])
	assert.Equal(t, "8443", loaded.Vars["PORT"])
}

func TestRunRollback_DryRun_DoesNotWrite(t *testing.T) {
	dir := t.TempDir()
	cur := writeRollbackSnap(t, dir, "cur.json", map[string]string{"HOST": "prod"})
	base := writeRollbackSnap(t, dir, "base.json", map[string]string{"HOST": "staging"})
	out := filepath.Join(dir, "out.json")

	err := runRollback([]string{cur, base}, map[string]string{"output": out, "dry-run": ""})
	require.NoError(t, err)
	_, statErr := os.Stat(out)
	assert.True(t, os.IsNotExist(statErr), "dry-run should not create output file")
}

func TestRunRollback_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	cur := writeRollbackSnap(t, dir, "cur.json", map[string]string{"A": "new", "B": "extra"})
	base := writeRollbackSnap(t, dir, "base.json", map[string]string{"A": "old"})
	out := filepath.Join(dir, "out.json")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runRollback([]string{cur, base}, map[string]string{"output": out, "format": "json", "dry-run": ""})
	require.NoError(t, err)

	w.Close()
	os.Stdout = old

	var result map[string]interface{}
	require.NoError(t, json.NewDecoder(r).Decode(&result))
	restored, ok := result["restored"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "old", restored["A"])
}

func TestRunRollback_NoArgs(t *testing.T) {
	err := runRollback([]string{}, map[string]string{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "usage")
}

func TestRunRollback_MissingFile(t *testing.T) {
	dir := t.TempDir()
	base := writeRollbackSnap(t, dir, "base.json", map[string]string{"A": "v"})
	err := runRollback([]string{"/nonexistent/cur.json", base}, map[string]string{})
	require.Error(t, err)
}
