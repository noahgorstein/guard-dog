package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/noahgorstein/stardog-go/stardog"
)

type activeView int

const (
	listView activeView = iota
	detailsView
)

const (
	columnKeyName      = "name"
	columnKeyEnabled   = "enabled"
	columnKeySuperuser = "superuser"
	columnKeyRoles     = "roles"
)

const (
	columnKeyAction       = "action"
	columnKeyResourceType = "resource type"
	columnKeyResource     = "resource"
	columnKeyExplicit     = "explicit"
)

// Bubble represents the properties of the UI.
type Bubble struct {
	width                int
	height               int
	list                 list.Model
	userDetailsTable     table.Model
	userPermissionsTable table.Model
	viewport             viewport.Model
	user                 stardog.User
	userDetails          stardog.GetUserDetailsResponse
	endpoint             string
	activeView           activeView
	columnSortKey        string
	statusBar            string
	connection           *stardog.ConnectionDetails
}

type item struct {
	title, desc string
}

func (i item) Title() string { return i.title }

func (i item) Description() string { return i.desc }

func (i item) FilterValue() string { return i.title }

// New creates a new instance of the UI.
func New() Bubble {

	stardogConnection := stardog.NewConnectionDetails("http://localhost:5820", "admin", "admin")
	users := stardog.GetUsers(*stardogConnection)

	items := []list.Item{}
	for _, user := range users {
		items = append(items, item{title: user.Name})
	}

	userDetailTable := table.New([]table.Column{
		table.NewColumn(columnKeyEnabled, "Enabled", 15),
		table.NewColumn(columnKeySuperuser, "Superuser", 15),
		table.NewFlexColumn(columnKeyRoles, "Roles", 1),
	})

	userPermissionsTable := table.New([]table.Column{
		table.NewColumn(columnKeyAction, "Action", 10).WithFiltered(true),
		table.NewColumn(columnKeyResourceType, "Resource Type", 20).WithFiltered(true),
		table.NewFlexColumn(columnKeyResource, "Resource", 1).WithFiltered(true),
		table.NewColumn(columnKeyExplicit, "Explicit Permission", 20).WithFiltered(true),
	}).Focused(true).HighlightStyle(lipgloss.NewStyle().Background(lipgloss.Color("#004B60")))

	itemDelegate := list.NewDefaultDelegate()
	itemDelegate.Styles.SelectedTitle.Foreground(lipgloss.Color("202"))
	itemDelegate.Styles.SelectedTitle.BorderLeftForeground(lipgloss.Color("202"))
	itemDelegate.Styles.SelectedDesc = lipgloss.NewStyle()
	itemDelegate.SetSpacing(0)

	userList := list.New(items, itemDelegate, 0, 0)

	b := Bubble{
		list:                 userList,
		userDetailsTable:     userDetailTable,
		userPermissionsTable: userPermissionsTable,
		activeView:           listView,
		connection:           stardogConnection,
	}

	b.list.Title = "Users"
	b.list.Styles.Title.MarginTop(1)
	return b
}
