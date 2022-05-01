package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/noahgorstein/stardog-go/stardog"
)

var connectionDetails = stardog.NewConnectionDetails("http://localhost:5820", "admin", "admin")

var (
	listStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("202")).
			Border(lipgloss.RoundedBorder(), true)

	listTitleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("202")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Bold(true).
			Italic(true)

	inactiveListStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), true)

	viewportStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("202"))

	inactiveViewportStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder())

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

	statusBarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).Height(1)

	usernameStatusBarStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Margin(1, 1)

	endpointStatusBarStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("202")).Background(lipgloss.Color("#FAFAFA"))
)

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
	columnKeyAction       = "action"
	columnKeyResourceType = "resource type"
	columnKeyResource     = "resource"
	columnKeyExplicit     = "explicit"
)

type model struct {
	width                int
	height               int
	list                 list.Model
	userDetailsTable     table.Model
	userPermissionsTable table.Model
	viewport             viewport.Model
	user                 stardog.User
	userDetails          stardog.GetUserDetailsResponse
	state                sessionState
	columnSortKey        string
	statusBar            string
}

func newModel() tea.Model {
	users := stardog.GetUsers(*connectionDetails)

	items := []list.Item{}
	for _, user := range users {
		items = append(items, item{title: user.Name})
	}

	userDetailTable := table.New([]table.Column{
		table.NewColumn(columnKeyEnabled, "Enabled", 15),
		table.NewColumn(columnKeySuperuser, "Superuser", 15),
		table.NewFlexColumn(columnKeyRoles, "Roles", 30),
	})

	userPermissionsTable := table.New([]table.Column{
		table.NewColumn(columnKeyAction, "Action", 10).WithFiltered(true),
		table.NewColumn(columnKeyResourceType, "Resource Type", 20).WithFiltered(true),
		table.NewFlexColumn(columnKeyResource, "Resource", 30).WithFiltered(true),
		table.NewColumn(columnKeyExplicit, "Explicit Permission", 20).WithFiltered(true),
	}).Focused(true).HighlightStyle(lipgloss.NewStyle().Background(lipgloss.Color("#004B60")))

	itemDelegate := list.NewDefaultDelegate()
	itemDelegate.Styles.SelectedTitle.Foreground(lipgloss.Color("202"))
	itemDelegate.Styles.SelectedTitle.BorderLeftForeground(lipgloss.Color("202"))
	itemDelegate.Styles.SelectedDesc = lipgloss.NewStyle()
	// itemDelegate.Styles.NormalTitle.Width(70)
	itemDelegate.SetSpacing(0)

	m := model{
		list:                 list.New(items, itemDelegate, 0, 0),
		userDetailsTable:     userDetailTable,
		userPermissionsTable: userPermissionsTable,
		state:                listView,
	}

	m.list.Styles.Title = listTitleStyle
	m.list.Title = "Stardog Users"

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) updateUserDetailsTable(selectedUser string) {
	m.user = stardog.User{Name: selectedUser}
	m.userDetails = stardog.GetUserDetails(*connectionDetails, m.user)
	userDetailsRow := []table.Row{table.NewRow(table.RowData{
		columnKeyEnabled:   m.userDetails.Enabled,
		columnKeySuperuser: m.userDetails.Superuser,
		columnKeyRoles:     strings.Join(m.userDetails.Roles, ", "),
	})}
	m.userDetailsTable = m.userDetailsTable.WithRows(userDetailsRow).WithTargetWidth(m.viewport.Width - viewportStyle.GetHorizontalFrameSize())
}

func (m *model) updateUserPermissionsTable(selectedUser string) {
	m.user = stardog.User{Name: selectedUser}
	m.userDetails = stardog.GetUserDetails(*connectionDetails, m.user)
	rows := []table.Row{}
	for _, permission := range m.userDetails.Permissions {
		rows = append(rows, table.NewRow(table.RowData{
			columnKeyAction:       permission.Action,
			columnKeyResourceType: permission.ResourceType,
			columnKeyResource:     strings.Join(permission.Resource, ", "),
			columnKeyExplicit:     strconv.FormatBool(permission.Explicit),
		}))
	}
	m.userPermissionsTable = m.userPermissionsTable.WithRows(rows).Filtered(true).WithPageSize(15).WithTargetWidth(m.viewport.Width - viewportStyle.GetHorizontalFrameSize())
}

func (m *model) toggleActiveView() {
	if m.state == listView {
		m.state = detailsView
	} else {
		m.state = listView
	}
}

func (m model) permissionsView() string {
	selectedUser := m.list.SelectedItem().(item).title
	m.updateUserDetailsTable(selectedUser)
	m.updateUserPermissionsTable(selectedUser)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		usernameStyle.Render(selectedUser),
		m.userDetailsTable.View(),
		permissionsStyle.Render("Permissions"),
		m.userPermissionsTable.View())
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
		if m.list.FilterState() != list.Filtering {
			if msg.String() == "b" {
				m.toggleActiveView()
			}
		}
		if m.state == detailsView {

			if msg.String() == "s" {
				m.columnSortKey = columnKeyAction
				m.userPermissionsTable = m.userPermissionsTable.SortByAsc(m.columnSortKey)
			}

			m.userDetailsTable, cmd = m.userDetailsTable.Update(msg)
			cmds = append(cmds, cmd)

			m.userPermissionsTable, cmd = m.userPermissionsTable.Update(msg)
			cmds = append(cmds, cmd)

			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)

			if m.list.SelectedItem() != nil {
				selectedUser := m.list.SelectedItem().(item).title
				m.updateUserDetailsTable(selectedUser)
				m.updateUserPermissionsTable(selectedUser)

				m.viewport.SetContent(
					lipgloss.JoinVertical(lipgloss.Left,
						usernameStyle.Render(selectedUser),
						m.userDetailsTable.View(),
						permissionsStyle.Render("Permissions"),
						m.userPermissionsTable.View()))
			}

			return m, tea.Batch(cmds...)

		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		height := m.height - 5

		// listViewWidth := int(0.3 * float64(m.width))
		// listWidth := listViewWidth - listStyle.GetHorizontalFrameSize()
		// m.list.SetSize(listWidth, height)
		m.list.SetSize(m.width/2, height)

		m.viewport = viewport.New(m.width-m.list.Width(), height)
		m.viewport.SetContent(m.permissionsView())
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	if m.list.SelectedItem() != nil {
		m.viewport.SetContent(m.permissionsView())
	}

	return m, tea.Batch(cmds...)
}

func (m model) statusView() string {
	statusBarStyle.Width(m.width - statusBarStyle.GetHorizontalFrameSize())
	status := endpointStatusBarStyle.Render(connectionDetails.Username + "@" + connectionDetails.Endpoint)
	return statusBarStyle.Render(status)
}

func (m model) View() string {

	m.viewport.SetContent(m.permissionsView())

	if m.state == listView {
		return lipgloss.JoinVertical(lipgloss.Left, lipgloss.JoinHorizontal(
			lipgloss.Top,
			listStyle.Render(m.list.View()),
			inactiveViewportStyle.Render(m.viewport.View()),
		), m.statusView())
	} else {
		return lipgloss.JoinVertical(lipgloss.Left, lipgloss.JoinHorizontal(
			lipgloss.Top,
			inactiveListStyle.Render(m.list.View()),
			viewportStyle.Render(m.viewport.View()),
		), m.statusView())
	}
}

func main() {

	p := tea.NewProgram(newModel(), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
