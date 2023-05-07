package addpermissionprompt

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/noahgorstein/go-stardog/stardog"
)

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch b.State {
		case SelectionActionState:
			b.actionSelector, cmd = b.actionSelector.Update(msg)
			cmds = append(cmds, cmd)

			if key.Matches(msg, previousKey) {
				a, ok := b.actionSelector.GetSelected().(stardog.PermissionAction)
				if ok {
					b.updatePermissionAction(a)
				}
				b.State = SubmitState
				b.actionSelector.SetIsActive(false)
				b.focusIndex = 1
			}

			if msg.Type == tea.KeyEnter || key.Matches(msg, nextKey) {
				a, ok := b.actionSelector.GetSelected().(stardog.PermissionAction)
				if ok {
					b.updatePermissionAction(a)
				}
				b.State = SubmitState
				b.State = SelectingResourceTypeState
				b.resourceTypeSelector.SetIsActive(true)
				b.actionSelector.SetIsActive(false)
			}
		case SelectingResourceTypeState:

			b.resourceTypeSelector, cmd = b.resourceTypeSelector.Update(msg)
			cmds = append(cmds, cmd)

			if key.Matches(msg, previousKey) {
				rt, ok := b.resourceTypeSelector.GetSelected().(stardog.PermissionResourceType)
				if ok {
					b.updatePermissionResourceType(rt)
				}
				b.State = SelectionActionState
				b.resourceTypeSelector.SetIsActive(false)
				b.actionSelector.SetIsActive(true)
			}

			if msg.Type == tea.KeyEnter || key.Matches(msg, nextKey) {
				rt, ok := b.resourceTypeSelector.GetSelected().(stardog.PermissionResourceType)
				if ok {
					b.updatePermissionResourceType(rt)
					b.updateResourcePromptPlaceholder(rt)
				}

				b.State = SelectingResourceState
				b.resourceTypeSelector.SetIsActive(false)

				b.resourceInput.PromptStyle = b.Styles.focusedStyle
				b.resourceInput.TextStyle = b.Styles.focusedStyle
				b.resourceInput.Focus()
			}
		case SelectingResourceState:

			b.resourceInput, cmd = b.resourceInput.Update(msg)
			cmds = append(cmds, cmd)

			if key.Matches(msg, previousKey) {
				b.updatePermissionResource(b.resourceInput.Value())
				b.State = SelectingResourceTypeState
				b.resourceInput.PromptStyle = b.Styles.noStyle
				b.resourceInput.TextStyle = b.Styles.noStyle
				b.resourceTypeSelector.SetIsActive(true)
			}

			if msg.Type == tea.KeyEnter || key.Matches(msg, nextKey) {
				b.updatePermissionResource(b.resourceInput.Value())
				b.resourceInput.PromptStyle = b.Styles.noStyle
				b.resourceInput.TextStyle = b.Styles.noStyle
				b.State = SubmitState
			}

		case SubmitState:

			if key.Matches(msg, previousKey) {
				b.State = SelectingResourceState
				b.resourceInput.PromptStyle = b.Styles.focusedStyle
				b.resourceInput.TextStyle = b.Styles.focusedStyle
				b.resourceInput.Focus()
			}

			if key.Matches(msg, nextKey) {
				b.focusIndex = 0
				b.State = SelectionActionState
				b.actionSelector.SetIsActive(true)
			}
		}
	}

	return b, tea.Batch(cmds...)
}
