package statusbar

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	UserListView    key.Binding
	UserDetailsView key.Binding
	RoleListView    key.Binding
	RoleDetailsView key.Binding

	Help  key.Binding
	Quit  key.Binding
	Roles key.Binding
	Users key.Binding
}

var Keys = KeyMap{
	UserListView: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "user list"),
	),
	UserDetailsView: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "user details"),
	),
	RoleListView: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "role list"),
	),
	RoleDetailsView: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "role details"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Roles: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "roles")),
	Users: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "users"),
	),
}
