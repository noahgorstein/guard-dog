package userdetails

import "github.com/charmbracelet/bubbles/key"

var (
	helpKey = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"))

	tableModeKey = key.NewBinding(
		key.WithKeys("T"),
		key.WithHelp("T", "activate table"),
	)
	grantUserPermissionKey = key.NewBinding(
		key.WithKeys("ctrl+g"),
		key.WithHelp("ctrl+g", "grant user permission"))

	revokeUserPermissionKey = key.NewBinding(
		key.WithKeys("ctrl+x"),
		key.WithHelp("ctrl+x", "revoke user permission"))

	assignRoleToUser = key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "assign user to role"))

	removeRoleFromUser = key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "unassign user from role"))
)
