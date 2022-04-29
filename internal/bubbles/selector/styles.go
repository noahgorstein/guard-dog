package selector

import (
	"github.com/charmbracelet/lipgloss"
)

// https://www.nordtheme.com/docs/colors-and-palettes
var (
	// Aurora
	nord12 = lipgloss.Color("#D08770") // orange
	nord14 = lipgloss.Color("#A3BE8C") // green

)

type Styles struct {
	grey lipgloss.Color

	CheckmarkStyle      lipgloss.Style
	CheckedChoiceStyle  lipgloss.Style
	SelectedChoiceStyle lipgloss.Style

	PromptStyle   lipgloss.Style
	ActiveStyle   lipgloss.Style
	InactiveStyle lipgloss.Style
}

func DefaultStyles() (s Styles) {
	s.grey = lipgloss.Color("240")

	s.SelectedChoiceStyle = lipgloss.NewStyle().Foreground(nord12)
	s.CheckmarkStyle = lipgloss.NewStyle().Foreground(nord14).Bold(true)
	s.CheckedChoiceStyle = lipgloss.NewStyle().Foreground(nord14).Bold(true)
	s.PromptStyle = lipgloss.NewStyle().Bold(true)

	s.ActiveStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(nord12).
		Padding(0, 1)

	s.InactiveStyle = s.ActiveStyle.Copy().
		BorderForeground(s.grey)

	return s
}
