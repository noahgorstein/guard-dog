package userdetails

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

func (b *Bubble) SetCurrentUser(user string) {
	b.selectedUser = user
}

func (b Bubble) generateContent(width int) string {
	width = width - b.viewport.Style.GetHorizontalFrameSize()

	var sb strings.Builder

	sb.WriteString(b.Styles.header.Render(b.selectedUser))
	sb.WriteRune('\n')
	sb.WriteString(b.Styles.secondaryHeader.Render("Assigned Roles"))
	sb.WriteRune('\n')

	roles := wordwrap.NewWriter(b.viewport.Width)
	if len(b.assignedRoles) > 0 {
		for _, role := range b.assignedRoles {
			_, _ = roles.Write([]byte(lipgloss.NewStyle().
				Bold(true).
				Foreground(nord14).
				Render(role)))
			_, _ = roles.Write([]byte(" "))
		}
	} else {
		sb.WriteString("No assigned roles.")
	}
	roles.Close()
	sb.WriteString(roles.String())
	sb.WriteString("\n\n")

	switch b.State {
	case IdleState:
		sb.WriteString(b.Styles.secondaryHeader.Render("Controls"))
		sb.WriteRune('\n')
		sb.WriteString(b.getIdleStateHelp())
		sb.WriteRune('\n')
		sb.WriteString(b.Styles.secondaryHeader.Render("Explicit Permissions"))
		sb.WriteRune('\n')
		if b.permissionsTable.TotalRows() > 0 ||
			(b.permissionsTable.TotalRows() == 0 && b.permissionsTable.GetIsFilterActive()) {
			sb.WriteString(b.permissionsTable.WithTargetWidth(width).View())
		}
		if b.permissionsTable.TotalRows() == 0 && !b.permissionsTable.GetIsFilterActive() {
			sb.WriteString("User " + b.selectedUser + " has no permissions.")
			sb.WriteRune('\n')
			sb.WriteRune('\n')
		}
	case TableState:
		sb.WriteString(b.Styles.secondaryHeader.Render("Explicit Permissions"))
		sb.WriteRune('\n')
		if b.permissionsTable.TotalRows() > 0 ||
			(b.permissionsTable.TotalRows() == 0 && b.permissionsTable.GetIsFilterActive()) {
			sb.WriteString(b.getTableModeHelp())
			sb.WriteRune('\n')
			sb.WriteString(b.permissionsTable.WithTargetWidth(width).View())
		}
		if b.permissionsTable.TotalRows() == 0 && !b.permissionsTable.GetIsFilterActive() {
			sb.WriteString("User " + b.selectedUser + " has no permissions.")
			sb.WriteRune('\n')
			sb.WriteRune('\n')
			sb.WriteString(lipgloss.NewStyle().
				Bold(true).
				Foreground(nord8).
				Render("esc") + " to go back\n")
			sb.WriteString(lipgloss.NewStyle().
				Bold(true).
				Foreground(nord8).
				Render("ctrl+g") + " to grant a permission to user")
		}
	case AddingUserPermission:
		sb.WriteString(b.addPermissionsPrompt.View())
	case AssigningRoleState:
		sb.WriteString(b.Styles.secondaryHeader.Render("Add role to user: ") + "\n")
		sb.WriteString(b.getRemoveUserFromRoleHelp())

		if len(b.addRoleSelector.GetChoices()) > 0 {
			if b.focusIndex > 0 {
				b.submitButton = lipgloss.NewStyle().Foreground(nord12).Render("[ Submit ]")
			} else {
				b.submitButton = lipgloss.NewStyle().Foreground(b.Styles.grey).Render("[ Submit ]")
			}
			sb.WriteString(lipgloss.JoinVertical(
				lipgloss.Center,
				b.addRoleSelector.View(),
				b.submitButton,
			))
		} else {
			sb.WriteString("No roles available to assign.\n")
			sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(nord8).Render("esc") + " to go back")
		}
	case RemovingRoleState:
		sb.WriteString(b.Styles.secondaryHeader.Render("Remove role from user: ") + "\n")
		sb.WriteString(b.getRemoveUserFromRoleHelp())
		sb.WriteRune('\n')
		if len(b.removeRoleSelector.GetChoices()) > 0 {
			if b.focusIndex > 0 {
				b.submitButton = lipgloss.NewStyle().Foreground(nord12).Render("[ Submit ]")
			} else {
				b.submitButton = lipgloss.NewStyle().Foreground(b.Styles.grey).Render("[ Submit ]")
			}
			sb.WriteString(lipgloss.JoinVertical(
				lipgloss.Center,
				b.removeRoleSelector.View(),
				b.submitButton,
			))
		} else {
			sb.WriteString("User is not assigned any roles.\n")
			sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(nord8).Render("esc") + " to go back")
		}
	}
	emptyLinesPastLastContent := 10
	return lipgloss.NewStyle().
		Width(width).
		Height(lipgloss.Height(sb.String()) + emptyLinesPastLastContent).
		Render(sb.String())
}

func (b *Bubble) SetIsActive(active bool) {
	b.active = active
}

func (b *Bubble) SetSize(width, height int) {
	b.viewport.Width = width
	b.viewport.Height = height

	b.viewport.SetContent(
		b.generateContent(
			b.viewport.Width,
		))
}

func (b *Bubble) Reset() {
	b.State = IdleState
	b.viewport.GotoTop()
	b.resetPermissionsTable()
	b.addPermissionsPrompt.Reset()
	b.viewport.SetContent(
		b.generateContent(
			b.viewport.Width,
		))
}

func (b *Bubble) resetPermissionsTable() {
	b.permissionsTable = b.permissionsTable.
		WithBaseStyle(lipgloss.NewStyle().BorderForeground(b.Styles.grey))
}

func (b Bubble) getIdleStateHelp() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("↑/↓"),
		b.Styles.helpTextDescription.Render("......scroll viewport")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("T"),
		b.Styles.helpTextDescription.Render("........navigate/edit permissions")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("ctrl+g"),
		b.Styles.helpTextDescription.Render("...grant user permission")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("ctrl+r"),
		b.Styles.helpTextDescription.Render("...assign role to user")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("ctrl+u"),
		b.Styles.helpTextDescription.Render("...unassign role from user")))

	return sb.String()
}

func (b Bubble) getTableModeHelp() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("↑/↓/←/→"),
		b.Styles.helpTextDescription.Render("...navigate table")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("/"),
		b.Styles.helpTextDescription.Render(".........filter by action, resource type, or resource")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("esc"),
		b.Styles.helpTextDescription.Render(".......clear filter/back")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("ctrl+x"),
		b.Styles.helpTextDescription.Render("....revoke highlighted permission")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("ctrl+g"),
		b.Styles.helpTextDescription.Render("....grant user permission")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("ctrl+r"),
		b.Styles.helpTextDescription.Render("....assign role to user")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("ctrl+u"),
		b.Styles.helpTextDescription.Render("....unassign role from user")))

	return sb.String()
}

func (b Bubble) getRemoveUserFromRoleHelp() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("esc"),
		b.Styles.helpTextDescription.Render("......back")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("enter"),
		b.Styles.helpTextDescription.Render("....make selection and jump to next prompt")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("space"),
		b.Styles.helpTextDescription.Render("....make selection")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("ctrl+b"),
		b.Styles.helpTextDescription.Render("...previous prompt")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.helpTextKey.Render("ctrl+n"),
		b.Styles.helpTextDescription.Render("...next prompt")))
	return sb.String()
}
