package userlist

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
	case GetUsersMsg:
		cmd = b.list.SetItems(msg)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch b.State {
		case AddUserState:
			if msg.Type == tea.KeyEsc {
				b.State = IdleState
				b.resetAddUserInputs()
				return b, nil
			}
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeyUp || msg.Type == tea.KeyDown {
				s := msg.String()

				if s == "enter" && b.focusIndex == len(b.addUserInputs) {
					b.State = IdleState
					createStardogUserCmd := b.CreateUserCmd(b.addUserInputs[0].Value(), b.addUserInputs[1].Value())
					b.addUserInputs[0].Reset()
					b.addUserInputs[1].Reset()
					cmds = append(cmds, createStardogUserCmd)
				}

				if s == "up" {
					b.focusIndex--
				} else {
					b.focusIndex++
				}

				if b.focusIndex > len(b.addUserInputs) {
					b.focusIndex = 0
				} else if b.focusIndex < 0 {
					b.focusIndex = len(b.addUserInputs)
				}

				cmds := make([]tea.Cmd, len(b.addUserInputs))
				for i := 0; i <= len(b.addUserInputs)-1; i++ {
					if i == b.focusIndex {
						cmds[i] = b.addUserInputs[i].Focus()
						b.addUserInputs[i].PromptStyle = b.Styles.focusedStyle
						b.addUserInputs[i].TextStyle = b.Styles.focusedStyle
						continue
					}
					b.addUserInputs[i].Blur()
					b.addUserInputs[i].PromptStyle = b.Styles.noStyle
					b.addUserInputs[i].TextStyle = b.Styles.noStyle
				}
			}

		case ChangeUserPasswordState:
			if msg.Type == tea.KeyEsc {
				b.State = IdleState
				b.resetChangeUserPasswordInputs()
				return b, nil
			}
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeyUp || msg.Type == tea.KeyDown {
				s := msg.String()
				if s == "enter" && b.focusIndex == len(b.changeUserPasswordInputs) {
					b.State = IdleState
					changeUserPasswordCmd := b.ChangeUserPasswordCmd(
						b.changeUserPasswordInputs[0].Value(),
						b.changeUserPasswordInputs[1].Value())
					b.changeUserPasswordInputs[0].Reset()
					b.changeUserPasswordInputs[1].Reset()
					cmds = append(cmds, changeUserPasswordCmd)
				}

				if s == "up" {
					b.focusIndex--
				} else {
					b.focusIndex++
				}

				if b.focusIndex > len(b.changeUserPasswordInputs) {
					b.focusIndex = 0
				} else if b.focusIndex < 0 {
					b.focusIndex = len(b.changeUserPasswordInputs)
				}

				cmds := make([]tea.Cmd, len(b.changeUserPasswordInputs))
				for i := 0; i <= len(b.changeUserPasswordInputs)-1; i++ {
					if i == b.focusIndex {
						cmds[i] = b.changeUserPasswordInputs[i].Focus()
						b.changeUserPasswordInputs[i].PromptStyle = b.Styles.focusedStyle
						b.changeUserPasswordInputs[i].TextStyle = b.Styles.focusedStyle
						continue
					}
					b.changeUserPasswordInputs[i].Blur()
					b.changeUserPasswordInputs[i].PromptStyle = b.Styles.noStyle
					b.changeUserPasswordInputs[i].TextStyle = b.Styles.noStyle
				}
			}
		case IdleState:
			switch {
			case key.Matches(msg, addUserKey):
				username := b.addUserInputs[0]
				username.Focus()
				username.PromptStyle = b.Styles.focusedStyle
				b.State = AddUserState
				return b, textinput.Blink
			case key.Matches(msg, changeUserPasswordKey):
				password := b.changeUserPasswordInputs[0]
				password.Focus()
				password.PromptStyle = b.Styles.focusedStyle
				b.State = ChangeUserPasswordState
				return b, textinput.Blink
			case key.Matches(msg, deleteUserKey):
				deleteStardogUserCmd := b.DeleteUserCmd()
				cmds = append(cmds, deleteStardogUserCmd)
			case key.Matches(msg, enableUserKey):
				enabledStardogUserCmd := b.EnableUserCmd(!b.GetSelectedUser().enabled)
				cmds = append(cmds, enabledStardogUserCmd)
			}
		}
	}

	if b.active {
		switch b.State {
		case IdleState:
			b.list, cmd = b.list.Update(msg)
			cmds = append(cmds, cmd)
		case AddUserState:
			cmd := b.updateAddUserInputs(msg)
			cmds = append(cmds, cmd)
		case ChangeUserPasswordState:
			cmd := b.updateChangeUserPasswordInputs(msg)
			cmds = append(cmds, cmd)
		}
	}

	return b, tea.Batch(cmds...)
}
