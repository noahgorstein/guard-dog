package addpermissionprompt

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/noahgorstein/go-stardog/stardog"
	"github.com/noahgorstein/guard-dog/internal/bubbles/selector"
)

type State int

const (
	SelectionActionState = iota
	SelectingResourceTypeState
	SelectingResourceState
	SubmitState
)

type Bubble struct {
	State  State
	Styles Styles

	resourceInput        textinput.Model
	actionSelector       selector.Model
	resourceTypeSelector selector.Model
	submitButton         string
	focusIndex           int

	Permission stardog.Permission
}

func New() Bubble {
	styles := DefaultStyles()

	resourceInput := textinput.New()
	resourceInput.CharLimit = 64
	resourceInput.Placeholder = "Resource"

	b := Bubble{
		State:                SelectionActionState,
		Styles:               styles,
		resourceInput:        resourceInput,
		actionSelector:       selector.New("Select an action: ", getStardogActions(), 0, 0),
		resourceTypeSelector: selector.New("Select a resource type: ", getStardogResourceTypes(), 0, 0),
	}

	b.actionSelector.SetIsActive(true)

	return b

}

func getStardogActions() []string {
	return []string{"ALL", "READ", "WRITE", "CREATE", "DELETE", "GRANT", "REVOKE", "EXECUTE"}
}

func getStardogResourceTypes() []string {
	return []string{"*", "USER", "ROLE", "DB", "METADATA", "NAMED-GRAPH",
		"VIRTUAL-GRAPH", "DATA-SOURCE", "DBMS-ADMIN", "ADMIN", "ICV-CONSTRAINTS", "SENSITIVE-PROPERTIES", "STORED-QUERY"}
}
