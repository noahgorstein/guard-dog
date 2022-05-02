package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// View returns a string representation of the UI.
func (b Bubble) View() string {
	leftBox := b.list.View()
	rightBox := b.viewport.View()

	switch b.activeView {
	case listView:
		leftBox = listStyle.Render(b.list.View())
		rightBox = inactiveViewportStyle.Render(b.viewport.View())
	case detailsView:
		leftBox = inactiveListStyle.Render(b.list.View())
		rightBox = viewportStyle.Render(b.viewport.View())
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftBox,
			rightBox,
		),
		b.statusBar)
}
