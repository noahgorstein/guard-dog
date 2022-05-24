package userlist

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
	stardog "github.com/noahgorstein/stardog-go/stardog"
)

type GetStardogUsersMsg []list.Item

func GetStardogUsersCmd(connectionDetails stardog.ConnectionDetails) tea.Cmd {
	return func() tea.Msg {
		var items []list.Item
		users := stardog.GetUsers(connectionDetails)
		for _, user := range users {
			items = append(items, item{title: user.Name, desc: ""})
		}
		return GetStardogUsersMsg(items)
	}
}

func (b *Bubble) CreateStardogCmd(username, password string) tea.Cmd {
	return func() tea.Msg {
		credentials := stardog.Credentials{
			Username: username,
			Password: strings.Split(password, ""),
		}
		stardog.AddUser(b.connection, credentials)
		return nil
	}
}

func (b *Bubble) DeleteStardogUserCmd(username string) tea.Cmd {
	return func() tea.Msg {
		stardog.DeleteUser(b.connection, stardog.User{Name: username})
		return nil
	}
}
