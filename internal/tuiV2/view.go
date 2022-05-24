package tuiv2

import "github.com/charmbracelet/lipgloss"

// View returns a string representation of the UI.
func (b Bubble) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, b.userList.View(), b.userDetails.View()),
		b.statusBar.View())
}
