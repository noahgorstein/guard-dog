package userdetails

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/noahgorstein/stardog-go/stardog"
)

type GetUserPermissionsMsg []table.Row

func (b *Bubble) GetUserPermissionsCmd(user stardog.User) tea.Cmd {
	return func() tea.Msg {
		rows := []table.Row{}
		userDetails := stardog.GetUserDetails(b.connection, user)

		for _, permission := range userDetails.Permissions {
			rows = append(rows, table.NewRow(table.RowData{
				columnKeyAction:       permission.Action,
				columnKeyResourceType: permission.ResourceType,
				columnKeyResource:     strings.Join(permission.Resource, ", "),
				columnKeyExplicit:     strconv.FormatBool(permission.Explicit),
			}))
		}

		return GetUserPermissionsMsg(rows)
	}
}

type DeleteUserPermissionMsg struct{ success bool }

func (b *Bubble) DeleteUserPermissionCmd(user stardog.User) tea.Cmd {
	return func() tea.Msg {

		row := b.permissionsTable.HighlightedRow()
		permission := stardog.NewPermission(
			row.Data[columnKeyAction].(string),
			row.Data[columnKeyResourceType].(string),
			[]string{row.Data[columnKeyResource].(string)})
		stardog.DeleteUserPermission(b.connection, user, *permission)
		return DeleteUserPermissionMsg{success: true}
	}

}

func (b *Bubble) AddUserPermissionCmd(user stardog.User, action string, resourceType string, resource string) tea.Cmd {
	return func() tea.Msg {
		permission := stardog.NewPermission(action, resourceType, []string{resource})
		stardog.AddUserPermission(b.connection, user, *permission)
		return nil
	}

}
