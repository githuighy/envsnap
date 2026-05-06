package snapshot_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/user/envsnap/internal/snapshot"
)

func savePromoteSnap(t *testing.T, dir, name string, s snapshot.Snapshot) string {
	t.Helper()
	p := filepath.Join(dir, name)
	require.NoError(t, snapshot.Save(p, s))
	return p
}

func TestPromote_RoundTripFromDisk(t *testing.T) {
	dir := t.TempDir()
	srcPath := savePromoteSnap(t, dir, "src.json", snapshot.Snapshot{
		"APP_HOST": "prod.example.com",
		"APP_PORT": "443",
		"DB_PASS":  "secret",
	})
	dstPath := savePromoteSnap(t, dir, "dst.json", snapshot.Snapshot{
		"APP_HOST": "staging.example.com",
	})

	src, err := snapshot.Load(srcPath)
	require.NoError(t, err)
	dst, err := snapshot.Load(dstPath)
	require.NoError(t, err)

	_, promoted, err := snapshot.Promote(src, dst, snapshot.PromoteOptions{
		Prefixes:  []string{"APP_"},
		Overwrite: true,
	})
	require.NoError(t, err)
	require.NoError(t, snapshot.Save(dstPath, promoted))

	loaded, err := snapshot.Load(dstPath)
	require.NoError(t, err)
	assert.Equal(t, "prod.example.com", loaded["APP_HOST"])
	assert.Equal(t, "443", loaded["APP_PORT"])
	_, hasDB := loaded["DB_PASS"]
	assert.False(t, hasDB, "DB_PASS should not have been promoted")
}

func TestPromote_SummaryAccurate(t *testing.T) {
	src := snapshot.Snapshot{"A": "1", "B": "2", "C": "3"}
	dst := snapshot.Snapshot{"A": "old"}

	results, _, err := snapshot.Promote(src, dst, snapshot.PromoteOptions{Overwrite: false})
	require.NoError(t, err)

	summary := snapshot.PromoteSummary(results)
	assert.Contains(t, summary, "promoted=2")
	assert.Contains(t, summary, "skipped_exists=1")
}
