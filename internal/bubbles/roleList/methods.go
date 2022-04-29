package rolelist

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (b *Bubble) GetRoles() []list.Item {
	var items []list.Item
	roleList, _ := b.stardogClient.Security.ListRoles(context.Background())
	for _, role := range roleList {
		usersAssignedToRole, _ := b.stardogClient.Security.ListUsersAssignedToRole(context.Background(), role)
		desc := ""
		if len(usersAssignedToRole) == 1 {
			desc = fmt.Sprintf("%d user with role", len(usersAssignedToRole))
		} else {
			desc = fmt.Sprintf("%d users with role", len(usersAssignedToRole))

		}
		items = append(items,
			Item{
				title:       role,
				description: desc,
			})
	}
	return items
}

func (b *Bubble) IsFilterActive() bool {
	return b.list.FilterState() == list.Filtering
}

func (b *Bubble) Reset() {
	b.list.ResetSelected()
}

func (b *Bubble) resetAddRoleInputs() {
	b.focusIndex = 0
	for i := 0; i <= len(b.addRoleInputs)-1; i++ {
		b.addRoleInputs[i].Reset()
		b.addRoleInputs[i].Blur()
		b.addRoleInputs[i].PromptStyle = b.Styles.noStyle
	}
	b.addRoleInputs[0].PromptStyle = b.Styles.focusedStyle
	b.addRoleInputs[0].TextStyle = b.Styles.focusedStyle
	b.addRoleInputs[0].Focus()
}

func (b *Bubble) updateAddUserInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(b.addRoleInputs))

	for i := range b.addRoleInputs {
		b.addRoleInputs[i], cmds[i] = b.addRoleInputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (b *Bubble) enableSubmitButton(enabled bool) {
	if enabled {
		b.submitButton = b.Styles.focusedStyle.Copy().Render("[ Submit ]")
	} else {
		b.submitButton = fmt.Sprintf("[ %s ]", b.Styles.blurredStyle.Render("Submit"))
	}
}

func (b *Bubble) GetSelectedRole() string {
	if b.list.SelectedItem() != nil {
		return b.list.SelectedItem().(Item).Title()
	}
	return ""
}

func (b *Bubble) SetIsActive(active bool) {
	b.active = active
}

func (b *Bubble) SetSize(width, height int) {
	horizontal, vertical := b.Styles.listStyle.GetFrameSize()
	b.list.Styles.StatusBar.Width(width - horizontal)

	inputHeight := 0
	for i := 0; i <= len(b.addRoleInputs)-1; i++ {
		inputHeight += lipgloss.Height(b.addRoleInputs[i].View())
	}

	b.list.SetSize(
		width-horizontal-vertical,
		height-vertical-inputHeight-b.Styles.listStyle.GetVerticalFrameSize(),
	)
}

func (b *Bubble) SetListTitleBackgroundColor(color lipgloss.AdaptiveColor) {
	b.list.Styles.Title.Background(color)
}

func (b *Bubble) SetListTitleForegroundColor(color lipgloss.AdaptiveColor) {
	b.list.Styles.Title.Foreground(color)
}
