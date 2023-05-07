package addpermissionprompt

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/noahgorstein/go-stardog/stardog"
)

func (b *Bubble) generateView() string {
	var sb strings.Builder
	sb.WriteString(b.Styles.secondaryHeader.Render("Add Permission"))
	sb.WriteRune('\n')
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.emphasisTextStyle.Render("esc"),
		b.Styles.normalTextStyle.Render("......back")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.emphasisTextStyle.Render("ctrl+b"),
		b.Styles.normalTextStyle.Render("...previous prompt")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.emphasisTextStyle.Render("ctrl+n"),
		b.Styles.normalTextStyle.Render("...next prompt")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.emphasisTextStyle.Render("enter"),
		b.Styles.normalTextStyle.Render("....make selection and jump to next prompt")))
	sb.WriteString(fmt.Sprintf("• %s%s\n",
		b.Styles.emphasisTextStyle.Render("space"),
		b.Styles.normalTextStyle.Render("....make selection")))

	if b.actionSelector.Width() > 0 {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(b.Styles.grey).Render(strings.Repeat("-", b.actionSelector.Width())))
	}

	sb.WriteRune('\n')

	currPermissionAction := ""
	a, ok := b.actionSelector.GetSelected().(stardog.PermissionAction)
	if ok {
		currPermissionAction = a.String()
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", b.Styles.selectionTextStyle.Render("Action:"), b.Styles.normalTextStyle.Render(currPermissionAction)))

	currPermissionResouceType := ""
	r, ok := b.resourceTypeSelector.GetSelected().(stardog.PermissionResourceType)
	if ok {
		currPermissionResouceType = r.String()
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", b.Styles.selectionTextStyle.Render("Resource Type:"), b.Styles.normalTextStyle.Render(currPermissionResouceType)))

	sb.WriteString(fmt.Sprintf("%s %s\n", b.Styles.selectionTextStyle.Render("Resource:"), b.Styles.normalTextStyle.Render(b.resourceInput.Value())))

	switch b.State {
	case SelectionActionState:
		sb.WriteString(b.actionSelector.View())
	case SelectingResourceTypeState:
		sb.WriteString(b.resourceTypeSelector.View())
	case SelectingResourceState:
		var resourceInput strings.Builder
		resourceInput.WriteString(b.Styles.resourceInputsPrompt.Render("Enter a resource: "))
		resourceInput.WriteRune('\n')
		resourceInput.WriteString(b.resourceInput.View())
		resourceInput.WriteRune('\n')

		var resourceView string
		resourceView = b.Styles.activeResourceInputs.Render(resourceInput.String())

		sb.WriteString(resourceView)

	case SubmitState:
		b.enableSubmitButton(true)
		sb.WriteString(b.submitButton)
	}

	return sb.String()
}

func (b Bubble) View() string {
	return b.generateView()
}
