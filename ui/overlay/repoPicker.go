package overlay

import (
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const repoPickerMaxVisible = 10

// RepoPicker lets the user choose a git repository by typing a filter.
type RepoPicker struct {
	repos  []string
	filter string
	cursor int
	width  int
}

func NewRepoPicker(repos []string) *RepoPicker {
	return &RepoPicker{repos: repos, width: 60}
}

func (p *RepoPicker) SetWidth(width int) {
	p.width = width
}

// HandleKeyPress returns done once the picker should close. selected is empty on cancel.
func (p *RepoPicker) HandleKeyPress(msg tea.KeyMsg) (selected string, done bool) {
	items := p.visibleRepos()
	switch msg.Type {
	case tea.KeyEsc:
		return "", true
	case tea.KeyEnter:
		if len(items) == 0 {
			return "", false
		}
		return items[p.cursor], true
	case tea.KeyUp:
		if p.cursor > 0 {
			p.cursor--
		}
	case tea.KeyDown:
		if p.cursor < len(items)-1 {
			p.cursor++
		}
	case tea.KeyBackspace:
		if runes := []rune(p.filter); len(runes) > 0 {
			p.filter = string(runes[:len(runes)-1])
			p.cursor = 0
		}
	case tea.KeyRunes:
		p.filter += string(msg.Runes)
		p.cursor = 0
	}
	return "", false
}

// visibleRepos matches the filter against the repo name only, every repo shares the parent path.
func (p *RepoPicker) visibleRepos() []string {
	filter := strings.ToLower(p.filter)
	var items []string
	for _, repo := range p.repos {
		if strings.Contains(strings.ToLower(filepath.Base(repo)), filter) {
			items = append(items, repo)
		}
	}
	return items
}

func (p *RepoPicker) Render() string {
	var s strings.Builder
	s.WriteString(bpLabelStyle.Render("Repository"))
	s.WriteString(bpFilterStyle.Render(" (filter: " + p.filter + "█)"))
	s.WriteString("\n\n")

	items := p.visibleRepos()
	if len(items) == 0 {
		s.WriteString(bpDimStyle.Render("  No matching repositories"))
	}
	start := max(0, p.cursor-repoPickerMaxVisible+1)
	end := min(len(items), start+repoPickerMaxVisible)
	for i := start; i < end; i++ {
		label := filepath.Base(items[i])
		if i == p.cursor {
			s.WriteString(bpSelectedStyle.Render("> " + label))
		} else {
			s.WriteString(bpDimStyle.Render("  " + label))
		}
		s.WriteString("\n")
	}
	s.WriteString("\n" + bpDimStyle.Render("↑/↓ select • enter confirm • esc cancel"))

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(p.width).
		Render(s.String())
}
