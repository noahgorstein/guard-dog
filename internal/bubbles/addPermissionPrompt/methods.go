package addpermissionprompt

import (
	"fmt"
	"strings"

	"github.com/noahgorstein/go-stardog/stardog"
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

func (b *Bubble) updatePermissionAction(action stardog.PermissionAction) {
	b.Permission.Action = action
}

func (b *Bubble) updatePermissionResourceType(resourceType stardog.PermissionResourceType) {
	b.Permission.ResourceType = resourceType
}

func (b *Bubble) updatePermissionResource(resource string) {
	permissionResourceArr := strings.Split(resource, "\\")

	b.Permission.Resource = permissionResourceArr
}

func (b *Bubble) updateResourcePromptPlaceholder(selectedResourceType stardog.PermissionResourceType) {

	switch selectedResourceType {
	case stardog.PermissionResourceTypeDatabase, stardog.PermissionResourceTypeMetadata, stardog.PermissionResourceTypeServeradmin, stardog.PermissionResourceTypeSensitiveProperty, stardog.PermissionResourceTypeStoredQuery:
		b.resourceInput.Placeholder = "myDatabase"
	case stardog.PermissionResourceTypeRole:
		b.resourceInput.Placeholder = "myRole"
	case stardog.PermissionResourceTypeUser:
		b.resourceInput.Placeholder = "myUser"
	case stardog.PermissionResourceTypeNamedGraph:
		b.resourceInput.Placeholder = "myDatabase\\https://my.graph.com"
	case stardog.PermissionResourceTypeVirtualGraph:
		b.resourceInput.Placeholder = "virtual://myVirtualGraph"
	case stardog.PermissionResourceTypeDataSource:
		b.resourceInput.Placeholder = "data-source://myDataSource"
	case stardog.PermissionResourceTypeDatabaseAdmin:
		b.resourceInput.Placeholder = "metrics"
	}

}
