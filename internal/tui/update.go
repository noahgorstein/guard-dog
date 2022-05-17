package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/noahgorstein/stardog-go/stardog"
)

var (
	f, _ = tea.LogToFile("debug.log", "debug")
)

func (b *Bubble) toggleActiveView() {
	if b.activeView == listView {
		b.activeView = detailsView
	} else {
		b.activeView = listView
	}
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

	b.userPermissionsTable = table.New([]table.Column{
		table.NewColumn(columnKeyAction, "Action", 10).WithFiltered(true),
		table.NewColumn(columnKeyResourceType, "Resource Type", 20).WithFiltered(true),
		table.NewFlexColumn(columnKeyResource, "Resource", 1).WithFiltered(true),
		table.NewColumn(columnKeyExplicit, "Explicit Permission", 20).WithFiltered(true),
	}).Focused(true).
		SelectableRows(true).
		WithBaseStyle(permissionsTableStyle).
		HighlightStyle(permissionsTableHighltedRowStyle).
		WithRows(rows).
		WithTargetWidth(b.viewport.Width - 5).
		WithPageSize(10).
		Filtered(true)

}

func (b *Bubble) updateUserPermissionsTableFooter() {

	highlightedRow := b.userPermissionsTable.HighlightedRow().Data
	f.WriteString(fmt.Sprintf("HIGHLIGHTED ROW: %s", highlightedRow) + "\n")

	selectedText := strings.Builder{}
	selectedIDs := []string{}

	for _, row := range b.userPermissionsTable.SelectedRows() {
		selectedIDs = append(selectedIDs, row.Data[columnKeyAction].(string))
	}

	selectedText.WriteString(fmt.Sprintf("SelectedIDs: %s\n", strings.Join(selectedIDs, ", ")))

	footerText := fmt.Sprintf(
		"Pg. %d/%d - Currently looking at ID: %s - selected: %s",
		b.userPermissionsTable.CurrentPage(),
		b.userPermissionsTable.MaxPages(),
		highlightedRow,
		selectedText.String(),
	)

	b.userPermissionsTable = b.userPermissionsTable.WithStaticFooter(footerText)

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
		switch b.activeView {
		case detailsView:

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

			selectedUser := b.list.SelectedItem().(item).title
			b.updateUserPermissionsTableFooter()

			if b.list.SelectedItem() != nil {
				b.viewport.SetContent(lipgloss.JoinVertical(
					lipgloss.Left,
					usernameStyle.Render(selectedUser),
					b.userDetailsTable.View(),
					permissionsStyle.Render("Permissions"),
					b.userPermissionsTable.View()))
			}

			return b, tea.Batch(cmds...)
		case listView:
			f.WriteString("In list view! \n")

			b.list, cmd = b.list.Update(msg)
			cmds = append(cmds, cmd)

			if b.list.FilterState() == list.FilterApplied ||
				b.list.FilterState() == list.Filtering ||
				b.list.FilterState() == list.Unfiltered {
				b.viewport = viewport.New(b.width-lipgloss.Width(b.list.View())-listStyle.GetWidth(), b.list.Height())
			}

			if b.list.SelectedItem() != nil {
				selectedUser := b.list.SelectedItem().(item).title
				b.updateUserDetailsTable(selectedUser)
				b.updateUserPermissionsTable(selectedUser)
				// b.updateUserPermissionsTableFooter()

				b.viewport.SetContent(lipgloss.JoinVertical(
					lipgloss.Left,
					usernameStyle.Render(selectedUser),
					b.userDetailsTable.View(),
					permissionsStyle.Render("Permissions"),
					b.userPermissionsTable.View()))

			}
			return b, tea.Batch(cmds...)
		}
	case tea.WindowSizeMsg:
		b.width, b.height = msg.Width, msg.Height
		listV, _ := listStyle.GetFrameSize()
		statusBarV, _ := statusBarStyle.GetFrameSize()

		b.list.SetSize(int(float64(msg.Width)*0.3), msg.Height-listV-statusBarV-statusBarStyle.GetHeight())
		b.viewport = viewport.New(msg.Width-lipgloss.Width(b.list.View())-listStyle.GetWidth(), b.list.Height())

		b.list, cmd = b.list.Update(msg)
		cmds = append(cmds, cmd)

		if b.list.SelectedItem() != nil {
			selectedUser := b.list.SelectedItem().(item).title
			b.updateUserDetailsTable(selectedUser)
			b.updateUserPermissionsTable(selectedUser)
			// b.updateUserPermissionsTableFooter()

			b.viewport.SetContent(lipgloss.JoinVertical(
				lipgloss.Left,
				usernameStyle.Render(selectedUser),
				b.userDetailsTable.View(),
				permissionsStyle.Render("Permissions"),
				b.userPermissionsTable.View()))
		}

		b.updateStatusBar()
		return b, tea.Batch(cmds...)

	default:
		b.list, cmd = b.list.Update(msg)
		cmds = append(cmds, cmd)

		b.userDetailsTable, cmd = b.userDetailsTable.Update(msg)
		cmds = append(cmds, cmd)

		b.userPermissionsTable, cmd = b.userPermissionsTable.Update(msg)
		cmds = append(cmds, cmd)

		return b, tea.Batch(cmds...)
	}

	b.list, cmd = b.list.Update(msg)
	cmds = append(cmds, cmd)

	b.userDetailsTable, cmd = b.userDetailsTable.Update(msg)
	cmds = append(cmds, cmd)

	b.userPermissionsTable, cmd = b.userPermissionsTable.Update(msg)
	cmds = append(cmds, cmd)
	return b, tea.Batch(cmds...)
}
