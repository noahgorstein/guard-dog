package userdetails

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/noahgorstein/go-stardog/stardog"
	addpermissionprompt "github.com/noahgorstein/guard-dog/internal/bubbles/addPermissionPrompt"
	selector "github.com/noahgorstein/guard-dog/internal/bubbles/selector"
)

type State int

const (
	IdleState State = iota
	TableState
	AddingUserPermission
	AssigningRoleState
	RemovingRoleState
)

type permissionTableColumnKeys struct {
	action       string
	resourceType string
	resource     string
}

type Bubble struct {
	State                     State
	Styles                    Styles
	viewport                  viewport.Model
	active                    bool
	permissionTableColumnKeys permissionTableColumnKeys
	permissionsTable          table.Model
	stardogClient             stardog.Client
	selectedUser              string
	assignedRoles             []string
	addPermissionsPrompt      addpermissionprompt.Bubble
	loggedInUser              string

	addRoleSelector    selector.Model
	removeRoleSelector selector.Model
	focusIndex         int
	submitButton       string
}

func New(stardogClient stardog.Client, loggedInUser string) Bubble {
	styles := DefaultStyles()

	viewport := viewport.New(0, 0)
	viewport.Style = styles.InactiveViewportStyle
	viewport.MouseWheelEnabled = false

	permissionTableColumnKeys := permissionTableColumnKeys{
		action:       "action",
		resourceType: "resource type",
		resource:     "resource",
	}

	permissionsTable := table.New([]table.Column{
		table.NewColumn(permissionTableColumnKeys.action, "Action", 10).WithFiltered(true),
		table.NewColumn(permissionTableColumnKeys.resourceType, "Resource Type", 20).WithFiltered(true),
		table.NewFlexColumn(permissionTableColumnKeys.resource, "Resource", 1).WithFiltered(true),
	}).Focused(true).Filtered(true).Border(styles.permissionsTableBorder).WithPageSize(10)

	permissionsTable = permissionsTable.
		HeaderStyle(styles.permissionsTableHeaderStyle).
		WithBaseStyle(lipgloss.NewStyle().
			BorderForeground(styles.grey).
			Foreground(lipgloss.AdaptiveColor{
				Light: string(nord0),
				Dark:  string(nord6),
			}))

	b := Bubble{
		active:                    false,
		Styles:                    styles,
		State:                     IdleState,
		stardogClient:             stardogClient,
		viewport:                  viewport,
		permissionTableColumnKeys: permissionTableColumnKeys,
		permissionsTable:          permissionsTable,
		addRoleSelector:           selector.New("Select a role: ", make(map[string]interface{}), 0, 0),
		removeRoleSelector:        selector.New("Remove a role: ", make(map[string]interface{}), 0, 0),
		addPermissionsPrompt:      addpermissionprompt.New(),
		loggedInUser:              loggedInUser,
	}

	return b
}
