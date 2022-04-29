package rolelist

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case GetRolesMsg:
		cmd = b.list.SetItems(msg)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch b.State {
		case IdleState:
			switch {
			case key.Matches(msg, addRoleKey):
				username := b.addRoleInputs[0]
				username.Focus()
				username.PromptStyle = b.Styles.focusedStyle
				b.State = AddRoleState
				return b, textinput.Blink
			case key.Matches(msg, deleteRoleKey):
				deleteRoleCmd := b.DeleteRoleCmd()
				cmds = append(cmds, deleteRoleCmd)
				return b, tea.Batch(cmds...)
			case key.Matches(msg, forceDeleteRoleKey):
				forceDeleteRoleCmd := b.ForceDeleteRoleCmd()
				cmds = append(cmds, forceDeleteRoleCmd)
				return b, tea.Batch(cmds...)
			}
		case AddRoleState:
			if msg.Type == tea.KeyEsc {
				b.State = IdleState
				b.resetAddRoleInputs()
				return b, nil
			}
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeyUp || msg.Type == tea.KeyDown {
				s := msg.String()

				if s == "enter" && b.focusIndex == len(b.addRoleInputs) {
					b.State = IdleState
					createRoleCmd := b.CreateRoleCmd(b.addRoleInputs[0].Value())
					b.addRoleInputs[0].Reset()
					cmds = append(cmds, createRoleCmd)
				}

				if s == "up" {
					b.focusIndex--
				} else {
					b.focusIndex++
				}

				if b.focusIndex > len(b.addRoleInputs) {
					b.focusIndex = 0
				} else if b.focusIndex < 0 {
					b.focusIndex = len(b.addRoleInputs)
				}

				cmds := make([]tea.Cmd, len(b.addRoleInputs))
				for i := 0; i <= len(b.addRoleInputs)-1; i++ {
					if i == b.focusIndex {
						cmds[i] = b.addRoleInputs[i].Focus()
						b.addRoleInputs[i].PromptStyle = b.Styles.focusedStyle
						b.addRoleInputs[i].TextStyle = b.Styles.focusedStyle
						continue
					}
					b.addRoleInputs[i].Blur()
					b.addRoleInputs[i].PromptStyle = b.Styles.noStyle
					b.addRoleInputs[i].TextStyle = b.Styles.noStyle
				}
			}

		}
	}

	if b.active {
		switch b.State {
		case IdleState:
			b.list, cmd = b.list.Update(msg)
			cmds = append(cmds, cmd)
		case AddRoleState:
			cmd := b.updateAddUserInputs(msg)
			cmds = append(cmds, cmd)
		}

	}

	return b, tea.Batch(cmds...)
}
