package roledetails

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

func (b *Bubble) SetCurrentRole(role string) {
	b.selectedRole = role
}

func (b *Bubble) SetIsActive(active bool) {
	b.active = active
}

func (b Bubble) generateContent(width, height int) string {

	if b.selectedRole == "" {

		var noRolesText strings.Builder
		noRolesText.WriteString(b.Styles.textStyle.Render("No roles in the system."))
		noRolesText.WriteRune('\n')
		noRolesText.WriteRune('\n')
		noRolesText.WriteString(fmt.Sprintf("%s %s",
			b.Styles.emphasizedTextStyle.Render("ctrl+n"),
			b.Styles.textStyle.Render("to create a role")))

		return lipgloss.NewStyle().
			Width(b.viewport.Width).
			Height(b.viewport.Height).
			Render(noRolesText.String())
	}

	var sb strings.Builder

	sb.WriteString(b.Styles.header.Render(b.selectedRole))
	sb.WriteRune('\n')

	sb.WriteString(b.Styles.secondaryHeader.Render(fmt.Sprintf("Users assigned to %s", b.selectedRole)))
	sb.WriteRune('\n')
	users := wordwrap.NewWriter(b.viewport.Width)
	if len(b.usersAssignedToRole) > 0 {
		for _, user := range b.usersAssignedToRole {
			users.Write([]byte(lipgloss.NewStyle().
				Bold(true).
				Foreground(nord14).
				Render(user)))
			users.Write([]byte(" "))
		}
	} else {
		sb.WriteString(b.Styles.textStyle.Render("Role is not assigned to any user."))
	}
	sb.WriteString(users.String())
	sb.WriteRune('\n')
	sb.WriteRune('\n')

	switch b.State {
	case IdleState:
		sb.WriteString(b.Styles.secondaryHeader.Render("Controls"))
		sb.WriteRune('\n')
		sb.WriteString(b.getIdleStateHelp())
		sb.WriteRune('\n')

		sb.WriteString(b.Styles.secondaryHeader.Render("Role Permissions"))
		sb.WriteRune('\n')
		if b.permissionsTable.TotalRows() > 0 ||
			(b.permissionsTable.TotalRows() == 0 && b.permissionsTable.GetIsFilterActive()) {

			sb.WriteString(b.permissionsTable.WithTargetWidth(width).View())
			sb.WriteRune('\n')
		}
		if b.permissionsTable.TotalRows() == 0 && !b.permissionsTable.GetIsFilterActive() {
			sb.WriteString(b.Styles.textStyle.Render("The role has no permissions."))
			sb.WriteRune('\n')
			sb.WriteRune('\n')
			sb.WriteString(fmt.Sprintf("%s %s",
				b.Styles.emphasizedTextStyle.Render("ctrl+g"),
				b.Styles.textStyle.Render("to grant role permission.")))
		}
	case TableState:
		sb.WriteString(b.Styles.secondaryHeader.Render("Controls"))
		sb.WriteRune('\n')
		sb.WriteString(b.getTableStateHelp())
		sb.WriteRune('\n')

		sb.WriteString(b.Styles.secondaryHeader.Render("Role Permissions"))
		sb.WriteRune('\n')
		if b.permissionsTable.TotalRows() > 0 ||
			(b.permissionsTable.TotalRows() == 0 && b.permissionsTable.GetIsFilterActive()) {

			sb.WriteString(b.permissionsTable.WithTargetWidth(width).View())
			sb.WriteRune('\n')
		}
		if b.permissionsTable.TotalRows() == 0 && !b.permissionsTable.GetIsFilterActive() {
			sb.WriteString(b.Styles.textStyle.Render("The role has no permissions."))
			sb.WriteRune('\n')
			sb.WriteRune('\n')
			sb.WriteString(lipgloss.NewStyle().
				Bold(true).
				Foreground(nord8).
				Render("esc") + " to go back\n")
			sb.WriteString(fmt.Sprintf("%s %s",
				b.Styles.emphasizedTextStyle.Render("ctrl+g"),
				b.Styles.textStyle.Render("to grant role permission")))
		}
	case AddingRolePermissionState:
		sb.WriteString(b.addPermissionsPrompt.View())
	}

	return lipgloss.NewStyle().
		Width(b.viewport.Width).
		Height(b.viewport.Height).
		Render(sb.String())
}

func (b *Bubble) SetSize(width, height int) {
	b.viewport.Width = width - b.viewport.Style.GetHorizontalFrameSize()
	b.viewport.Height = height - b.viewport.Style.GetVerticalFrameSize()

	b.viewport.SetContent(b.generateContent(b.viewport.Width, b.viewport.Height))
}

func (b *Bubble) Reset() {
	b.State = IdleState
	b.viewport.GotoTop()
	b.addPermissionsPrompt.Reset()
	b.resetPermissionsTable()
	b.viewport.SetContent(
		b.generateContent(
			b.viewport.Width,
			b.viewport.Height,
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
		b.Styles.helpTextDescription.Render("...grant role permission")))

	return sb.String()
}

func (b Bubble) getTableStateHelp() string {
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
		b.Styles.helpTextDescription.Render("....grant role permission")))

	return sb.String()
}
