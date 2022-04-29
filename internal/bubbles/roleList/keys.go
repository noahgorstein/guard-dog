package rolelist

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	addRoleKey = key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "add role"))

	deleteRoleKey = key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "delete role"))

	forceDeleteRoleKey = key.NewBinding(
		key.WithKeys(tea.KeyCtrlF.String()),
		key.WithHelp(tea.KeyCtrlF.String(), "force delete role"))
)
