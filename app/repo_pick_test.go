package app

import (
	"claude-squad/config"
	"claude-squad/ui"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMultiRepoHome(t *testing.T) (*home, string) {
	root := t.TempDir()
	repo := filepath.Join(root, "shop")
	require.NoError(t, os.MkdirAll(filepath.Join(repo, ".git"), 0o755))
	sp := spinner.New()
	return &home{
		ctx:       context.Background(),
		state:     stateDefault,
		appConfig: config.DefaultConfig(),
		list:      ui.NewList(&sp, false),
		menu:      ui.NewMenu(),
		tabbedWindow: ui.NewTabbedWindow(
			ui.NewGridPane(&sp), ui.NewPreviewPane(), ui.NewDiffPane(), ui.NewTerminalPane()),
		repoRoot: root,
		keySent:  true, // skip the menu highlight round trip
	}, repo
}

func TestNewSessionOutsideRepoCreatesInstanceInPickedRepo(t *testing.T) {
	h, repo := newMultiRepoHome(t)

	h.handleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	assert.Equal(t, stateRepoPick, h.state)
	h.handleKeyPress(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, stateNew, h.state)
	require.Equal(t, 1, h.list.NumInstances())
	assert.Equal(t, repo, h.list.GetSelectedInstance().Path)
}

func TestRepoPickEscCreatesNoInstance(t *testing.T) {
	h, _ := newMultiRepoHome(t)

	h.handleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	h.handleKeyPress(tea.KeyMsg{Type: tea.KeyEsc})

	assert.Equal(t, stateDefault, h.state)
	assert.Equal(t, 0, h.list.NumInstances())
}

// press skips the menu highlight round trip that re-sends mapped keys.
func press(h *home, msg tea.KeyMsg) {
	h.keySent = true
	h.handleKeyPress(msg)
}

func TestSubmittingNameDoesNotRegisterRepoBeforeStart(t *testing.T) {
	h, _ := newMultiRepoHome(t)
	press(h, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	press(h, tea.KeyMsg{Type: tea.KeyEnter})
	calls := 0
	h.newInstanceFinalizer = func() { calls++ }

	press(h, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("ab")})
	press(h, tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, stateDefault, h.state)
	assert.Equal(t, 0, calls)
}
