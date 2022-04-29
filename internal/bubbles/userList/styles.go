package userlist

import (
	"github.com/charmbracelet/lipgloss"
)

// https://www.nordtheme.com/docs/colors-and-palettes
var (
	// Polar Night
	nord0 = lipgloss.Color("#2E3440")

	// Snow storm
	nord4 = lipgloss.Color("#D8DEE9")
	nord6 = lipgloss.Color("#ECEFF4")

	// Frost
	nord8  = lipgloss.Color("#88C0D0")
	nord10 = lipgloss.Color("#5E81AC")

	// Aurora
	nord11 = lipgloss.Color("#BF616A") // red
	nord12 = lipgloss.Color("#D08770") // orange
	nord14 = lipgloss.Color("#A3BE8C") // green
	nord15 = lipgloss.Color("#B48EAD") // purple
)

type Styles struct {
	grey lipgloss.Color

	focusedStyle      lipgloss.Style
	noStyle           lipgloss.Style
	blurredStyle      lipgloss.Style
	addUserInputStyle lipgloss.Style
	ActiveListStyle   lipgloss.Style
	InactiveListStyle lipgloss.Style
	listStyle         lipgloss.Style
}

func DefaultStyles() (s Styles) {

	s.grey = lipgloss.Color("240")

	s.focusedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{
			Light: string(nord12),
			Dark:  string(nord12),
		})

	s.noStyle = lipgloss.NewStyle()

	s.blurredStyle = lipgloss.NewStyle().Foreground(s.grey)

	s.addUserInputStyle = lipgloss.NewStyle().PaddingTop(1)

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

	return s
}
