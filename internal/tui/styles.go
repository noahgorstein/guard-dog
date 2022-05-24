package tui

import "github.com/charmbracelet/lipgloss"

var (
	listStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("202")).
			Border(lipgloss.NormalBorder(), true).
			MarginRight(1)

	listTitleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("202")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Bold(true).
			Italic(true)

	inactiveListStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), true).
				MarginRight(1)

	viewportStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("202"))

	inactiveViewportStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder())

	usernameStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#004B60")).
			Margin(1, 0, 1, 1).
			Padding(1, 2)

	permissionsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#004B60")).
				Margin(1, 0, 1, 1).
				Padding(1, 1)

	permissionsTableStyle = lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("#FAFAFA")).
				Foreground(lipgloss.Color("#FAFAFA")).
				Align(lipgloss.Left)

	permissionsTableHighltedRowStyle = lipgloss.NewStyle().
						Background(lipgloss.Color("#FAFAFA")).
						Foreground(lipgloss.Color("202"))

	statusBarStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true)

	usernameStatusBarStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Margin(1, 1)

	endpointStatusBarStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("202")).Background(lipgloss.Color("#FAFAFA"))
)
