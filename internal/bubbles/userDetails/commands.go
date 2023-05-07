package userdetails

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

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

type GetRolesMsg []string

func (b *Bubble) GetRolesCmd() tea.Cmd {
	return func() tea.Msg {
		roleList, _, err := b.stardogClient.Role.ListNames(context.Background())
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return GetRolesMsg(roleList)
	}
}

type GetAvailableRolesMsg []string

func (b *Bubble) GetAvailableRoles() tea.Cmd {
	return func() tea.Msg {

		availableRoles := map[string]struct{}{}

		roleList, _, err := b.stardogClient.Role.ListNames(context.Background())
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		for _, role := range roleList {
			availableRoles[role] = struct{}{}
		}

		assignedRoles, _, err := b.stardogClient.User.Roles(context.Background(), b.selectedUser)
		if err != nil {
			return errMsg{
				err: err,
			}
		}

		for _, assignedRole := range assignedRoles {
			delete(availableRoles, assignedRole)
		}

		var availableRolesList []string
		for availableRole := range availableRoles {
			availableRolesList = append(availableRolesList, availableRole)
		}
		return GetAvailableRolesMsg(availableRolesList)
	}
}

type GetAssignedRolesMsg []string

func (b *Bubble) GetRolesAssignedToUser() tea.Cmd {
	return func() tea.Msg {
		roleList, _, err := b.stardogClient.User.Roles(context.Background(), b.selectedUser)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return GetAssignedRolesMsg(roleList)
	}
}

func (b *Bubble) RemoveRoleFromUser(role string) tea.Cmd {
	return func() tea.Msg {

		_, err := b.stardogClient.User.UnassignRole(context.Background(), b.selectedUser, role)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully removed role: %s from: %s", role, b.selectedUser),
		}
	}
}

func (b *Bubble) AssignUserRole(role string) tea.Cmd {
	return func() tea.Msg {
		_, err := b.stardogClient.User.AssignRole(context.Background(), b.selectedUser, role)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully assigned role: %s to: %s", role, b.selectedUser),
		}
	}
}

type GetUserDetailsMsg struct {
	Permissions   []table.Row
	Roles         []string
	UpdateOccured bool
}

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

func (b *Bubble) GetUserDetailsCmd(updateOccured bool) tea.Cmd {
	return func() tea.Msg {
		rows := []table.Row{}
		userDetails, _, err := b.stardogClient.User.Get(context.Background(), b.selectedUser)
		if err != nil {
			return errMsg{
				err: err,
			}
		}

		for _, permission := range userDetails.EffectivePermissions {

			if permission.Explicit {

				rows = append(rows, table.NewRow(table.RowData{
					b.permissionTableColumnKeys.action:       permission.Action,
					b.permissionTableColumnKeys.resourceType: permission.ResourceType,
					b.permissionTableColumnKeys.resource:     strings.Join(permission.Resource, "\\"),
				}).WithStyle(getRowStyle(permission.Action)))
			}
		}

		return GetUserDetailsMsg{
			Permissions:   rows,
			Roles:         userDetails.Roles,
			UpdateOccured: updateOccured,
		}
	}
}

func (b *Bubble) DeleteUserPermissionCmd() tea.Cmd {
	return func() tea.Msg {
		row := b.permissionsTable.HighlightedRow()
		permission := stardog.Permission{
			Action:       row.Data[b.permissionTableColumnKeys.action].(stardog.PermissionAction),
			ResourceType: row.Data[b.permissionTableColumnKeys.resourceType].(stardog.PermissionResourceType),
			Resource:     []string{row.Data[b.permissionTableColumnKeys.resource].(string)}}

		_, err := b.stardogClient.User.RevokePermission(context.Background(), b.selectedUser, permission)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully removed permission: [%s] [%s] [%s]",
				permission.Action,
				permission.ResourceType,
				strings.Join(permission.Resource, "\\")),
		}
	}

}

func (b *Bubble) GrantUserPermissionCmd(permission stardog.Permission) tea.Cmd {
	return func() tea.Msg {
		_, err := b.stardogClient.User.GrantPermission(context.Background(), b.selectedUser, permission)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully added permission: [%s] [%s] [%s]",
				permission.Action,
				permission.ResourceType,
				strings.Join(permission.Resource, "\\")),
		}
	}

}
