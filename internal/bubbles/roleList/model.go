package rolelist

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/noahgorstein/go-stardog/stardog"
)

type State int

const (
	IdleState State = iota
	AddRoleState
)

type Bubble struct {
	State         State
	active        bool
	Styles        Styles
	list          list.Model
	stardogClient stardog.Client
	focusIndex    int
	addRoleInputs []textinput.Model
	submitButton  string
}

func New(stardogClient stardog.Client) Bubble {
	styles := DefaultStyles()

	items := []list.Item{}
	itemDelegate := list.NewDefaultDelegate()
	itemDelegate.SetSpacing(1)

	itemDelegate.Styles.NormalTitle.Bold(true)
	itemDelegate.Styles.SelectedTitle.
		Foreground(nord12).
		BorderLeftForeground(nord10).Bold(true)

	itemDelegate.Styles.NormalDesc.Foreground(lipgloss.AdaptiveColor{
		Light: string(nord8),
		Dark:  string(nord7),
	})
	itemDelegate.Styles.SelectedDesc.
		BorderLeftForeground(nord10).Foreground(nord14)

	roleList := list.New(items, itemDelegate, 0, 0)
	roleList.Title = "Roles"
	roleList.Styles.Title.Bold(true).Italic(true)
	roleList.SetShowHelp(false)
	roleList.DisableQuitKeybindings()
	roleList.SetStatusBarItemName("role", "roles")

	b := Bubble{
		State:         IdleState,
		Styles:        styles,
		list:          roleList,
		stardogClient: stardogClient,
		addRoleInputs: make([]textinput.Model, 1),
	}
	b.list.SetItems(b.GetRoles())

	b.addRoleInputs[0] = textinput.New()
	b.addRoleInputs[0].Placeholder = "Rolename"
	b.addRoleInputs[0].Focus()
	b.addRoleInputs[0].PromptStyle = b.Styles.focusedStyle
	b.addRoleInputs[0].TextStyle = b.Styles.focusedStyle

	return b
}
