package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"envsnap/internal/snapshot"
)

func saveRollbackSnap(t *testing.T, dir, name string, vars map[string]string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := snapshot.Save(snapshot.Snapshot{Vars: vars}, path)
	require.NoError(t, err)
	return path
}

func TestRollback_RoundTripFromDisk(t *testing.T) {
	dir := t.TempDir()
	currentPath := saveRollbackSnap(t, dir, "current.json", map[string]string{
		"HOST": "prod.host",
		"PORT": "443",
	})
	baselinePath := saveRollbackSnap(t, dir, "baseline.json", map[string]string{
		"HOST": "staging.host",
		"PORT": "8443",
	})

	current, err := snapshot.Load(currentPath)
	require.NoError(t, err)
	baseline, err := snapshot.Load(baselinePath)
	require.NoError(t, err)

	res, err := snapshot.Rollback(current, baseline, snapshot.RollbackOptions{})
	require.NoError(t, err)

	outPath := filepath.Join(dir, "rolled-back.json")
	err = snapshot.Save(res.Snapshot, outPath)
	require.NoError(t, err)

	loaded, err := snapshot.Load(outPath)
	require.NoError(t, err)
	assert.Equal(t, "staging.host", loaded.Vars["HOST"])
	assert.Equal(t, "8443", loaded.Vars["PORT"])
}

func TestRollback_SummaryAccurate(t *testing.T) {
	dir := t.TempDir()
	currentPath := saveRollbackSnap(t, dir, "current.json", map[string]string{
		"A": "new",
		"B": "same",
		"C": "extra",
	})
	baselinePath := saveRollbackSnap(t, dir, "baseline.json", map[string]string{
		"A": "old",
		"B": "same",
	})

	current, err := snapshot.Load(currentPath)
	require.NoError(t, err)
	baseline, err := snapshot.Load(baselinePath)
	require.NoError(t, err)

	res, err := snapshot.Rollback(current, baseline, snapshot.RollbackOptions{})
	require.NoError(t, err)

	assert.Equal(t, map[string]string{"A": "old"}, res.Restored)
	assert.Contains(t, res.Dropped, "C")
	assert.Contains(t, res.Unchanged, "B")

	_ = os.Remove(currentPath)
	_ = os.Remove(baselinePath)
}
