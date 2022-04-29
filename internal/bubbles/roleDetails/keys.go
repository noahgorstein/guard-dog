package roledetails

import "github.com/charmbracelet/bubbles/key"

var (
	tableModeKey = key.NewBinding(
		key.WithKeys("T"),
		key.WithHelp("T", "activate table"),
	)

	grantRolePermissionKey = key.NewBinding(
		key.WithKeys("ctrl+g"),
		key.WithHelp("ctrl+a", "grant role permission"))

	revokeRolePermissionKey = key.NewBinding(
		key.WithKeys("ctrl+x"),
		key.WithHelp("ctrl+x", "revoke role permission"))
)
