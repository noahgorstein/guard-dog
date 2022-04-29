package roledetails

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// https://www.nordtheme.com/docs/colors-and-palettes
var (
	// Polar Night
	nord0 = lipgloss.Color("#2E3440")
	nord3 = lipgloss.Color("#4C566A")

	// Snow storm
	nord6 = lipgloss.Color("#ECEFF4")

	// Frost
	nord7  = lipgloss.Color("#8FBCBB")
	nord8  = lipgloss.Color("#88C0D0")
	nord9  = lipgloss.Color("#81A1C1")
	nord10 = lipgloss.Color("#5E81AC")

	// Aurora
	nord11 = lipgloss.Color("#BF616A") // red
	nord12 = lipgloss.Color("#D08770") // orange
	nord13 = lipgloss.Color("#EBCB8B") // yellow
	nord14 = lipgloss.Color("#A3BE8C") // green
	nord15 = lipgloss.Color("#B48EAD") // purple
)

type Styles struct {
	white lipgloss.Color
	black lipgloss.Color
	grey  lipgloss.Color

	InactiveViewportStyle lipgloss.Style
	ActiveViewportStyle   lipgloss.Style
	textStyle             lipgloss.Style
	emphasizedTextStyle   lipgloss.Style
	helpTextKey           lipgloss.Style
	helpTextDescription   lipgloss.Style

	focusedStyle    lipgloss.Style
	noStyle         lipgloss.Style
	blurredStyle    lipgloss.Style
	header          lipgloss.Style
	secondaryHeader lipgloss.Style

	PermissionsTableBorder      table.Border
	PermissionsTableHeaderStyle lipgloss.Style
}

func DefaultStyles() (s Styles) {

	s.white = lipgloss.Color("15")
	s.black = lipgloss.Color("16")
	s.grey = lipgloss.Color("240")

	s.InactiveViewportStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Border(lipgloss.NormalBorder())

	s.ActiveViewportStyle = s.InactiveViewportStyle.Copy()

	s.textStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: string(nord0),
		Dark:  string(nord6),
	})
	s.emphasizedTextStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{
		Light: string(nord10),
		Dark:  string(nord8),
	})

	s.helpTextKey = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{
		Light: string(nord10),
		Dark:  string(nord8),
	})
	s.helpTextDescription = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{
		Light: string(nord0),
		Dark:  string(nord6),
	})

	s.focusedStyle = lipgloss.NewStyle().Foreground(nord12)

	s.noStyle = lipgloss.NewStyle()

	s.blurredStyle = lipgloss.NewStyle().Foreground(s.grey)

	s.header = lipgloss.NewStyle().
		MarginBottom(1).
		Background(nord12).
		Foreground(nord6).
		Bold(true).
		Italic(true).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.AdaptiveColor{
			Light: string(nord3),
			Dark:  string(nord6),
		})

	s.secondaryHeader = lipgloss.NewStyle().
		MarginBottom(1).
		Background(nord9).
		Foreground(nord6).
		Bold(true).
		Italic(true).
		Padding(0, 1)

	s.PermissionsTableHeaderStyle = lipgloss.NewStyle().Bold(true)

	s.PermissionsTableBorder = table.Border{
		Top:    "─",
		Left:   "│",
		Right:  "│",
		Bottom: "─",

		TopRight:    "╮",
		TopLeft:     "╭",
		BottomRight: "╯",
		BottomLeft:  "╰",

		TopJunction:    "┬",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "┴",
		InnerJunction:  "┼",

		InnerDivider: "│",
	}

	return s
}
