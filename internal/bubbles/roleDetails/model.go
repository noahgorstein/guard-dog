package roledetails

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/noahgorstein/go-stardog/stardog"
	permissionPrompt "github.com/noahgorstein/guard-dog/internal/bubbles/addPermissionPrompt"
)

type permissionTableColumnKeys struct {
	action       string
	resourceType string
	resource     string
}

type State int

const (
	IdleState State = iota
	AddingRolePermissionState
	TableState
)

type Bubble struct {
	State                     State
	active                    bool
	stardogClient             stardog.Client
	Styles                    Styles
	viewport                  viewport.Model
	permissionTableColumnKeys permissionTableColumnKeys
	permissionsTable          table.Model
	selectedRole              string
	usersAssignedToRole       []string

	addPermissionsPrompt permissionPrompt.Bubble
}

func New(stardogClient stardog.Client) Bubble {

	styles := DefaultStyles()

	viewport := viewport.New(0, 0)
	viewport.Style = styles.InactiveViewportStyle

	permissionTableColumnKeys := permissionTableColumnKeys{
		action:       "action",
		resourceType: "resource type",
		resource:     "resource",
	}

	permissionsTable := table.New([]table.Column{
		table.NewColumn(permissionTableColumnKeys.action, "Action", 10).WithFiltered(true),
		table.NewColumn(permissionTableColumnKeys.resourceType, "Resource Type", 20).WithFiltered(true),
		table.NewFlexColumn(permissionTableColumnKeys.resource, "Resource", 1).WithFiltered(true),
	}).Focused(true).Filtered(true).Border(styles.PermissionsTableBorder).WithPageSize(10)

	permissionsTable = permissionsTable.
		HeaderStyle(styles.PermissionsTableHeaderStyle).
		WithBaseStyle(lipgloss.NewStyle().
			BorderForeground(styles.grey).
			Foreground(lipgloss.AdaptiveColor{
				Light: string(styles.black),
				Dark:  string(styles.white),
			}))

	b := Bubble{
		State:                     IdleState,
		Styles:                    styles,
		stardogClient:             stardogClient,
		viewport:                  viewport,
		permissionTableColumnKeys: permissionTableColumnKeys,
		permissionsTable:          permissionsTable,

		addPermissionsPrompt: permissionPrompt.New(),
	}

	return b
}
