package userdetails

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
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
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	noStyle      = lipgloss.NewStyle()
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	addPermission = lipgloss.NewStyle().Background(lipgloss.Color("202")).
			Foreground(lipgloss.Color("#FAFAFA")).Bold(true).Render("Add Permission")
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

type sessionState int

const (
	idleState sessionState = iota
	addUserPermissionState
)

type Bubble struct {
	focusIndex              int
	state                   sessionState
	viewport                viewport.Model
	active                  bool
	permissionsTable        table.Model
	connection              stardog.ConnectionDetails
	selectedUser            stardog.User
	addUserPermissionInputs []textinput.Model
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
	viewport.SetContent("")

	permissionsTable := table.New([]table.Column{
		table.NewColumn(columnKeyAction, "Action", 10).WithFiltered(true),
		table.NewColumn(columnKeyResourceType, "Resource Type", 20),
		table.NewFlexColumn(columnKeyResource, "Resource", 1),
		table.NewColumn(columnKeyExplicit, "Explicit Permission", 20),
	}).Focused(true).Filtered(true).Border(customBorder)
	permissionsTable = permissionsTable.HeaderStyle(lipgloss.NewStyle().Bold(true))

	b := Bubble{
		state:                   idleState,
		viewport:                viewport,
		connection:              *connectionDetails,
		permissionsTable:        permissionsTable,
		addUserPermissionInputs: make([]textinput.Model, 3),
	}

	var addUserPermissionInput textinput.Model
	for i := range b.addUserPermissionInputs {
		addUserPermissionInput = textinput.New()
		addUserPermissionInput.CharLimit = 32

		switch i {
		case 0:
			addUserPermissionInput.Placeholder = "Action"
			addUserPermissionInput.Focus()
			addUserPermissionInput.PromptStyle = focusedStyle
			addUserPermissionInput.TextStyle = focusedStyle
		case 1:
			addUserPermissionInput.Placeholder = "Resource Type"
		case 2:
			addUserPermissionInput.Placeholder = "Resource Type"
		}
		b.addUserPermissionInputs[i] = addUserPermissionInput
	}

	return b
}

func (b Bubble) Init() tea.Cmd {
	return b.GetUserPermissionsCmd(b.selectedUser)
}

func (b Bubble) generateContent(width, height int, content string) string {
	header := lipgloss.NewStyle().
		MarginBottom(1).
		Background(lipgloss.Color("202")).
		Bold(true).
		Italic(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("#FAFAFA")).
		Render("Permissions")

	var builder strings.Builder
	var inputView string

	switch b.state {
	case addUserPermissionState:
		builder.WriteRune('\n')
		builder.WriteString(addPermission + "\n")
		for i := range b.addUserPermissionInputs {
			builder.WriteString(b.addUserPermissionInputs[i].View())
			if i < len(b.addUserPermissionInputs)-1 {
				builder.WriteRune('\n')

			}
		}
		button := &blurredButton
		if b.focusIndex == len(b.addUserPermissionInputs) {
			button = &focusedButton
		}
		builder.WriteString("\n" + *button)

		inputView = builder.String()
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Render(header + content + inputView)
}

func (b *Bubble) SetIsActive(active bool) {
	b.active = active
}

func (b *Bubble) SetSize(width, height int) {
	b.viewport.Width = width - b.viewport.Style.GetHorizontalFrameSize()
	b.viewport.Height = height - b.viewport.Style.GetVerticalFrameSize()

	b.viewport.SetContent(
		b.generateContent(
			b.viewport.Width,
			b.viewport.Height,
			b.permissionsTable.WithTargetWidth(b.viewport.Width).View()))
}

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case GetUserPermissionsMsg:
		if msg != nil {
			b.permissionsTable = b.permissionsTable.WithRows(msg).WithTargetWidth(b.viewport.Width).WithHighlightedRow(0)
			b.viewport.SetContent(
				b.generateContent(
					b.viewport.Width,
					b.viewport.Height,
					b.permissionsTable.View()))
		}
	case DeleteUserPermissionMsg:
		cmds = append(cmds, b.GetUserPermissionsCmd(b.selectedUser))

	case tea.KeyMsg:

		if msg.Type == tea.KeyCtrlA {
			action := b.addUserPermissionInputs[0]
			action.Focus()
			action.PromptStyle = focusedStyle
			b.state = addUserPermissionState
			b.viewport.SetContent(
				b.generateContent(
					b.viewport.Width,
					b.viewport.Height,
					b.permissionsTable.View()))
			return b, textinput.Blink
		}

		switch b.state {
		case addUserPermissionState:
			if msg.Type == tea.KeyEsc {
				b.state = idleState
				// b.resetAddUserInputs()
				b.viewport.SetContent(
					b.generateContent(
						b.viewport.Width,
						b.viewport.Height,
						b.permissionsTable.View()))
				return b, nil
			}
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeyUp || msg.Type == tea.KeyDown {
				s := msg.String()

				if s == "enter" && b.focusIndex == len(b.addUserPermissionInputs) {
					b.state = idleState
					// createStardogUserCmd := b.CreateStardogCmd(b.addUserInputs[0].Value(), b.addUserInputs[1].Value())
					addUserPermissionCmd := b.AddUserPermissionCmd(
						b.selectedUser,
						b.addUserPermissionInputs[0].Value(),
						b.addUserPermissionInputs[1].Value(),
						b.addUserPermissionInputs[2].Value())
					// b.addUserInputs[0].Reset()
					// b.addUserInputs[1].Reset()
					getUserPermissionsCmd := b.GetUserPermissionsCmd(b.selectedUser)

					cmds = append(cmds, tea.Sequentially(addUserPermissionCmd, getUserPermissionsCmd))
				}

				if s == "up" {
					b.focusIndex--
				} else {
					b.focusIndex++
				}

				if b.focusIndex > len(b.addUserPermissionInputs) {
					b.focusIndex = 0
				} else if b.focusIndex < 0 {
					b.focusIndex = len(b.addUserPermissionInputs)
				}

				cmds := make([]tea.Cmd, len(b.addUserPermissionInputs))
				for i := 0; i <= len(b.addUserPermissionInputs)-1; i++ {
					if i == b.focusIndex {
						cmds[i] = b.addUserPermissionInputs[i].Focus()
						b.addUserPermissionInputs[i].PromptStyle = focusedStyle
						b.addUserPermissionInputs[i].TextStyle = focusedStyle
						continue
					}
					b.addUserPermissionInputs[i].Blur()
					b.addUserPermissionInputs[i].PromptStyle = noStyle
					b.addUserPermissionInputs[i].TextStyle = noStyle
				}
			}
		}

		if msg.String() == "x" {
			cmds = append(cmds, b.DeleteUserPermissionCmd(b.selectedUser))
		}

		if b.active {
			switch b.state {
			case addUserPermissionState:
				cmd := b.updateAddUserPermissionInputs(msg)
				cmds = append(cmds, cmd)
			}

			b.permissionsTable, cmd = b.permissionsTable.Update(msg)
			cmds = append(cmds, cmd)

			b.viewport, cmd = b.viewport.Update(msg)
			cmds = append(cmds, cmd)

			b.viewport.SetContent(
				b.generateContent(
					b.viewport.Width,
					b.viewport.Height,
					b.permissionsTable.View()))

			return b, tea.Batch(cmds...)
		}
	}

	return b, tea.Batch(cmds...)
}

func (b *Bubble) updateAddUserPermissionInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(b.addUserPermissionInputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range b.addUserPermissionInputs {
		b.addUserPermissionInputs[i], cmds[i] = b.addUserPermissionInputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (b Bubble) View() string {

	if b.active {
		b.viewport.Style = b.viewport.Style.Copy().BorderForeground(lipgloss.Color("33"))
	} else {
		b.viewport.Style = b.viewport.Style.Copy().BorderForeground(lipgloss.Color("#FAFAFA"))
	}
	return b.viewport.View()
}
