package userlist

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (b *Bubble) GetStardogUsersWithDetails() ([]list.Item, error) {
	var items []list.Item
	userList, err := b.stardogClient.Security.ListUsers(context.Background())
	if err != nil {
		return nil, err
	}
	for _, user := range userList {

		userDetails, _ := b.stardogClient.Security.GetUserDetails(context.Background(), user)
		items = append(items,
			Item{
				title:     user,
				username:  user,
				superuser: userDetails.Superuser,
				enabled:   userDetails.Enabled,
				roles:     userDetails.Roles,
			})
	}
	return items, nil
}

func (b *Bubble) IsFilterActive() bool {
	return b.list.FilterState() == list.Filtering
}

func (b *Bubble) Reset() {
	b.list.ResetSelected()
}

func (b *Bubble) GetSelectedUser() Item {
	if b.list.SelectedItem() != nil {
		return b.list.SelectedItem().(Item)
	}
	return Item{}
}

func (b *Bubble) GetCurrentUserRoles() []string {
	return strings.Split(b.list.SelectedItem().(Item).Description(), ",")
}

func (b *Bubble) enableSubmitButton(enabled bool) {
	if enabled {
		b.submitButton = b.Styles.focusedStyle.Copy().Render("[ Submit ]")
	} else {
		b.submitButton = fmt.Sprintf("[ %s ]", b.Styles.blurredStyle.Render("Submit"))
	}
}

func (b *Bubble) SetIsActive(active bool) {
	b.active = active
}

func (b *Bubble) SetSize(width, height int) {
	horizontal, vertical := b.Styles.listStyle.GetFrameSize()
	b.list.Styles.StatusBar.Width(width - horizontal)

	inputHeight := 1
	for i := 0; i <= len(b.addUserInputs)-1; i++ {
		inputHeight += lipgloss.Height(b.addUserInputs[i].View())
	}
	b.list.SetSize(
		width-horizontal-vertical,
		height-vertical-inputHeight-b.Styles.listStyle.GetVerticalFrameSize(),
	)
}

func (b *Bubble) resetAddUserInputs() {
	b.focusIndex = 0
	for i := 0; i <= len(b.addUserInputs)-1; i++ {
		b.addUserInputs[i].Reset()
		b.addUserInputs[i].Blur()
		b.addUserInputs[i].PromptStyle = b.Styles.noStyle
	}
	b.addUserInputs[0].PromptStyle = b.Styles.focusedStyle
	b.addUserInputs[0].TextStyle = b.Styles.focusedStyle
	b.addUserInputs[0].Focus()
}

func (b *Bubble) resetChangeUserPasswordInputs() {
	b.focusIndex = 0
	for i := 0; i <= len(b.changeUserPasswordInputs)-1; i++ {
		b.changeUserPasswordInputs[i].Reset()
		b.changeUserPasswordInputs[i].Blur()
		b.changeUserPasswordInputs[i].PromptStyle = b.Styles.noStyle
	}
	b.changeUserPasswordInputs[0].PromptStyle = b.Styles.focusedStyle
	b.changeUserPasswordInputs[0].TextStyle = b.Styles.focusedStyle
	b.changeUserPasswordInputs[0].Focus()
}

func (b *Bubble) updateAddUserInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(b.addUserInputs))

	for i := range b.addUserInputs {
		b.addUserInputs[i], cmds[i] = b.addUserInputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (b *Bubble) updateChangeUserPasswordInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(b.changeUserPasswordInputs))

	for i := range b.changeUserPasswordInputs {
		b.changeUserPasswordInputs[i], cmds[i] = b.changeUserPasswordInputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (b *Bubble) SetListTitleBackgroundColor(color lipgloss.AdaptiveColor) {
	b.list.Styles.Title.Background(color)
}

func (b *Bubble) SetListTitleForegroundColor(color lipgloss.AdaptiveColor) {
	b.list.Styles.Title.Foreground(color)
}
