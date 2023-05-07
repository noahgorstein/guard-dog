package roledetails

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	permissionPrompt "github.com/noahgorstein/guard-dog/internal/bubbles/addPermissionPrompt"
)

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.addPermissionsPrompt.SetSize(b.viewport.Width)
	case GetRoleDetailsMsg:
		b.usersAssignedToRole = msg.roles
		b.permissionsTable = b.permissionsTable.WithRows(msg.permissions)
		b.viewport.SetContent(
			b.generateContent(
				b.viewport.Width,
				b.viewport.Height))
	case tea.KeyMsg:
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

			if key.Matches(msg, grantRolePermissionKey) {
				b.State = AddingRolePermissionState
				b.resetPermissionsTable()
				b.viewport.SetContent(
					b.generateContent(
						b.viewport.Width,
						b.viewport.Height))
			}

		case TableState:

			if key.Matches(msg, revokeRolePermissionKey) {
				if b.permissionsTable.TotalRows() > 0 {
					deleteRolePermissionsCmd := b.DeleteRolePermissionCmd()
					cmds = append(cmds, deleteRolePermissionsCmd)
				}
			}
			if key.Matches(msg, grantRolePermissionKey) {
				b.State = AddingRolePermissionState
				b.resetPermissionsTable()
				b.viewport.SetContent(
					b.generateContent(
						b.viewport.Width,
						b.viewport.Height))
			}

			if !b.permissionsTable.GetIsFilterActive() && msg.Type == tea.KeyEsc {
				b.viewport.GotoTop()
				b.State = IdleState
				b.resetPermissionsTable()
			}

			b.permissionsTable, cmd = b.permissionsTable.Update(msg)
			cmds = append(cmds, cmd)

		case AddingRolePermissionState:

			if msg.String() == "enter" && b.addPermissionsPrompt.State == permissionPrompt.SubmitState {
				addRolePermissionCmd := b.CreateRolePermissionCmd(b.addPermissionsPrompt.Permission)
				cmds = append(cmds, addRolePermissionCmd)

				b.addPermissionsPrompt.Reset()
				b.viewport.GotoTop()
				b.State = IdleState
			}

			if msg.String() == "esc" {
				b.viewport.GotoTop()
				b.addPermissionsPrompt.Reset()
				b.State = IdleState
			}

			b.addPermissionsPrompt, cmd = b.addPermissionsPrompt.Update(msg)
			cmds = append(cmds, cmd)
		}

	}

	b.viewport, cmd = b.viewport.Update(msg)
	cmds = append(cmds, cmd)

	b.viewport.SetContent(b.generateContent(b.viewport.Width, b.viewport.Height))

	return b, tea.Batch(cmds...)

}
