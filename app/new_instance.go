package app

import (
	"claude-squad/session"
	"claude-squad/session/git"
	"claude-squad/ui"
	"claude-squad/ui/overlay"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// createInstance adds an unnamed instance for repoPath and switches to name entry.
func (m *home) createInstance(repoPath string) (tea.Model, tea.Cmd) {
	instance, err := session.NewInstance(session.InstanceOptions{
		Title:   "",
		Path:    repoPath,
		Program: m.program,
	})
	if err != nil {
		return m, m.handleError(err)
	}

	m.newInstanceFinalizer = m.list.AddInstance(instance)
	m.list.SetSelectedInstance(m.list.NumInstances() - 1)
	m.state = stateNew
	m.menu.SetState(ui.StateNewInstance)

	return m, nil
}

// createInstanceWithPrompt is createInstance followed by the prompt overlay after naming.
func (m *home) createInstanceWithPrompt(repoPath string) (tea.Model, tea.Cmd) {
	model, cmd := m.createInstance(repoPath)
	if m.state != stateNew {
		return model, cmd
	}
	m.promptAfterName = true

	// Start a background fetch so branches are up to date by the time the picker opens
	return m, func() tea.Msg {
		git.FetchBranches(repoPath)
		return nil
	}
}

func (m *home) openRepoPicker(next func(repoPath string) (tea.Model, tea.Cmd)) (tea.Model, tea.Cmd) {
	repos, err := git.FindRepos(m.repoRoot)
	if err != nil {
		return m, m.handleError(err)
	}
	if len(repos) == 0 {
		return m, m.handleError(fmt.Errorf("no git repositories found in %s", m.repoRoot))
	}

	m.repoPicker = overlay.NewRepoPicker(repos)
	m.afterRepoPick = next
	m.state = stateRepoPick
	return m, nil
}

func (m *home) handleRepoPickState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	repoPath, done := m.repoPicker.HandleKeyPress(msg)
	if !done {
		return m, nil
	}

	next := m.afterRepoPick
	m.repoPicker = nil
	m.afterRepoPick = nil
	m.state = stateDefault
	if repoPath == "" {
		return m, nil
	}
	return next(repoPath)
}
