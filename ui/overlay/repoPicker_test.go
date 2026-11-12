package overlay

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

var testRepos = []string{"/dev/claude-squad", "/dev/zr-microservice-bundle", "/dev/zr-microservice-salesforce"}

func typeText(p *RepoPicker, text string) {
	p.HandleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(text)})
}

func TestRepoPickerEnterSelectsFirstRepo(t *testing.T) {
	p := NewRepoPicker(testRepos)

	selected, done := p.HandleKeyPress(tea.KeyMsg{Type: tea.KeyEnter})

	assert.True(t, done)
	assert.Equal(t, "/dev/claude-squad", selected)
}

func TestRepoPickerFilterIsCaseInsensitiveOnRepoName(t *testing.T) {
	p := NewRepoPicker(testRepos)

	typeText(p, "SALES")
	selected, done := p.HandleKeyPress(tea.KeyMsg{Type: tea.KeyEnter})

	assert.True(t, done)
	assert.Equal(t, "/dev/zr-microservice-salesforce", selected)
}

func TestRepoPickerFilterIgnoresParentPath(t *testing.T) {
	p := NewRepoPicker(testRepos)

	typeText(p, "dev")
	_, done := p.HandleKeyPress(tea.KeyMsg{Type: tea.KeyEnter})

	assert.False(t, done)
}

func TestRepoPickerDownMovesSelection(t *testing.T) {
	p := NewRepoPicker(testRepos)

	typeText(p, "zr")
	p.HandleKeyPress(tea.KeyMsg{Type: tea.KeyDown})
	selected, _ := p.HandleKeyPress(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, "/dev/zr-microservice-salesforce", selected)
}

func TestRepoPickerBackspaceWidensFilterAndResetsCursor(t *testing.T) {
	p := NewRepoPicker(testRepos)

	typeText(p, "zr")
	p.HandleKeyPress(tea.KeyMsg{Type: tea.KeyDown})
	p.HandleKeyPress(tea.KeyMsg{Type: tea.KeyBackspace})
	p.HandleKeyPress(tea.KeyMsg{Type: tea.KeyBackspace})
	selected, _ := p.HandleKeyPress(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, "/dev/claude-squad", selected)
}

func TestRepoPickerEscCancels(t *testing.T) {
	p := NewRepoPicker(testRepos)

	selected, done := p.HandleKeyPress(tea.KeyMsg{Type: tea.KeyEsc})

	assert.True(t, done)
	assert.Empty(t, selected)
}
