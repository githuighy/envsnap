package snapshot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var baseSrc = Snapshot{
	"APP_HOST":  "prod.example.com",
	"APP_PORT":  "443",
	"DB_PASS":   "secret",
	"LOG_LEVEL": "info",
}

var baseDst = Snapshot{
	"APP_HOST":  "staging.example.com",
	"LOG_LEVEL": "debug",
}

func TestPromote_AllKeys_Overwrite(t *testing.T) {
	_, out, err := Promote(baseSrc, baseDst, PromoteOptions{Overwrite: true})
	require.NoError(t, err)
	assert.Equal(t, "prod.example.com", out["APP_HOST"])
	assert.Equal(t, "443", out["APP_PORT"])
	assert.Equal(t, "secret", out["DB_PASS"])
	assert.Equal(t, "info", out["LOG_LEVEL"])
}

func TestPromote_NoOverwrite_ExistingPreserved(t *testing.T) {
	_, out, err := Promote(baseSrc, baseDst, PromoteOptions{Overwrite: false})
	require.NoError(t, err)
	assert.Equal(t, "staging.example.com", out["APP_HOST"], "existing key should not be overwritten")
	assert.Equal(t, "443", out["APP_PORT"], "new key should be added")
}

func TestPromote_PrefixFilter(t *testing.T) {
	_, out, err := Promote(baseSrc, baseDst, PromoteOptions{
		Prefixes:  []string{"APP_"},
		Overwrite: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "prod.example.com", out["APP_HOST"])
	assert.Equal(t, "443", out["APP_PORT"])
	_, hasDB := out["DB_PASS"]
	assert.False(t, hasDB, "DB_PASS should be skipped by prefix filter")
}

func TestPromote_ExcludeFilter(t *testing.T) {
	_, out, err := Promote(baseSrc, baseDst, PromoteOptions{
		Exclude:   []string{"DB_"},
		Overwrite: true,
	})
	require.NoError(t, err)
	_, hasDB := out["DB_PASS"]
	assert.False(t, hasDB, "DB_PASS should be excluded")
	assert.Equal(t, "prod.example.com", out["APP_HOST"])
}

func TestPromote_DryRun_DoesNotMutateDst(t *testing.T) {
	results, out, err := Promote(baseSrc, baseDst, PromoteOptions{Overwrite: true, DryRun: true})
	require.NoError(t, err)
	assert.Equal(t, baseDst, out, "dry run should not change destination")
	promoted := 0
	for _, r := range results {
		if r.Action == "promoted" {
			promoted++
		}
	}
	assert.Greater(t, promoted, 0, "should still report what would be promoted")
}

func TestPromoteSummary(t *testing.T) {
	results := []PromoteResult{
		{Key: "A", Action: "promoted"},
		{Key: "B", Action: "promoted"},
		{Key: "C", Action: "skipped_exists"},
		{Key: "D", Action: "skipped_filter"},
	}
	summary := PromoteSummary(results)
	assert.Equal(t, "promoted=2 skipped_exists=1 skipped_filter=1", summary)
}
