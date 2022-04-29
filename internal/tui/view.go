package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/noahgorstein/guard-dog/internal/mode"
)

func (b Bubble) View() string {

	if b.activeView == mode.UserListMode || b.activeView == mode.UserDetailsMode {
		if b.helpEnabled {
			return lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(lipgloss.Top, b.userList.View(), b.userDetails.View(), b.help.View()),
				b.statusBar.View())
		}
		return lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top, b.userList.View(), b.userDetails.View()),
			b.statusBar.View())
	} else {
		if b.helpEnabled {
			return lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(lipgloss.Top, b.roleList.View(), b.roleDetails.View(), b.help.View()),
				b.statusBar.View())
		}
		return lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top, b.roleList.View(), b.roleDetails.View()),
			b.statusBar.View())
	}
}
