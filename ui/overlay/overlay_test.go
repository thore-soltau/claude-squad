package overlay

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
)

func TestPlaceOverlayKeepsBackgroundHyperlinksIntact(t *testing.T) {
	link := ansi.SetHyperlink("https://github.com/ZooRoyal/claude-squad/pull/1234567890") + "PR" + ansi.ResetHyperlink()
	bgLine := "ab " + link + strings.Repeat(" right side", 4)
	bg := strings.Join([]string{bgLine, bgLine, bgLine}, "\n")
	plain := ansi.Strip(bgLine)

	out := PlaceOverlay(10, 1, "XX", bg, false, false)

	line := strings.Split(out, "\n")[1]
	assert.Equal(t, plain[:10]+"XX"+plain[12:], ansi.Strip(line))
	opened := strings.Count(line, "\x1b]8;") - strings.Count(line, ansi.ResetHyperlink())
	assert.Equal(t, opened, strings.Count(line, ansi.ResetHyperlink()), "every link must be closed")
}
