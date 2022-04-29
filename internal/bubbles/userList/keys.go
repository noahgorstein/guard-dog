package userlist

import "github.com/charmbracelet/bubbles/key"

var (
	addUserKey = key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "add user"))

	deleteUserKey = key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "delete user"))

	changeUserPasswordKey = key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("ctrl+p", "change password"))

	enableUserKey = key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "enable/disable user"))
)
