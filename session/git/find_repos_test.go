package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindReposReturnsDirectSubdirectoriesWithGit(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "repo", ".git"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "worktree"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "worktree", ".git"), []byte("gitdir: x"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "plain"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "group", "nested", ".git"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "file.txt"), nil, 0o644))

	repos, err := FindRepos(root)

	require.NoError(t, err)
	assert.Equal(t, []string{filepath.Join(root, "repo"), filepath.Join(root, "worktree")}, repos)
}

func TestFindReposFailsForMissingRoot(t *testing.T) {
	_, err := FindRepos(filepath.Join(t.TempDir(), "missing"))

	assert.Error(t, err)
}
