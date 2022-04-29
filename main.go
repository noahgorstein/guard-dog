package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/noahgorstein/stardog-go/stardog"
)

var connectionDetails = stardog.NewConnectionDetails("http://localhost:5820", "admin", "admin")

var listStyle = lipgloss.NewStyle().
	Margin(1, 2).
	BorderForeground(lipgloss.Color("202")).
	Border(lipgloss.RoundedBorder(), true)
var inactiveListStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.RoundedBorder(), true)

var viewportStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("202")).
	Margin(1, 2).
	Padding(1, 1)
var inactiveViewportStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	Margin(1, 2).
	Padding(1, 1)

type item struct {
	title, desc string
}

func (i item) Title() string { return i.title }

func (i item) Description() string { return i.desc }

func (i item) FilterValue() string { return i.title }

type sessionState uint

const (
	listView sessionState = iota
	detailsView
)

const (
	columnKeyName      = "name"
	columnKeyEnabled   = "enabled"
	columnKeySuperuser = "superuser"
	columnKeyRoles     = "roles"
)

const (
	columnKeyAction       = "actioon"
	columnKeyResourceType = "resource type"
	columnKeyResource     = "resource"
	columnKeyExplicit     = "explicit"
)

type model struct {
	list                 list.Model
	userDetailsTable     table.Model
	userPermissionsTable table.Model
	viewport             viewport.Model
	user                 stardog.User
	userDetails          stardog.GetUserDetailsResponse
	state                sessionState
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) updateUserDetailsTable(selectedUser string) {
	m.user = stardog.User{Name: selectedUser}
	m.userDetails = stardog.GetUserDetails(*connectionDetails, m.user)
	userDetailsRow := []table.Row{table.NewRow(table.RowData{
		columnKeyName:      m.user.Name,
		columnKeyEnabled:   m.userDetails.Enabled,
		columnKeySuperuser: m.userDetails.Superuser,
		columnKeyRoles:     strings.Join(m.userDetails.Roles, ", "),
	})}
	m.userDetailsTable = m.userDetailsTable.WithRows(userDetailsRow)
}

func (m *model) updateUserPermissionsTable(selectedUser string) {
	m.user = stardog.User{Name: selectedUser}
	m.userDetails = stardog.GetUserDetails(*connectionDetails, m.user)
	rows := []table.Row{}
	for _, permission := range m.userDetails.Permissions {
		rows = append(rows, table.NewRow(table.RowData{
			columnKeyAction:       permission.Action,
			columnKeyResourceType: permission.ResourceType,
			columnKeyResource:     permission.Resource,
			columnKeyExplicit:     permission.Explicit,
		}))
	}
	m.userPermissionsTable = m.userPermissionsTable.WithRows(rows)
}

func (m *model) toggleActiveView() {
	if m.state == listView {
		m.state = detailsView
	} else {
		m.state = listView
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// user may want to use that char
		if m.list.FilterState() != list.Filtering {
			if msg.String() == "b" {
				m.toggleActiveView()
			}
		}
		if m.state == listView {
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}
		if m.state == detailsView {
			m.userDetailsTable, cmd = m.userDetailsTable.Update(msg)
			cmds = append(cmds, cmd)

			m.userPermissionsTable, cmd = m.userPermissionsTable.Update(msg)
			cmds = append(cmds, cmd)

			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		// m.viewport.HighPerformanceRendering = true
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.viewport.Height = msg.Height - h
	}

	if m.list.SelectedItem() != nil {
		selectedUser := m.list.SelectedItem().(item).title
		m.updateUserDetailsTable(selectedUser)
		m.updateUserPermissionsTable(selectedUser)
		m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, m.userDetailsTable.View(), m.userPermissionsTable.View()))
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.state == listView {
		return lipgloss.JoinHorizontal(lipgloss.Top, listStyle.Render(m.list.View()), inactiveViewportStyle.Render(m.viewport.View()))
	} else {
		return lipgloss.JoinHorizontal(lipgloss.Top, inactiveListStyle.Render(m.list.View()), viewportStyle.Render(m.viewport.View()))
	}
}

func main() {
	users := stardog.GetUsers(*connectionDetails)

	items := []list.Item{}
	for _, user := range users {
		items = append(items, item{title: user.Name})
	}

	userDetailTable := table.New([]table.Column{
		table.NewColumn(columnKeyName, "Username", 10).WithFiltered(true),
		table.NewColumn(columnKeyEnabled, "Enabled", 10).WithFiltered(true),
		table.NewColumn(columnKeySuperuser, "Superuser", 10).WithFiltered(true),
		table.NewColumn(columnKeyRoles, "Roles", 10).WithFiltered(true),
	}).Focused(true)

	userPermissionsTable := table.New([]table.Column{
		table.NewColumn(columnKeyAction, "Action", 10).WithFiltered(true),
		table.NewColumn(columnKeyResourceType, "Resource Type", 20).WithFiltered(true),
		table.NewColumn(columnKeyResource, "Resource", 30).WithFiltered(true),
		table.NewColumn(columnKeyExplicit, "Explicit Permission", 20).WithFiltered(true),
	}).Focused(true)

	theViewport := viewport.New(100, 40)
	itemDelegate := list.NewDefaultDelegate()
	itemDelegate.Styles.SelectedTitle.Foreground(lipgloss.Color("202"))
	itemDelegate.Styles.SelectedTitle.BorderLeftForeground(lipgloss.Color("202"))
	itemDelegate.Styles.SelectedDesc = lipgloss.NewStyle()
	itemDelegate.SetSpacing(0)

	m := model{
		list:                 list.New(items, itemDelegate, 0, 0),
		viewport:             theViewport,
		userDetailsTable:     userDetailTable,
		userPermissionsTable: userPermissionsTable,
		state:                listView,
	}
	m.list.Title = "Users"

	// m.list.Styles.Title.Foreground(lipgloss.Color("888B7E"))
	// m.list.Styles.TitleBar.Background(lipgloss.Color("#888B7E"))
	m.list.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("202")).
		Foreground(lipgloss.Color("202")).
		Padding(0, 1)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
