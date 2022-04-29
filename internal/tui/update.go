package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	roledetails "github.com/noahgorstein/guard-dog/internal/bubbles/roleDetails"
	rolelist "github.com/noahgorstein/guard-dog/internal/bubbles/roleList"
	userdetails "github.com/noahgorstein/guard-dog/internal/bubbles/userDetails"
	userlist "github.com/noahgorstein/guard-dog/internal/bubbles/userList"
	"github.com/noahgorstein/guard-dog/internal/mode"
)

func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case userlist.SuccessMsg:
		cmd = b.statusBar.NewStatusMessage(msg.Message, true)
		cmds = append(cmds, cmd)

		cmd = b.userList.GetUsersCmd()
		cmds = append(cmds, cmd)

		return b, tea.Batch(cmds...)

	case userdetails.SuccessMsg:
		cmd = b.statusBar.NewStatusMessage(msg.Message, true)
		cmds = append(cmds, cmd)

		cmds = append(cmds, b.userDetails.GetUserDetailsCmd(true))
	case rolelist.SuccessMsg:
		cmd = b.statusBar.NewStatusMessage(msg.Message, true)
		cmds = append(cmds, cmd)

		cmd = b.roleList.GetRolesCmd()
		cmds = append(cmds, cmd)

	case roledetails.SuccessMsg:
		cmd = b.statusBar.NewStatusMessage(msg.Message, true)
		cmds = append(cmds, cmd)

		cmds = append(cmds, b.roleDetails.GetRoleDetailsCmd())
	case error:
		cmd = b.statusBar.NewStatusMessage(msg.Error(), false)
		cmds = append(cmds, cmd)
	case rolelist.GetRolesMsg:
		b.roleList, cmd = b.roleList.Update(msg)
		cmds = append(cmds, cmd)

		selectedRole := b.roleList.GetSelectedRole()
		b.roleDetails.SetCurrentRole(selectedRole)

		getRoleDetailsCmd := b.roleDetails.GetRoleDetailsCmd()
		cmds = append(cmds, getRoleDetailsCmd)

		return b, tea.Batch(cmds...)
	case userlist.GetUsersMsg:
		b.userList, cmd = b.userList.Update(msg)
		cmds = append(cmds, cmd)

		selectedUser := b.userList.GetSelectedUser().Username()
		b.userDetails.SetCurrentUser(selectedUser)

		getUserDetailsCmd := b.userDetails.GetUserDetailsCmd(false)
		cmds = append(cmds, getUserDetailsCmd)

		return b, tea.Batch(cmds...)
	case userdetails.GetUserDetailsMsg:
		if msg.UpdateOccured {
			cmd = b.userList.GetUsersCmd()
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:

		b.width = msg.Width
		b.height = msg.Height

		b.statusBar.SetWidth(msg.Width)

		if b.helpEnabled {
			b.roleList.SetSize(int(float64(msg.Width)*0.3), msg.Height-lipgloss.Height(b.statusBar.View()))
			b.userList.SetSize(int(float64(msg.Width)*0.3), msg.Height-lipgloss.Height(b.statusBar.View()))
			b.userDetails.SetSize(int(float64(msg.Width)*0.4), msg.Height-lipgloss.Height(b.statusBar.View()))
			b.roleDetails.SetSize(int(float64(msg.Width)*0.4), msg.Height-lipgloss.Height(b.statusBar.View()))
			b.help.SetSize(int(float64(msg.Width)*0.3), msg.Height-lipgloss.Height(b.statusBar.View()))
		} else {
			b.roleList.SetSize(int(float64(msg.Width)*0.3), msg.Height-lipgloss.Height(b.statusBar.View()))
			b.userList.SetSize(int(float64(msg.Width)*0.3), msg.Height-lipgloss.Height(b.statusBar.View()))
			b.userDetails.SetSize(int(float64(msg.Width)*0.7), msg.Height-lipgloss.Height(b.statusBar.View()))
			b.roleDetails.SetSize(int(float64(msg.Width)*0.7), msg.Height-lipgloss.Height(b.statusBar.View()))
		}

		selectedUser := b.userList.GetSelectedUser().Username()
		b.userDetails.SetCurrentUser(selectedUser)

		cmd = b.userDetails.GetUserDetailsCmd(false)
		cmds = append(cmds, cmd)

		b.statusBar, cmd = b.statusBar.Update(msg)
		cmds = append(cmds, cmd)

		b.userList, cmd = b.userList.Update(msg)
		cmds = append(cmds, cmd)

		b.userDetails, cmd = b.userDetails.Update(msg)
		cmds = append(cmds, cmd)

		b.roleList, cmd = b.roleList.Update(msg)
		cmds = append(cmds, cmd)

		b.roleDetails, cmd = b.roleDetails.Update(msg)
		cmds = append(cmds, cmd)

		return b, tea.Batch(cmds...)

	case tea.KeyMsg:

		if b.activeView == mode.UserDetailsMode || b.activeView == mode.UserListMode ||
			b.activeView == mode.RoleDetailsMode || b.activeView == mode.RoleListMode {
			if msg.String() == "ctrl+c" {
				return b, tea.Quit
			}
			if msg.String() == "?" {
				if b.helpEnabled {
					b.enableHelp(false)
					b.helpEnabled = false
				} else {
					b.enableHelp(true)
					b.helpEnabled = true
				}
			}
		}

		if b.activeView == mode.UserDetailsMode || b.activeView == mode.UserListMode {
			if msg.String() == "tab" {
				b.userDetails.Reset()
				b.toggleActiveView()

				// kind of hacky but need to update help bubbles viewport content
				if b.helpEnabled {
					b.userList.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))
					b.userDetails.SetSize(int(float64(b.width)*0.4), b.height-lipgloss.Height(b.statusBar.View()))
					b.help.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))
				}
			}

			if msg.String() == "r" && b.userList.State == userlist.IdleState &&
				!b.userList.IsFilterActive() && b.userDetails.State == userdetails.IdleState {
				b.activeView = mode.RoleListMode
				b.roleList.SetIsActive(true)
				b.roleDetails.SetIsActive(false)
				if b.helpEnabled {
					b.help.Entries = getRoleListViewHelpEntries()
					b.roleList.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))
					b.roleDetails.SetSize(int(float64(b.width)*0.4), b.height-lipgloss.Height(b.statusBar.View()))
					b.help.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))
				}
				getRolesCmd := b.roleList.GetRolesCmd()
				cmds = append(cmds, getRolesCmd)

				b.statusBar.UpdateMode(b.activeView)

			}
		}
		if b.activeView == mode.RoleDetailsMode || b.activeView == mode.RoleListMode {

			if msg.String() == "tab" {
				b.roleDetails.Reset()
				b.toggleActiveView()

				if b.helpEnabled {
					b.roleList.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))
					b.roleDetails.SetSize(int(float64(b.width)*0.4), b.height-lipgloss.Height(b.statusBar.View()))
					b.help.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))
				}
			}

			if msg.String() == "u" && b.roleList.State == rolelist.IdleState &&
				!b.roleList.IsFilterActive() && b.roleDetails.State == roledetails.IdleState {
				b.activeView = mode.UserListMode
				b.userList.SetIsActive(true)
				b.userDetails.SetIsActive(false)
				if b.helpEnabled {
					b.help.Entries = getUserListViewHelpEntries()
					b.userList.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))
					b.userDetails.SetSize(int(float64(b.width)*0.4), b.height-lipgloss.Height(b.statusBar.View()))
					b.help.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))
				}
				getUsersCmd := b.userList.GetUsersCmd()
				cmds = append(cmds, getUsersCmd)

				b.statusBar.UpdateMode(b.activeView)
			}
		}

		if b.activeView == mode.UserDetailsMode {

			b.userDetails, cmd = b.userDetails.Update(msg)
			cmds = append(cmds, cmd)

			return b, tea.Batch(cmds...)
		}

		if b.activeView == mode.UserListMode {

			b.userList, cmd = b.userList.Update(msg)
			cmds = append(cmds, cmd)

			selectedUser := b.userList.GetSelectedUser().Username()
			b.userDetails.SetCurrentUser(selectedUser)

			// if we're deleting the user, we'll get an error trying
			// to retrieve details for it once it's deleted
			if msg.Type != tea.KeyCtrlD {
				getUserDetailsCmd := b.userDetails.GetUserDetailsCmd(false)
				cmds = append(cmds, getUserDetailsCmd)
			}

			return b, tea.Batch(cmds...)
		}

		if b.activeView == mode.RoleDetailsMode {
			b.roleDetails, cmd = b.roleDetails.Update(msg)
			cmds = append(cmds, cmd)

			return b, tea.Batch(cmds...)
		}

		if b.activeView == mode.RoleListMode {
			b.roleList, cmd = b.roleList.Update(msg)
			cmds = append(cmds, cmd)

			selectedRole := b.roleList.GetSelectedRole()
			b.roleDetails.SetCurrentRole(selectedRole)

			// if we're deleting the role, we'll get an error trying
			// to retrieve details for it once it's deleted
			if msg.Type != tea.KeyCtrlF && msg.Type != tea.KeyCtrlD {
				cmd = b.roleDetails.GetRoleDetailsCmd()
				cmds = append(cmds, cmd)
			}

			return b, tea.Batch(cmds...)

		}
	}

	switch b.activeView {
	case mode.UserListMode:
		b.userList, cmd = b.userList.Update(msg)
		cmds = append(cmds, cmd)

		b.userDetails, cmd = b.userDetails.Update(msg)
		cmds = append(cmds, cmd)
	case mode.UserDetailsMode:
		b.userList, cmd = b.userList.Update(msg)
		cmds = append(cmds, cmd)

		b.userDetails, cmd = b.userDetails.Update(msg)
		cmds = append(cmds, cmd)
	case mode.RoleListMode:
		b.roleList, cmd = b.roleList.Update(msg)
		cmds = append(cmds, cmd)

		b.roleDetails, cmd = b.roleDetails.Update(msg)
		cmds = append(cmds, cmd)
	case mode.RoleDetailsMode:
		b.roleDetails, cmd = b.roleDetails.Update(msg)
		cmds = append(cmds, cmd)
	}

	b.statusBar, cmd = b.statusBar.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)

}
