package rolelist

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (b Bubble) View() string {

	if b.active {
		b.Styles.listStyle = b.Styles.ActiveListStyle
	} else {
		b.Styles.listStyle = b.Styles.InactiveListStyle
	}

	var inputView string
	var builder strings.Builder

	switch b.State {
	case IdleState:
		inputView = "\n"
	case AddRoleState:

		for i := range b.addRoleInputs {
			builder.WriteString(b.addRoleInputs[i].View())
		}
		b.enableSubmitButton(false)
		if b.focusIndex == len(b.addRoleInputs) {
			b.enableSubmitButton(true)
		}
		builder.WriteString("\n" + b.submitButton)

		inputView = builder.String()
	}
	return b.Styles.listStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			b.list.View(),
			b.Styles.addRoleInputStyle.Render(inputView)))
}
