package userdetails

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	addpermissionprompt "github.com/noahgorstein/guard-dog/internal/bubbles/addPermissionPrompt"
)

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.addPermissionsPrompt.SetSize(b.viewport.Width)
	case GetUserDetailsMsg:
		b.assignedRoles = msg.Roles
		b.permissionsTable = b.permissionsTable.WithRows(msg.Permissions)
		b.viewport.SetContent(
			b.generateContent(
				b.viewport.Width,
				b.viewport.Height,
			))
	case GetAvailableRolesMsg:
		b.addRoleSelector.SetChoices(msg)
		b.viewport.SetContent(
			b.generateContent(
				b.viewport.Width,
				b.viewport.Height,
			))
	case GetAssignedRolesMsg:
		b.removeRoleSelector.SetChoices(msg)
		b.viewport.SetContent(
			b.generateContent(
				b.viewport.Width,
				b.viewport.Height,
			))
	case tea.KeyMsg:

		if b.active {
			if key.Matches(msg, helpKey) {
				b.viewport, cmd = b.viewport.Update(msg)
				cmds = append(cmds, cmd)

				return b, tea.Batch(cmds...)
			}
			switch b.State {
			case IdleState:

				if key.Matches(msg, tableModeKey) {
					b.State = TableState
					b.permissionsTable = b.permissionsTable.
						WithBaseStyle(lipgloss.NewStyle().BorderForeground(nord12))
					b.viewport.SetContent(
						b.generateContent(
							b.viewport.Width,
							b.viewport.Height,
						))
				}

				if key.Matches(msg, grantUserPermissionKey) {
					b.State = AddingUserPermission
					b.viewport.SetContent(
						b.generateContent(
							b.viewport.Width,
							b.viewport.Height,
						))
				}
				if key.Matches(msg, assignRoleToUser) {
					b.State = AssigningRoleState
					b.addRoleSelector.SetIsActive(true)
					GetAvailableRolesCmd := b.GetAvailableRoles()
					cmds = append(cmds, GetAvailableRolesCmd)
				}
				if key.Matches(msg, removeRoleFromUser) {
					b.State = RemovingRoleState
					b.removeRoleSelector.SetIsActive(true)
					GetAssignedRolesCmd := b.GetRolesAssignedToUser()
					cmds = append(cmds, GetAssignedRolesCmd)
				}
			case TableState:

				if !b.permissionsTable.GetIsFilterActive() && msg.Type == tea.KeyEsc {
					b.State = IdleState
					b.resetPermissionsTable()
				}

				if key.Matches(msg, grantUserPermissionKey) {
					b.State = AddingUserPermission
					b.resetPermissionsTable()
					b.viewport.SetContent(
						b.generateContent(
							b.viewport.Width,
							b.viewport.Height,
						))
				}
				if key.Matches(msg, revokeUserPermissionKey) {
					if b.permissionsTable.TotalRows() > 0 {
						cmds = append(cmds, b.DeleteUserPermissionCmd())
					}
				}
				if key.Matches(msg, assignRoleToUser) {
					b.State = AssigningRoleState
					b.resetPermissionsTable()
					b.addRoleSelector.SetIsActive(true)
					GetAvailableRolesCmd := b.GetAvailableRoles()
					cmds = append(cmds, GetAvailableRolesCmd)
				}
				if key.Matches(msg, removeRoleFromUser) {
					b.State = RemovingRoleState
					b.resetPermissionsTable()
					b.removeRoleSelector.SetIsActive(true)
					GetAssignedRolesCmd := b.GetRolesAssignedToUser()
					cmds = append(cmds, GetAssignedRolesCmd)
				}

				b.permissionsTable, cmd = b.permissionsTable.Update(msg)
				cmds = append(cmds, cmd)

				b.viewport.SetContent(
					b.generateContent(
						b.viewport.Width,
						b.viewport.Height,
					))

				return b, tea.Batch(cmds...)
			case AddingUserPermission:

				if msg.Type == tea.KeyEsc {
					b.addPermissionsPrompt.Reset()
					b.State = IdleState
				}

				if b.addPermissionsPrompt.State == addpermissionprompt.SubmitState &&
					msg.Type == tea.KeyEnter {
					addUserPermissionCmd := b.GrantUserPermissionCmd(
						b.addPermissionsPrompt.Permission.Action,
						b.addPermissionsPrompt.Permission.ResourceType,
						b.addPermissionsPrompt.Permission.Resource,
					)
					cmds = append(cmds, addUserPermissionCmd)
					b.addPermissionsPrompt.Reset()
					b.viewport.GotoTop()
					b.State = IdleState
					return b, tea.Batch(cmds...)
				}
				b.addPermissionsPrompt, cmd = b.addPermissionsPrompt.Update(msg)
				cmds = append(cmds, cmd)
			case AssigningRoleState:
				switch msg.Type {
				case tea.KeyEsc:
					b.focusIndex = 0
					b.State = IdleState
					b.addRoleSelector.Reset()
				case tea.KeyCtrlB:
					b.focusIndex = 0
					b.addRoleSelector.SetIsActive(true)
				case tea.KeyCtrlN:
					b.focusIndex = 1
					b.addRoleSelector.SetIsActive(false)
				case tea.KeyEnter, tea.KeySpace:
					if b.focusIndex > 0 {
						addRoleToUserCmd := b.AssignUserRole(b.addRoleSelector.GetSelected())
						cmds = append(cmds, addRoleToUserCmd)
						b.focusIndex = 0
						b.State = IdleState
						b.viewport.GotoTop()
						b.addRoleSelector.Reset()
					} else {
						b.addRoleSelector, cmd = b.addRoleSelector.Update(msg)
						cmds = append(cmds, cmd)

						b.focusIndex++
						b.addRoleSelector.SetIsActive(false)
					}
				default:
					b.addRoleSelector, cmd = b.addRoleSelector.Update(msg)
					cmds = append(cmds, cmd)
				}
			case RemovingRoleState:
				switch msg.Type {
				case tea.KeyEsc:
					b.focusIndex = 0
					b.State = IdleState
					b.removeRoleSelector.Reset()
				case tea.KeyCtrlB:
					b.focusIndex = 0
					b.removeRoleSelector.SetIsActive(true)
				case tea.KeyCtrlN:
					b.focusIndex = 1
					b.removeRoleSelector.SetIsActive(false)
				case tea.KeyEnter, tea.KeySpace:
					if b.focusIndex > 0 {
						removeRoleFromUserCmd := b.RemoveRoleFromUser(b.removeRoleSelector.GetSelected())
						cmds = append(cmds, removeRoleFromUserCmd)
						b.focusIndex = 0
						b.State = IdleState
						b.viewport.GotoTop()
						b.removeRoleSelector.Reset()
					} else {
						b.removeRoleSelector, cmd = b.removeRoleSelector.Update(msg)
						cmds = append(cmds, cmd)

						b.focusIndex++
						b.removeRoleSelector.SetIsActive(false)
					}
				default:
					b.removeRoleSelector, cmd = b.removeRoleSelector.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
			b.viewport, cmd = b.viewport.Update(msg)
			cmds = append(cmds, cmd)

			b.viewport.SetContent(
				b.generateContent(
					b.viewport.Width,
					b.viewport.Height,
				))

			return b, tea.Batch(cmds...)
		}

	}

	return b, tea.Batch(cmds...)
}
