package userlist

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
		inputView = "\n\n"
	case AddUserState:

		for i := range b.addUserInputs {
			builder.WriteString(b.addUserInputs[i].View())
			if i < len(b.addUserInputs)-1 {
				builder.WriteRune('\n')

			}
		}
		b.enableSubmitButton(false)
		if b.focusIndex == len(b.addUserInputs) {
			b.enableSubmitButton(true)
		}
		builder.WriteString("\n" + b.submitButton)

		inputView = builder.String()
	case ChangeUserPasswordState:
		for i := range b.changeUserPasswordInputs {
			builder.WriteString(b.changeUserPasswordInputs[i].View())
			if i < len(b.changeUserPasswordInputs)-1 {
				builder.WriteRune('\n')
			}
		}
		b.enableSubmitButton(false)
		if b.focusIndex == len(b.changeUserPasswordInputs) {
			b.enableSubmitButton(true)
		}
		builder.WriteString("\n" + b.submitButton)
		inputView = builder.String()
	}

	return b.Styles.listStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			b.list.View(),
			b.divider,
			b.Styles.addUserInputStyle.Render(inputView)))
}
