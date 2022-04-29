package rolelist

import (
	"github.com/charmbracelet/lipgloss"
)

// https://www.nordtheme.com/docs/colors-and-palettes
var (
	// Frost
	nord7  = lipgloss.Color("#8FBCBB")
	nord8  = lipgloss.Color("#88C0D0")
	nord10 = lipgloss.Color("#5E81AC")

	// Aurora
	nord12 = lipgloss.Color("#D08770") // orange
	nord14 = lipgloss.Color("#A3BE8C") // green
)

type Styles struct {
	white lipgloss.Color
	black lipgloss.Color
	grey  lipgloss.Color

	listStyle         lipgloss.Style
	InactiveListStyle lipgloss.Style
	ActiveListStyle   lipgloss.Style

	focusedStyle      lipgloss.Style
	noStyle           lipgloss.Style
	blurredStyle      lipgloss.Style
	addRoleInputStyle lipgloss.Style
}

func DefaultStyles() (s Styles) {

	s.white = lipgloss.Color("15")
	s.black = lipgloss.Color("16")
	s.grey = lipgloss.Color("240")

	s.InactiveListStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.NormalBorder())

	s.ActiveListStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.NormalBorder())

	s.listStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.NormalBorder())

	s.focusedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{
			Light: string(nord12),
			Dark:  string(nord12),
		})

	s.noStyle = lipgloss.NewStyle()

	s.blurredStyle = lipgloss.NewStyle().Foreground(s.grey)

	s.addRoleInputStyle = lipgloss.NewStyle().PaddingTop(1)

	return s
}
