package addpermissionprompt

import (
	"github.com/charmbracelet/lipgloss"
)

// https://www.nordtheme.com/docs/colors-and-palettes
var (
	// Polar Night
	nord0 = lipgloss.Color("#2E3440")

	// Snow storm
	nord6 = lipgloss.Color("#ECEFF4")

	// Frost
	nord8  = lipgloss.Color("#88C0D0")
	nord9  = lipgloss.Color("#81A1C1")
	nord10 = lipgloss.Color("#5E81AC")

	// Aurora
	nord12 = lipgloss.Color("#D08770") // orange
	nord15 = lipgloss.Color("#B48EAD") // purple
)

type Styles struct {
	white lipgloss.Color
	black lipgloss.Color
	grey  lipgloss.Color

	focusedStyle    lipgloss.Style
	noStyle         lipgloss.Style
	blurredStyle    lipgloss.Style
	secondaryHeader lipgloss.Style

	selectionTextStyle     lipgloss.Style
	emphasisTextStyle      lipgloss.Style
	normalTextStyle        lipgloss.Style
	activeResourceInputs   lipgloss.Style
	inactiveResourceInputs lipgloss.Style
	resourceInputsPrompt   lipgloss.Style
}

func DefaultStyles() (s Styles) {
	s.white = lipgloss.Color("15")
	s.black = lipgloss.Color("16")
	s.grey = lipgloss.Color("240")

	s.focusedStyle = lipgloss.NewStyle().Foreground(nord12)

	s.noStyle = lipgloss.NewStyle()

	s.blurredStyle = lipgloss.NewStyle().Foreground(s.grey)

	s.secondaryHeader = lipgloss.NewStyle().
		MarginBottom(1).
		Background(nord9).
		Foreground(lipgloss.AdaptiveColor{
			Light: string(nord6),
			Dark:  string(nord6),
		}).
		Bold(true).
		Italic(true).
		Padding(0, 1)

	s.emphasisTextStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{
		Light: string(nord10),
		Dark:  string(nord8),
	})
	s.normalTextStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{
		Light: string(nord0),
		Dark:  string(nord6),
	})

	s.selectionTextStyle = lipgloss.NewStyle().Bold(true).Foreground(nord15)

	s.activeResourceInputs = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(nord12).
		Padding(1, 1)

	s.inactiveResourceInputs = s.activeResourceInputs.Copy().
		BorderForeground(s.grey)

	s.resourceInputsPrompt = lipgloss.NewStyle().
		Bold(true).
		MarginBottom(1)

	return s
}
