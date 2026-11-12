package ui

import (
	"claude-squad/session"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

const gridSlots = 4

const ansiReset = "\x1b[0m"

var (
	gridCellStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#BBBBBB", Dark: "#555555"})
	selectedGridCellStyle = gridCellStyle.BorderForeground(highlightColor)
)

// GridPane shows live snapshots of the most recently active sessions in a 2x2 grid.
type GridPane struct {
	width  int
	height int

	spinner  *spinner.Model
	slots    [gridSlots]*session.Instance
	content  [gridSlots]string
	selected *session.Instance
}

func NewGridPane(spinner *spinner.Model) *GridPane {
	return &GridPane{spinner: spinner}
}

func (g *GridPane) SetSize(width, height int) {
	g.width = width
	g.height = height
}

// CellContentSize is the tmux size that makes a session fit one cell (minus border and title row).
func (g *GridPane) CellContentSize() (width, height int) {
	return g.width/2 - gridCellStyle.GetHorizontalFrameSize(),
		g.height/2 - gridCellStyle.GetVerticalFrameSize() - 1
}

// UpdateContent reassigns the slots and captures the pane output of every slotted session.
func (g *GridPane) UpdateContent(instances []*session.Instance, selected *session.Instance) error {
	var active []*session.Instance
	for _, inst := range instances {
		if inst.Started() && !inst.Paused() {
			active = append(active, inst)
		}
	}
	g.slots = assignSlots(g.slots, active)
	g.selected = selected

	for i, inst := range g.slots {
		g.content[i] = ""
		if inst == nil {
			continue
		}
		content, err := inst.Preview()
		if err != nil {
			return err
		}
		g.content[i] = content
	}
	return nil
}

// Instance returns the session in the given cell (0-based), or nil.
func (g *GridPane) Instance(cell int) *session.Instance {
	if cell < 0 || cell >= gridSlots {
		return nil
	}
	return g.slots[cell]
}

func (g *GridPane) String() string {
	if g.width == 0 || g.height == 0 {
		return ""
	}
	var cells [gridSlots]string
	for i := range g.slots {
		cells[i] = g.renderCell(i)
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, cells[0], cells[1]),
		lipgloss.JoinHorizontal(lipgloss.Top, cells[2], cells[3]),
	)
}

func (g *GridPane) renderCell(i int) string {
	width, height := g.CellContentSize()
	style := gridCellStyle.Width(width).Height(height + 1)
	inst := g.slots[i]
	if inst == nil {
		return style.Render("")
	}
	if inst == g.selected {
		style = selectedGridCellStyle.Width(width).Height(height + 1)
	}

	title := truncate.StringWithTail(g.statusIcon(inst)+inst.Title, uint(width), "...")
	lines := strings.Split(strings.TrimRight(g.content[i], "\n"), "\n")
	if len(lines) > height {
		lines = lines[len(lines)-height:]
	}
	for j, line := range lines {
		// Reset so colors cut off mid-line don't bleed into the cell border.
		lines[j] = truncate.String(line, uint(width)) + ansiReset
	}
	return style.Render(title + "\n" + strings.Join(lines, "\n"))
}

func (g *GridPane) statusIcon(inst *session.Instance) string {
	switch inst.Status {
	case session.Running, session.Loading:
		return g.spinner.View() + " "
	case session.Ready:
		return readyStyle.Render(readyIcon)
	default:
		return ""
	}
}

// assignSlots keeps sessions in their slot so the grid doesn't jump, and lets a
// newly active session evict the one that has been quiet the longest.
func assignSlots(current [gridSlots]*session.Instance, candidates []*session.Instance) [gridSlots]*session.Instance {
	alive := make(map[*session.Instance]bool, len(candidates))
	for _, inst := range candidates {
		alive[inst] = true
	}

	slots := current
	slotted := make(map[*session.Instance]bool, gridSlots)
	for i, inst := range slots {
		if inst != nil && !alive[inst] {
			slots[i] = nil
		} else if inst != nil {
			slotted[inst] = true
		}
	}

	var waiting []*session.Instance
	for _, inst := range candidates {
		if !slotted[inst] {
			waiting = append(waiting, inst)
		}
	}
	sort.SliceStable(waiting, func(a, b int) bool {
		return waiting[a].LastActivity.After(waiting[b].LastActivity)
	})

	for _, inst := range waiting {
		i := stalestSlot(slots)
		if slots[i] != nil && !inst.LastActivity.After(slots[i].LastActivity) {
			break
		}
		slots[i] = inst
	}
	return slots
}

// stalestSlot returns the first empty slot, otherwise the one with the oldest activity.
func stalestSlot(slots [gridSlots]*session.Instance) int {
	stalest := 0
	for i, inst := range slots {
		if inst == nil {
			return i
		}
		if inst.LastActivity.Before(slots[stalest].LastActivity) {
			stalest = i
		}
	}
	return stalest
}
