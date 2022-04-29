package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/noahgorstein/guard-dog/internal/mode"
)

func (b *Bubble) toggleActiveView() {

	switch b.activeView {
	case mode.UserListMode:
		b.activeView = mode.UserDetailsMode
		b.userDetails.SetIsActive(true)
		b.userList.SetIsActive(false)
		b.help.Entries = getUserDetailsViewHelpEntries()
	case mode.UserDetailsMode:
		b.activeView = mode.UserListMode
		b.userDetails.SetIsActive(false)
		b.userList.SetIsActive(true)
		b.help.Entries = getUserListViewHelpEntries()
	case mode.RoleListMode:
		b.activeView = mode.RoleDetailsMode
		b.roleDetails.SetIsActive(true)
		b.roleList.SetIsActive(false)
		b.help.Entries = getRoleDetailsViewHelpEntries()
	case mode.RoleDetailsMode:
		b.activeView = mode.RoleListMode
		b.roleDetails.SetIsActive(false)
		b.roleList.SetIsActive(true)
		b.help.Entries = getRoleListViewHelpEntries()
	}
	b.statusBar.UpdateMode(b.activeView)

}

func (b *Bubble) enableHelp(enabled bool) {

	if !enabled {
		b.roleList.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))
		b.userList.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))
		b.userDetails.SetSize(int(float64(b.width)*0.7), b.height-lipgloss.Height(b.statusBar.View()))
		b.roleDetails.SetSize(int(float64(b.width)*0.7), b.height-lipgloss.Height(b.statusBar.View()))
	} else {
		b.roleList.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))
		b.userList.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))
		b.userDetails.SetSize(int(float64(b.width)*0.4), b.height-lipgloss.Height(b.statusBar.View()))
		b.roleDetails.SetSize(int(float64(b.width)*0.4), b.height-lipgloss.Height(b.statusBar.View()))
		b.help.SetSize(int(float64(b.width)*0.3), b.height-lipgloss.Height(b.statusBar.View()))

	}
}
