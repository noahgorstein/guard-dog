package statusbar

import "github.com/charmbracelet/lipgloss"

// https://www.nordtheme.com/docs/colors-and-palettes
var (
	// Polar Night
	nord0 = lipgloss.Color("#2E3440")

	// Snow storm
	nord6 = lipgloss.Color("#ECEFF4")

	// Frost
	nord8  = lipgloss.Color("#88C0D0")
	nord10 = lipgloss.Color("#5E81AC")

	// Aurora
	nord11 = lipgloss.Color("#BF616A") // red
	nord12 = lipgloss.Color("#D08770") // orange
	nord14 = lipgloss.Color("#A3BE8C") // green
)

var (
	grey = lipgloss.Color("240")
)

type Styles struct {
	StatusBarStyle      lipgloss.Style
	EndpointStyle       lipgloss.Style
	HelpKeyStyle        lipgloss.Style
	HelpTextStyle       lipgloss.Style
	ErrorMessageStyle   lipgloss.Style
	SuccessMessageStyle lipgloss.Style
}

func DefaultStyles() (s Styles) {

	s.StatusBarStyle = lipgloss.NewStyle()

	s.EndpointStyle = lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{
			Light: string(nord11),
			Dark:  string(nord11),
		}).
		Padding(0, 1).
		Foreground(nord6).
		Bold(true)

	s.HelpKeyStyle = lipgloss.NewStyle().Bold(true).Foreground(
		lipgloss.AdaptiveColor{
			Light: string(nord10),
			Dark:  string(nord8),
		},
	)

	s.HelpTextStyle = lipgloss.NewStyle().Foreground(
		lipgloss.AdaptiveColor{
			Light: string(nord0),
			Dark:  string(nord6),
		},
	)

	s.ErrorMessageStyle = lipgloss.NewStyle().Foreground(
		lipgloss.AdaptiveColor{
			Light: string(nord11),
			Dark:  string(nord11),
		},
	)

	s.SuccessMessageStyle = lipgloss.NewStyle().Foreground(
		lipgloss.AdaptiveColor{
			Light: string(nord14),
			Dark:  string(nord14),
		},
	)

	return s
}
