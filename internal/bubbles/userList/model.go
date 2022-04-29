package userlist

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/noahgorstein/go-stardog/stardog"
)

type State int

const (
	IdleState State = iota
	AddUserState
	ChangeUserPasswordState
)

type Bubble struct {
	Styles                   Styles
	State                    State
	list                     list.Model
	active                   bool
	width                    int
	height                   int
	stardogClient            stardog.Client
	addUserInputs            []textinput.Model
	changeUserPasswordInputs []textinput.Model
	divider                  string
	focusIndex               int
	submitButton             string
}

func New(stardogClient stardog.Client) Bubble {
	styles := DefaultStyles()

	items := []list.Item{}
	itemDelegate := list.NewDefaultDelegate()

	itemDelegate.Styles.NormalTitle.Bold(true).
		Foreground(lipgloss.AdaptiveColor{
			Light: string(nord0),
			Dark:  string(nord4),
		})
	itemDelegate.Styles.SelectedTitle.
		Foreground(nord12).
		BorderLeftForeground(nord10).Bold(true)

	itemDelegate.Styles.NormalDesc.Foreground(nord8)
	itemDelegate.Styles.SelectedDesc.
		BorderLeftForeground(nord10).Foreground(nord8)

	userList := list.New(items, itemDelegate, 0, 0)
	userList.Title = "Users"
	userList.Styles.Title.Bold(true).Italic(true)
	userList.SetShowHelp(false)
	userList.DisableQuitKeybindings()
	userList.SetStatusBarItemName("user", "users")

	b := Bubble{
		Styles:                   styles,
		State:                    IdleState,
		list:                     userList,
		stardogClient:            stardogClient,
		addUserInputs:            make([]textinput.Model, 2),
		changeUserPasswordInputs: make([]textinput.Model, 2),
	}

	initialUserList, _ := b.GetStardogUsersWithDetails()
	b.list.SetItems(initialUserList)

	var addUserInput textinput.Model
	for i := range b.addUserInputs {
		addUserInput = textinput.New()
		addUserInput.CharLimit = 32

		switch i {
		case 0:
			addUserInput.Placeholder = "Username"
			addUserInput.Focus()
			addUserInput.PromptStyle = b.Styles.focusedStyle
			addUserInput.TextStyle = b.Styles.focusedStyle
		case 1:
			addUserInput.Placeholder = "Password"
			addUserInput.EchoMode = textinput.EchoPassword
			addUserInput.EchoCharacter = '•'
		}
		b.addUserInputs[i] = addUserInput
	}

	var changeUserPasswordInput textinput.Model
	for i := range b.changeUserPasswordInputs {
		changeUserPasswordInput = textinput.New()
		changeUserPasswordInput.CharLimit = 32

		switch i {
		case 0:
			changeUserPasswordInput.Placeholder = "Enter password"
			changeUserPasswordInput.EchoMode = textinput.EchoPassword
			changeUserPasswordInput.EchoCharacter = '•'
			changeUserPasswordInput.Focus()
			changeUserPasswordInput.PromptStyle = b.Styles.focusedStyle
			changeUserPasswordInput.TextStyle = b.Styles.focusedStyle

		case 1:
			changeUserPasswordInput.Placeholder = "Confirm password"
			changeUserPasswordInput.EchoMode = textinput.EchoPassword
			changeUserPasswordInput.EchoCharacter = '•'
		}
		b.changeUserPasswordInputs[i] = changeUserPasswordInput
	}

	return b

}
