package roledetails

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/noahgorstein/go-stardog/stardog"
)

type SuccessMsg struct {
	Message string
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

var (
	red      = lipgloss.NewStyle().Bold(true).Foreground(nord11)
	orange   = lipgloss.NewStyle().Bold(true).Foreground(nord12)
	yellow   = lipgloss.NewStyle().Bold(true).Foreground(nord13)
	green    = lipgloss.NewStyle().Bold(true).Foreground(nord14)
	teal     = lipgloss.NewStyle().Bold(true).Foreground(nord7)
	blue     = lipgloss.NewStyle().Bold(true).Foreground(nord8)
	purple   = lipgloss.NewStyle().Bold(true).Foreground(nord15)
	darkBlue = lipgloss.NewStyle().Bold(true).Foreground(nord9)
)

func getRowStyle(action stardog.PermissionAction) lipgloss.Style {
	switch action {
	case stardog.PermissionActionAll:
		return red
	case stardog.PermissionActionRead:
		return orange
	case stardog.PermissionActionWrite:
		return yellow
	case stardog.PermissionActionCreate:
		return green
	case stardog.PermissionActionDelete:
		return teal
	case stardog.PermissionActionExecute:
		return blue
	case stardog.PermissionActionGrant:
		return purple
	case stardog.PermissionActionRevoke:
		return darkBlue
	default:
		return lipgloss.NewStyle()
	}
}

type GetRoleDetailsMsg struct {
	permissions []table.Row
	roles       []string
}

func (b *Bubble) GetRoleDetailsCmd() tea.Cmd {
	return func() tea.Msg {
		rows := []table.Row{}
		rolePermissions, _, err := b.stardogClient.Role.Permissions(context.Background(), b.selectedRole)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		usersWithRole, _, err := b.stardogClient.User.ListNamesAssignedRole(context.Background(), b.selectedRole)
		if err != nil {
			return errMsg{
				err: err,
			}
		}

		for _, permission := range rolePermissions {

			rows = append(rows, table.NewRow(table.RowData{
				b.permissionTableColumnKeys.action:       permission.Action,
				b.permissionTableColumnKeys.resourceType: permission.ResourceType,
				b.permissionTableColumnKeys.resource:     strings.Join(permission.Resource, "\\"),
			}).WithStyle(getRowStyle(permission.Action)))
		}

		return GetRoleDetailsMsg{
			permissions: rows,
			roles:       usersWithRole,
		}
	}
}

func (b *Bubble) DeleteRolePermissionCmd() tea.Cmd {
	return func() tea.Msg {

		row := b.permissionsTable.HighlightedRow()
		permission := &stardog.Permission{
			Action:       row.Data[b.permissionTableColumnKeys.action].(stardog.PermissionAction),
			ResourceType: row.Data[b.permissionTableColumnKeys.resourceType].(stardog.PermissionResourceType),
			Resource:     []string{row.Data[b.permissionTableColumnKeys.resource].(string)},
		}
		// permission := stardog.NewPermission(
		// 	row.Data[b.permissionTableColumnKeys.action].(string),
		// 	row.Data[b.permissionTableColumnKeys.resourceType].(string),
		// 	[]string{row.Data[b.permissionTableColumnKeys.resource].(string)})

		_, err := b.stardogClient.Role.RevokePermission(context.Background(), b.selectedRole, *permission)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully removed permission: [%s] [%s] [%s] from role: %s",
				permission.Action,
				permission.ResourceType,
				strings.Join(permission.Resource, "\\"),
				b.selectedRole),
		}
	}
}

func (b *Bubble) CreateRolePermissionCmd(p stardog.Permission) tea.Cmd {
	return func() tea.Msg {
		_, err := b.stardogClient.Role.GrantPermission(context.Background(), b.selectedRole, p)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully added permission: [%s] [%s] [%s] to role: %s",
				p.Action,
				p.ResourceType,
				strings.Join(p.Resource, "\\"),
				b.selectedRole),
		}
	}

}
