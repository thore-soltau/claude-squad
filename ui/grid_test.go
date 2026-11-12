package ui

import (
	"claude-squad/session"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func instanceActiveAt(title string, secondsAgo int) *session.Instance {
	return &session.Instance{Title: title, LastActivity: time.Now().Add(-time.Duration(secondsAgo) * time.Second)}
}

func TestAssignSlotsFillsEmptySlotsInListOrder(t *testing.T) {
	a, b := &session.Instance{Title: "a"}, &session.Instance{Title: "b"}

	slots := assignSlots([gridSlots]*session.Instance{}, []*session.Instance{a, b})

	assert.Equal(t, [gridSlots]*session.Instance{a, b, nil, nil}, slots)
}

func TestAssignSlotsKeepsSessionsInTheirSlot(t *testing.T) {
	a, b, c := instanceActiveAt("a", 30), instanceActiveAt("b", 20), instanceActiveAt("c", 10)
	current := [gridSlots]*session.Instance{c, nil, a, b}

	a.LastActivity = time.Now()
	slots := assignSlots(current, []*session.Instance{a, b, c})

	assert.Equal(t, current, slots)
}

func TestAssignSlotsReplacesLeastRecentlyActiveSession(t *testing.T) {
	a, b, c, d := instanceActiveAt("a", 40), instanceActiveAt("b", 10), instanceActiveAt("c", 20), instanceActiveAt("d", 30)
	e := instanceActiveAt("e", 0)
	current := [gridSlots]*session.Instance{a, b, c, d}

	slots := assignSlots(current, []*session.Instance{a, b, c, d, e})

	assert.Equal(t, [gridSlots]*session.Instance{e, b, c, d}, slots)
}

func TestAssignSlotsKeepsSlotsWhenNewcomerIsOlder(t *testing.T) {
	a, b, c, d := instanceActiveAt("a", 40), instanceActiveAt("b", 10), instanceActiveAt("c", 20), instanceActiveAt("d", 30)
	e := instanceActiveAt("e", 50)
	current := [gridSlots]*session.Instance{a, b, c, d}

	slots := assignSlots(current, []*session.Instance{a, b, c, d, e})

	assert.Equal(t, current, slots)
}

func TestAssignSlotsFreesSlotOfRemovedSession(t *testing.T) {
	a, b, c := instanceActiveAt("a", 30), instanceActiveAt("b", 20), instanceActiveAt("c", 10)
	current := [gridSlots]*session.Instance{a, b, nil, nil}

	slots := assignSlots(current, []*session.Instance{b, c})

	assert.Equal(t, [gridSlots]*session.Instance{c, b, nil, nil}, slots)
}

func TestGridInstanceReturnsSessionOfCell(t *testing.T) {
	a, b := &session.Instance{Title: "a"}, &session.Instance{Title: "b"}
	grid := &GridPane{slots: [gridSlots]*session.Instance{a, nil, b, nil}}
	window := NewTabbedWindow(grid, NewPreviewPane(), NewDiffPane(), NewTerminalPane())

	assert.Equal(t, b, window.GridInstance(2))
	assert.Nil(t, window.GridInstance(1))
	assert.Nil(t, window.GridInstance(gridSlots))
}

func TestGridInstanceIsNilOutsideGridTab(t *testing.T) {
	a := &session.Instance{Title: "a"}
	grid := &GridPane{slots: [gridSlots]*session.Instance{a}}
	window := NewTabbedWindow(grid, NewPreviewPane(), NewDiffPane(), NewTerminalPane())

	window.Toggle()

	assert.Nil(t, window.GridInstance(0))
}
