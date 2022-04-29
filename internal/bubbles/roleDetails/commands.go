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

func getRowStyle(action string) lipgloss.Style {
	switch action {
	case "ALL":
		return red
	case "READ":
		return orange
	case "WRITE":
		return yellow
	case "CREATE":
		return green
	case "DELETE":
		return teal
	case "EXECUTE":
		return blue
	case "GRANT":
		return purple
	case "REVOKE":
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
		rolePermissions, err := b.stardogClient.Security.GetRolePermissions(context.Background(), b.selectedRole)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		usersWithRole, err := b.stardogClient.Security.ListUsersAssignedToRole(context.Background(), b.selectedRole)
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
		permission := stardog.NewPermission(
			row.Data[b.permissionTableColumnKeys.action].(string),
			row.Data[b.permissionTableColumnKeys.resourceType].(string),
			[]string{row.Data[b.permissionTableColumnKeys.resource].(string)})

		_, err := b.stardogClient.Security.RevokeRolePermission(context.Background(), b.selectedRole, *permission)
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

func (b *Bubble) CreateRolePermissionCmd(action string, resourceType string, resource []string) tea.Cmd {
	return func() tea.Msg {
		permission := stardog.NewPermission(action, resourceType, resource)
		_, err := b.stardogClient.Security.GrantRolePermission(context.Background(), b.selectedRole, *permission)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully added permission: [%s] [%s] [%s] to role: %s",
				permission.Action,
				permission.ResourceType,
				strings.Join(permission.Resource, "\\"),
				b.selectedRole),
		}
	}

}
