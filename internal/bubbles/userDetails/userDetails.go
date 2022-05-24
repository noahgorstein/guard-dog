package userdetails

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/noahgorstein/stardog-go/internal/config"
	"github.com/noahgorstein/stardog-go/stardog"
)

var (
	f, _ = tea.LogToFile("debug.log", "debug")
)

var (
	customBorder = table.Border{
		Top:    "─",
		Left:   "│",
		Right:  "│",
		Bottom: "─",

		TopRight:    "╮",
		TopLeft:     "╭",
		BottomRight: "╯",
		BottomLeft:  "╰",

		TopJunction:    "╥",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "╨",
		InnerJunction:  "╫",

		InnerDivider: "║",
	}
)

type Bubble struct {
	viewport         viewport.Model
	active           bool
	permissionsTable table.Model
	connection       stardog.ConnectionDetails
	selectedUser     stardog.User
}

const (
	padding = 1
)

const (
	columnKeyAction       = "action"
	columnKeyResourceType = "resource type"
	columnKeyResource     = "resource"
	columnKeyExplicit     = "explicit"
)

func (b *Bubble) SetCurrentUser(user stardog.User) {
	b.selectedUser = user
}

func New(config config.Config) Bubble {
	connectionDetails := stardog.NewConnectionDetails(config.Endpoint, config.Username, config.Password)
	border := lipgloss.RoundedBorder()

	viewport := viewport.New(0, 0)
	viewport.Style = lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border)
	viewport.SetContent(generateContent(0, 0, ""))

	permissionsTable := table.New([]table.Column{
		table.NewColumn(columnKeyAction, "Action", 10).WithFiltered(true),
		table.NewColumn(columnKeyResourceType, "Resource Type", 20),
		table.NewFlexColumn(columnKeyResource, "Resource", 1),
		table.NewColumn(columnKeyExplicit, "Explicit Permission", 20),
	}).Focused(true).Filtered(true).Border(customBorder)
	permissionsTable = permissionsTable.HeaderStyle(lipgloss.NewStyle().Bold(true))

	return Bubble{
		viewport:         viewport,
		connection:       *connectionDetails,
		permissionsTable: permissionsTable,
	}

}

func (b Bubble) Init() tea.Cmd {
	return b.GetUserPermissionsCmd(b.selectedUser)
}

func generateContent(width, height int, content string) string {
	header := lipgloss.NewStyle().
		MarginBottom(1).
		Background(lipgloss.Color("202")).
		Bold(true).
		Italic(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("#FAFAFA")).
		Render("Permissions")

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Render(header + content)
}

func (b *Bubble) SetIsActive(active bool) {
	b.active = active
}

func (b *Bubble) SetSize(width, height int) {
	b.viewport.Width = width - b.viewport.Style.GetHorizontalFrameSize()
	b.viewport.Height = height - b.viewport.Style.GetVerticalFrameSize()

	b.viewport.SetContent(
		generateContent(
			b.viewport.Width,
			b.viewport.Height,
			b.permissionsTable.WithTargetWidth(b.viewport.Width).View()))
}

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	// f.WriteString("msg recieved! \n")
	// f.WriteString("current user is: " + b.selectedUser.Name + "\n")

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case GetUserPermissionsMsg:
		if msg != nil {
			b.permissionsTable = b.permissionsTable.WithRows(msg).WithTargetWidth(b.viewport.Width).WithHighlightedRow(0)
			b.viewport.SetContent(
				generateContent(
					b.viewport.Width,
					b.viewport.Height,
					b.permissionsTable.View()))
		}
	case DeleteUserPermissionMsg:
		cmds = append(cmds, b.GetUserPermissionsCmd(b.selectedUser))

	case tea.KeyMsg:

		if msg.String() == "x" {
			cmds = append(cmds, b.DeleteUserPermissionCmd(b.selectedUser))
		}

		b.permissionsTable, cmd = b.permissionsTable.Update(msg)
		cmds = append(cmds, cmd)

		b.viewport, cmd = b.viewport.Update(msg)
		cmds = append(cmds, cmd)

		b.viewport.SetContent(
			generateContent(
				b.viewport.Width,
				b.viewport.Height,
				b.permissionsTable.View()))

		return b, tea.Batch(cmds...)

	}

	return b, tea.Batch(cmds...)

}

func (b Bubble) View() string {
	if b.active {
		b.viewport.Style = b.viewport.Style.Copy().BorderForeground(lipgloss.Color("33"))
	} else {
		b.viewport.Style = b.viewport.Style.Copy().BorderForeground(lipgloss.Color("#FAFAFA"))
	}
	return b.viewport.View()
}
