package rolelist

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/noahgorstein/go-stardog/stardog"
)

type SuccessMsg struct {
	Message string
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type GetRolesMsg []list.Item

func (b *Bubble) GetRolesCmd() tea.Cmd {
	return func() tea.Msg {
		listItems := b.GetRoles()
		return GetRolesMsg(listItems)
	}
}

func (b *Bubble) CreateRoleCmd(rolename string) tea.Cmd {
	return func() tea.Msg {
		_, err := b.stardogClient.Role.Create(context.Background(), rolename)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully created role: %s", rolename),
		}
	}
}

func (b *Bubble) ForceDeleteRoleCmd() tea.Cmd {
	return func() tea.Msg {
		roleToDelete := b.GetSelectedRole()
		opts := stardog.DeleteRoleOptions{Force: true}
		_, err := b.stardogClient.Role.Delete(context.Background(), roleToDelete, &opts)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Forcefully deleted role: %s", roleToDelete),
		}
	}
}

func (b *Bubble) DeleteRoleCmd() tea.Cmd {
	return func() tea.Msg {
		roleToDelete := b.GetSelectedRole()
		opts := stardog.DeleteRoleOptions{Force: false}
		_, err := b.stardogClient.Role.Delete(context.Background(), roleToDelete, &opts)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully deleted role: %s", roleToDelete),
		}
	}
}
