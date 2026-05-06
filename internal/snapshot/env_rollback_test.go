package snapshot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var baseRollbackCurrent = Snapshot{Vars: map[string]string{
	"APP_HOST": "prod.example.com",
	"APP_PORT": "9090",
	"DB_URL":   "postgres://prod/db",
	"NEW_KEY":  "only-in-current",
}}

var baseRollbackBaseline = Snapshot{Vars: map[string]string{
	"APP_HOST": "staging.example.com",
	"APP_PORT": "8080",
	"DB_URL":   "postgres://prod/db", // same as current
}}

func TestRollback_AllKeys(t *testing.T) {
	res, err := Rollback(baseRollbackCurrent, baseRollbackBaseline, RollbackOptions{})
	require.NoError(t, err)
	assert.Equal(t, "staging.example.com", res.Snapshot.Vars["APP_HOST"])
	assert.Equal(t, "8080", res.Snapshot.Vars["APP_PORT"])
	assert.Equal(t, map[string]string{"APP_HOST": "staging.example.com", "APP_PORT": "8080"}, res.Restored)
	assert.Contains(t, res.Dropped, "NEW_KEY")
	assert.NotContains(t, res.Snapshot.Vars, "NEW_KEY")
	assert.Contains(t, res.Unchanged, "DB_URL")
}

func TestRollback_ExplicitKeys(t *testing.T) {
	res, err := Rollback(baseRollbackCurrent, baseRollbackBaseline, RollbackOptions{
		Keys: []string{"APP_HOST"},
	})
	require.NoError(t, err)
	assert.Equal(t, "staging.example.com", res.Snapshot.Vars["APP_HOST"])
	// APP_PORT should remain unchanged from current
	assert.Equal(t, "9090", res.Snapshot.Vars["APP_PORT"])
	assert.Equal(t, map[string]string{"APP_HOST": "staging.example.com"}, res.Restored)
}

func TestRollback_ByPrefix(t *testing.T) {
	res, err := Rollback(baseRollbackCurrent, baseRollbackBaseline, RollbackOptions{
		Prefixes: []string{"APP_"},
	})
	require.NoError(t, err)
	assert.Equal(t, "staging.example.com", res.Snapshot.Vars["APP_HOST"])
	assert.Equal(t, "8080", res.Snapshot.Vars["APP_PORT"])
	// DB_URL and NEW_KEY should be untouched
	assert.Equal(t, "postgres://prod/db", res.Snapshot.Vars["DB_URL"])
	assert.Equal(t, "only-in-current", res.Snapshot.Vars["NEW_KEY"])
}

func TestRollback_DryRun_DoesNotMutate(t *testing.T) {
	res, err := Rollback(baseRollbackCurrent, baseRollbackBaseline, RollbackOptions{DryRun: true})
	require.NoError(t, err)
	// DryRun: snapshot should still reflect current values
	assert.Equal(t, "prod.example.com", res.Snapshot.Vars["APP_HOST"])
	assert.Contains(t, res.Snapshot.Vars, "NEW_KEY")
	// But Restored and Dropped still report what would happen
	assert.NotEmpty(t, res.Restored)
	assert.Contains(t, res.Dropped, "NEW_KEY")
}

func TestRollback_NilCurrentReturnsError(t *testing.T) {
	_, err := Rollback(Snapshot{}, baseRollbackBaseline, RollbackOptions{})
	require.Error(t, err)
}

func TestRollback_NilBaselineReturnsError(t *testing.T) {
	_, err := Rollback(baseRollbackCurrent, Snapshot{}, RollbackOptions{})
	require.Error(t, err)
}
