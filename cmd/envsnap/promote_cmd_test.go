package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/user/envsnap/internal/snapshot"
)

func writePromoteSnap(t *testing.T, dir, name string, snap snapshot.Snapshot) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, snapshot.Save(path, snap))
	return path
}

func TestRunPromote_BasicPromotion(t *testing.T) {
	dir := t.TempDir()
	src := writePromoteSnap(t, dir, "src.json", snapshot.Snapshot{
		"APP_HOST": "prod.example.com",
		"APP_PORT": "443",
	})
	dst := writePromoteSnap(t, dir, "dst.json", snapshot.Snapshot{
		"APP_HOST": "staging.example.com",
	})

	err := runPromote([]string{src, dst, "--overwrite"})
	require.NoError(t, err)

	loaded, err := snapshot.Load(dst)
	require.NoError(t, err)
	assert.Equal(t, "prod.example.com", loaded["APP_HOST"])
	assert.Equal(t, "443", loaded["APP_PORT"])
}

func TestRunPromote_DryRun_DoesNotWrite(t *testing.T) {
	dir := t.TempDir()
	src := writePromoteSnap(t, dir, "src.json", snapshot.Snapshot{"KEY": "new"})
	dst := writePromoteSnap(t, dir, "dst.json", snapshot.Snapshot{"KEY": "old"})

	err := runPromote([]string{src, dst, "--overwrite", "--dry-run"})
	require.NoError(t, err)

	loaded, err := snapshot.Load(dst)
	require.NoError(t, err)
	assert.Equal(t, "old", loaded["KEY"], "dry run should not modify dst")
}

func TestRunPromote_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	src := writePromoteSnap(t, dir, "src.json", snapshot.Snapshot{"X": "1"})
	dst := writePromoteSnap(t, dir, "dst.json", snapshot.Snapshot{})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runPromote([]string{src, dst, "--format=json", "--dry-run"})
	w.Close()
	os.Stdout = old
	require.NoError(t, err)

	var out map[string]interface{}
	require.NoError(t, json.NewDecoder(r).Decode(&out))
	assert.Contains(t, out, "results")
	assert.Contains(t, out, "summary")
	assert.Equal(t, true, out["dry_run"])
}

func TestRunPromote_NoArgs(t *testing.T) {
	err := runPromote([]string{})
	assert.Error(t, err)
}

func TestRunPromote_MissingFile(t *testing.T) {
	dir := t.TempDir()
	dst := writePromoteSnap(t, dir, "dst.json", snapshot.Snapshot{})
	err := runPromote([]string{"/nonexistent/src.json", dst})
	assert.Error(t, err)
}
