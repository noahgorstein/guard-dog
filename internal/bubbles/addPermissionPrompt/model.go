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

	actionSelectorOptions := make(map[string]interface{})
	for k, v := range getStardogActions() {
		actionSelectorOptions[k] = v
	}

	resourceTypeSelectorOptions := make(map[string]interface{})
	for k, v := range getStardogResourceTypes() {
		resourceTypeSelectorOptions[k] = v
	}

	b := Bubble{
		State:                SelectionActionState,
		Styles:               styles,
		resourceInput:        resourceInput,
		actionSelector:       selector.New("Select an action: ", actionSelectorOptions, 0, 0),
		resourceTypeSelector: selector.New("Select a resource type: ", resourceTypeSelectorOptions, 0, 0),
	}

	b.actionSelector.SetIsActive(true)

	return b

}

func getStardogActions() map[string]stardog.PermissionAction {
	return map[string]stardog.PermissionAction{
		"ALL":     stardog.PermissionActionAll,
		"READ":    stardog.PermissionActionRead,
		"WRITE":   stardog.PermissionActionWrite,
		"CREATE":  stardog.PermissionActionCreate,
		"DELETE":  stardog.PermissionActionDelete,
		"GRANT":   stardog.PermissionActionGrant,
		"REVOKE":  stardog.PermissionActionRevoke,
		"EXECUTE": stardog.PermissionActionExecute,
	}
}

func getStardogResourceTypes() map[string]stardog.PermissionResourceType {
	return map[string]stardog.PermissionResourceType{
		"*":                    stardog.PermissionResourceTypeAll,
		"DB":                   stardog.PermissionResourceTypeDatabase,
		"USER":                 stardog.PermissionResourceTypeUser,
		"ROLE":                 stardog.PermissionResourceTypeRole,
		"METADATA":             stardog.PermissionResourceTypeMetadata,
		"NAMED-GRAPH":          stardog.PermissionResourceTypeNamedGraph,
		"VIRTUAL-GRAPH":        stardog.PermissionResourceTypeVirtualGraph,
		"DATA-SOURCE":          stardog.PermissionResourceTypeDataSource,
		"DBMS-ADMIN":           stardog.PermissionResourceTypeServeradmin,
		"ADMIN":                stardog.PermissionResourceTypeDatabaseAdmin,
		"SENSITIVE-PROPERTIES": stardog.PermissionResourceTypeSensitiveProperty,
		"STORED-QUERY":         stardog.PermissionResourceTypeStoredQuery,
	}
}

// return []string{"*", "USER", "ROLE", "DB", "METADATA", "NAMED-GRAPH",
//"VIRTUAL-GRAPH", "DATA-SOURCE", "DBMS-ADMIN", "ADMIN", "ICV-CONSTRAINTS", "SENSITIVE-PROPERTIES", "STORED-QUERY"}
