package addpermissionprompt

import (
	"fmt"
	"strings"
)

func (b *Bubble) Reset() {
	b.State = SelectionActionState

	b.actionSelector.Reset()
	b.actionSelector.SetIsActive(true)

	b.resourceTypeSelector.Reset()
	b.resourceTypeSelector.SetIsActive(false)

	b.resetResourceInput()
}

func (b *Bubble) enableSubmitButton(enabled bool) {
	if enabled {
		b.submitButton = b.Styles.focusedStyle.Copy().Render("[ Submit ]")
	} else {
		b.submitButton = fmt.Sprintf("[ %s ]", b.Styles.blurredStyle.Render("Submit"))
	}
}

func (b *Bubble) resetResourceInput() {
	b.focusIndex = 0
	b.resourceInput.Reset()
	b.resourceInput.Blur()
	b.resourceInput.PromptStyle = b.Styles.noStyle
}

func (b *Bubble) SetSize(width int) {

	b.actionSelector.SetWidth(width - 5)
	b.resourceTypeSelector.SetWidth(width - 5)
	b.Styles.inactiveResourceInputs.Width(width - 5)
	b.Styles.activeResourceInputs.Width(width - 5)
}

func (b *Bubble) updatePermissionAction(action string) {
	b.Permission.Action = action
}

func (b *Bubble) updatePermissionResourceType(resourceType string) {
	b.Permission.ResourceType = resourceType
}

func (b *Bubble) updatePermissionResource(resource string) {
	permissionResourceArr := strings.Split(resource, "\\")

	b.Permission.Resource = permissionResourceArr
}

func (b *Bubble) updateResourcePromptPlaceholder(selectedResourceType string) {

	switch selectedResourceType {
	case "DB", "METADATA", "ADMIN", "ICV-CONSTRAINTS", "SENSITIVE-PROPERTIES", "STORED-QUERY":
		b.resourceInput.Placeholder = "myDatabase"
	case "ROLE":
		b.resourceInput.Placeholder = "myRole"
	case "USER":
		b.resourceInput.Placeholder = "myUser"
	case "NAMED-GRAPH":
		b.resourceInput.Placeholder = "myDatabase\\https://my.graph.com"
	case "VIRTUAL-GRAPH":
		b.resourceInput.Placeholder = "virtual://myVirtualGraph"
	case "DATA-SOURCE":
		b.resourceInput.Placeholder = "data-source://myDataSource"
	case "DBMS-ADMIN":
		b.resourceInput.Placeholder = "metrics"
	}

}
