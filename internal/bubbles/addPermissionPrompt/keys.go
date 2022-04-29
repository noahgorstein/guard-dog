package addpermissionprompt

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	nextKey = key.NewBinding(
		key.WithKeys(tea.KeyCtrlN.String()),
		key.WithHelp("->", "next input"))
	previousKey = key.NewBinding(
		key.WithKeys(tea.KeyCtrlB.String()),
		key.WithHelp("down", "next input"))
)
