package userlist

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
)

type SuccessMsg struct {
	Message string
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type GetUsersMsg []list.Item

func (b *Bubble) GetUsersCmd() tea.Cmd {
	return func() tea.Msg {
		listItems, err := b.GetStardogUsersWithDetails()
		if err != nil {
			os.Exit(1)
		}
		return GetUsersMsg(listItems)
	}
}

func (b *Bubble) CreateUserCmd(username, password string) tea.Cmd {
	return func() tea.Msg {
		_, err := b.stardogClient.User.Create(context.Background(), username, password)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully created user: %s", username),
		}
	}
}

func (b *Bubble) DeleteUserCmd() tea.Cmd {
	return func() tea.Msg {
		_, err := b.stardogClient.User.Delete(context.Background(), b.GetSelectedUser().Username())
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully deleted user: %s", b.GetSelectedUser().Username()),
		}

	}
}

func (b *Bubble) DisableUserCmd() tea.Cmd {
	return func() tea.Msg {
		_, err := b.stardogClient.User.Disable(context.Background(), b.GetSelectedUser().Username())
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully disabled user: %s", b.GetSelectedUser().Username()),
		}
	}
}

func (b *Bubble) EnableUserCmd() tea.Cmd {
	return func() tea.Msg {
		_, err := b.stardogClient.User.Enable(context.Background(), b.GetSelectedUser().Username())
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully enabled user: %s", b.GetSelectedUser().Username()),
		}
	}
}

func (b *Bubble) ChangeUserPasswordCmd(newPassword, confirmedPassword string) tea.Cmd {
	return func() tea.Msg {

		if newPassword != confirmedPassword {
			return errMsg{
				err: fmt.Errorf("passwords must match"),
			}
		}

		_, err := b.stardogClient.User.ChangePassword(context.Background(), b.GetSelectedUser().Username(), newPassword)
		if err != nil {
			return errMsg{
				err: err,
			}
		}
		return SuccessMsg{
			Message: fmt.Sprintf("Successfully changed password for %s", b.GetSelectedUser().Username()),
		}
	}
}
