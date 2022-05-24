package tuiv2

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	userlist "github.com/noahgorstein/stardog-go/internal/bubbles/userList"
	"github.com/noahgorstein/stardog-go/stardog"
)

var (
	f, _ = tea.LogToFile("debug.log", "debug")
)

func (b *Bubble) toggleActiveView() {
	if b.activeView == listView {
		b.activeView = detailsView
		b.userDetails.SetIsActive(true)
		b.userList.SetIsActive(false)
	} else {
		b.activeView = listView
		b.userDetails.SetIsActive(false)
		b.userList.SetIsActive(true)
	}
}

func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.userList.SetSize(int(float64(msg.Width)*0.3), msg.Height-lipgloss.Height(b.statusBar.View()))
		b.userDetails.SetSize(int(float64(msg.Width)*0.7), msg.Height-lipgloss.Height(b.statusBar.View()))

		user := stardog.User{Name: b.userList.GetCurrentUser()}
		b.userDetails.SetCurrentUser(user)

		cmd = b.userDetails.GetUserPermissionsCmd(user)
		cmds = append(cmds, cmd)

		b.statusBar, cmd = b.statusBar.Update(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		if msg.Type == tea.KeyTab {
			b.toggleActiveView()
		}
		if b.activeView == detailsView {
			switch msg.String() {
			case "ctrl+c":
				return b, tea.Quit
			}

			b.userDetails, cmd = b.userDetails.Update(msg)
			cmds = append(cmds, cmd)

			return b, tea.Batch(cmds...)
		}
		if b.activeView == listView {

			switch msg.String() {
			case "ctrl+r":
				cmd = userlist.GetStardogUsersCmd(b.stardogConnection)
				cmds = append(cmds, cmd)
			}

			b.userList, cmd = b.userList.Update(msg)
			cmds = append(cmds, cmd)

			user := stardog.User{Name: b.userList.GetCurrentUser()}
			b.userDetails.SetCurrentUser(user)

			cmd = b.userDetails.GetUserPermissionsCmd(user)
			cmds = append(cmds, cmd)

			return b, tea.Batch(cmds...)

		}
	}

	b.userList, cmd = b.userList.Update(msg)
	cmds = append(cmds, cmd)

	b.userDetails, cmd = b.userDetails.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)

}
