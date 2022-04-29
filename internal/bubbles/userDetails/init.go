package userdetails

import tea "github.com/charmbracelet/bubbletea"

func (b Bubble) Init() tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	cmd = b.GetUserDetailsCmd(false)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}
