package userdetails

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// https://www.nordtheme.com/docs/colors-and-palettes
var (
	// Polar Night
	nord0 = lipgloss.Color("#2E3440")

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
	grey lipgloss.Color

	InactiveViewportStyle lipgloss.Style
	ActiveViewportStyle   lipgloss.Style

	header                      lipgloss.Style
	secondaryHeader             lipgloss.Style
	roleStyle                   lipgloss.Style
	permissionsTableHeaderStyle lipgloss.Style
	helpTextKey                 lipgloss.Style
	helpTextDescription         lipgloss.Style

	permissionsTableBorder table.Border
}

func DefaultStyles() (s Styles) {

	s.InactiveViewportStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Border(lipgloss.NormalBorder())

	s.ActiveViewportStyle = s.InactiveViewportStyle.Copy()

	s.header = lipgloss.NewStyle().
		MarginBottom(1).
		Background(nord12).
		Foreground(lipgloss.AdaptiveColor{
			Light: string(nord6),
			Dark:  string(nord6),
		}).
		Bold(true).
		Italic(true).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.AdaptiveColor{
			Light: string(nord0),
			Dark:  string(nord6),
		})

	s.secondaryHeader = lipgloss.NewStyle().
		MarginBottom(1).
		Background(nord9).
		Foreground(nord6).
		Bold(true).
		Italic(true).
		Padding(0, 1)

	s.roleStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Background(nord15).
		Foreground(nord6)

	s.permissionsTableHeaderStyle = lipgloss.NewStyle().Bold(true)

	s.helpTextKey = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{
		Light: string(nord10),
		Dark:  string(nord8),
	})
	s.helpTextDescription = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{
		Light: string(nord0),
		Dark:  string(nord6),
	})

	s.permissionsTableBorder = table.Border{
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
