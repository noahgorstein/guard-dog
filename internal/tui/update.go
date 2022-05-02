package tui

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/noahgorstein/stardog-go/stardog"
)

func (b *Bubble) toggleActiveView() {
	if b.activeView == listView {
		b.activeView = detailsView
	} else {
		b.activeView = listView
	}
}

func (b *Bubble) selectedUserPermissions() string {
	selectedUser := b.list.SelectedItem().(item).title
	b.updateUserDetailsTable(selectedUser)
	b.updateUserPermissionsTable(selectedUser)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		usernameStyle.Render(selectedUser),
		b.userDetailsTable.View(),
		permissionsStyle.Render("Permissions"),
		b.userPermissionsTable.View())
}

func (b *Bubble) updateUserDetailsTable(selectedUser string) {
	b.user = stardog.User{Name: selectedUser}
	b.userDetails = stardog.GetUserDetails(*b.connection, b.user)
	userDetailsRow := []table.Row{table.NewRow(table.RowData{
		columnKeyEnabled:   b.userDetails.Enabled,
		columnKeySuperuser: b.userDetails.Superuser,
		columnKeyRoles:     strings.Join(b.userDetails.Roles, ", "),
	})}
	b.userDetailsTable = b.userDetailsTable.WithRows(userDetailsRow).WithTargetWidth(b.viewport.Width - 5)

}

func (b *Bubble) updateUserPermissionsTable(selectedUser string) {
	b.user = stardog.User{Name: selectedUser}
	b.userDetails = stardog.GetUserDetails(*b.connection, b.user)
	rows := []table.Row{}
	for _, permission := range b.userDetails.Permissions {
		rows = append(rows, table.NewRow(table.RowData{
			columnKeyAction:       permission.Action,
			columnKeyResourceType: permission.ResourceType,
			columnKeyResource:     strings.Join(permission.Resource, ", "),
			columnKeyExplicit:     strconv.FormatBool(permission.Explicit),
		}))
	}
	b.userPermissionsTable = b.userPermissionsTable.WithRows(rows).Filtered(true).WithPageSize(15).WithTargetWidth(b.viewport.Width - 5)
}

func (b *Bubble) updateStatusBar() {
	statusBarStyle.Width(lipgloss.Width(b.list.View()) + lipgloss.Width(b.viewport.View()) + listStyle.GetHorizontalFrameSize() + viewportStyle.GetHorizontalFrameSize() - statusBarStyle.GetHorizontalFrameSize())
	b.statusBar = statusBarStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			b.connection.Username,
			"@",
			b.connection.Endpoint,
		))
}

func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return b, tea.Quit
		}
		if b.list.FilterState() != list.Filtering {
			if msg.String() == "b" {
				b.toggleActiveView()
			}
		}
		if b.activeView == detailsView {

			if msg.String() == "s" {
				b.columnSortKey = columnKeyAction
				b.userPermissionsTable = b.userPermissionsTable.SortByAsc(b.columnSortKey)
			}

			b.userDetailsTable, cmd = b.userDetailsTable.Update(msg)
			cmds = append(cmds, cmd)

			b.userPermissionsTable, cmd = b.userPermissionsTable.Update(msg)
			cmds = append(cmds, cmd)

			b.viewport, cmd = b.viewport.Update(msg)
			cmds = append(cmds, cmd)

			if b.list.SelectedItem() != nil {
				b.viewport.SetContent(b.selectedUserPermissions())
			}

			return b, tea.Batch(cmds...)

		}

	case tea.WindowSizeMsg:
		b.width, b.height = msg.Width, msg.Height
		listV, _ := listStyle.GetFrameSize()
		statusBarV, _ := statusBarStyle.GetFrameSize()

		b.list.SetSize(int(float64(msg.Width)*0.3), msg.Height-listV-statusBarV-statusBarStyle.GetHeight())
		b.viewport = viewport.New(msg.Width-lipgloss.Width(b.list.View())-listStyle.GetWidth(), b.list.Height())

		if b.list.SelectedItem() != nil {
			b.viewport.SetContent(b.selectedUserPermissions())
		}

		b.updateStatusBar()

	}

	if b.list.FilterState() == list.FilterApplied ||
		b.list.FilterState() == list.Filtering ||
		b.list.FilterState() == list.Unfiltered {
		b.viewport = viewport.New(b.width-lipgloss.Width(b.list.View())-listStyle.GetWidth(), b.list.Height())
	}

	b.list, cmd = b.list.Update(msg)
	cmds = append(cmds, cmd)

	if b.list.SelectedItem() != nil {
		b.viewport.SetContent(b.selectedUserPermissions())
	}

	return b, tea.Batch(cmds...)
}
