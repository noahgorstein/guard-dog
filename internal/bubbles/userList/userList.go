package userlist

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/noahgorstein/stardog-go/internal/config"
	"github.com/noahgorstein/stardog-go/stardog"
)

var (
	f, _ = tea.LogToFile("debug.log", "debug")
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	noStyle      = lipgloss.NewStyle()
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedButton     = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton     = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	addUserInputStyle = lipgloss.NewStyle().PaddingTop(1)
	bubbleStyle       = lipgloss.NewStyle().
				PaddingLeft(1).
				PaddingRight(1).
				BorderStyle(lipgloss.RoundedBorder())
)

type sessionState int

const (
	idleState sessionState = iota
	addUserState
	changeUserPasswordState
)

type Bubble struct {
	state                    sessionState
	list                     list.Model
	active                   bool
	width                    int
	height                   int
	connection               stardog.ConnectionDetails
	addUserInputs            []textinput.Model
	changeUserPasswordInputs []textinput.Model

	focusIndex int
}

func New(config config.Config) Bubble {
	connectionDetails := stardog.NewConnectionDetails(config.Endpoint, config.Username, config.Password)
	users := stardog.GetUsers(*connectionDetails)

	items := []list.Item{}
	for _, user := range users {
		items = append(items, item{title: user.Name})
	}

	itemDelegate := list.NewDefaultDelegate()
	itemDelegate.Styles.SelectedTitle.Foreground(lipgloss.Color("202"))
	itemDelegate.Styles.SelectedTitle.BorderLeftForeground(lipgloss.Color("202"))
	itemDelegate.Styles.SelectedDesc = lipgloss.NewStyle()
	itemDelegate.SetSpacing(0)

	userList := list.New(items, itemDelegate, 0, 0)
	userList.Title = "Users"
	userList.Styles.Title.Bold(true).Italic(true).Background(lipgloss.Color("202"))

	b := Bubble{
		state:                    idleState,
		list:                     userList,
		connection:               *connectionDetails,
		addUserInputs:            make([]textinput.Model, 2),
		changeUserPasswordInputs: make([]textinput.Model, 2),
	}

	var addUserInput textinput.Model
	for i := range b.addUserInputs {
		addUserInput = textinput.New()
		addUserInput.CharLimit = 32

		switch i {
		case 0:
			addUserInput.Placeholder = "Username"
			addUserInput.Focus()
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
		case 1:
			changeUserPasswordInput.Placeholder = "Confirm password"
			changeUserPasswordInput.EchoMode = textinput.EchoPassword
			changeUserPasswordInput.EchoCharacter = '•'
		}
		b.changeUserPasswordInputs[i] = changeUserPasswordInput
	}

	return b

}

func (b Bubble) Init() tea.Cmd {
	return nil
}

func (b *Bubble) GetCurrentUser() string {
	selectedItem := b.list.SelectedItem()
	if selectedItem != nil {
		return selectedItem.(item).title
	} else {
		return ""
	}
}

func (b *Bubble) SetIsActive(active bool) {
	b.active = active
}

func (b *Bubble) SetSize(width, height int) {
	horizontal, vertical := bubbleStyle.GetFrameSize()
	b.list.Styles.StatusBar.Width(width - horizontal)

	inputHeight := 1
	for i := 0; i <= len(b.addUserInputs)-1; i++ {
		inputHeight += lipgloss.Height(b.addUserInputs[i].View())
	}
	b.list.SetSize(
		width-horizontal-vertical,
		height-vertical-inputHeight-addUserInputStyle.GetVerticalFrameSize(),
	)
}

func (b *Bubble) resetAddUserInputs() {
	for i := 0; i <= len(b.addUserInputs)-1; i++ {
		b.addUserInputs[i].Reset()
	}
}

func (b *Bubble) resetChangeUserPasswordInputs() {
	for i := 0; i <= len(b.changeUserPasswordInputs)-1; i++ {
		b.changeUserPasswordInputs[i].Reset()
	}
}

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
	case GetStardogUsersMsg:
		if msg != nil {
			cmd = b.list.SetItems(msg)
			cmds = append(cmds, cmd)
		}
	case tea.KeyMsg:

		if msg.Type == tea.KeyCtrlC {
			return b, tea.Quit
		}

		switch b.state {
		case addUserState:
			if msg.Type == tea.KeyEsc {
				b.state = idleState
				b.resetAddUserInputs()
				return b, nil
			}
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeyUp || msg.Type == tea.KeyDown {
				s := msg.String()

				if s == "enter" && b.focusIndex == len(b.addUserInputs) {
					b.state = idleState
					createStardogUserCmd := b.CreateStardogCmd(b.addUserInputs[0].Value(), b.addUserInputs[1].Value())
					getStardogUsersCmd := b.GetStardogUsersCmd()
					b.addUserInputs[0].Reset()
					b.addUserInputs[1].Reset()
					cmds = append(cmds, tea.Sequentially(createStardogUserCmd, getStardogUsersCmd))
				}

				if s == "up" {
					b.focusIndex--
				} else {
					b.focusIndex++
				}

				if b.focusIndex > len(b.addUserInputs) {
					b.focusIndex = 0
				} else if b.focusIndex < 0 {
					b.focusIndex = len(b.addUserInputs)
				}

				cmds := make([]tea.Cmd, len(b.addUserInputs))
				for i := 0; i <= len(b.addUserInputs)-1; i++ {
					if i == b.focusIndex {
						// Set focused state
						cmds[i] = b.addUserInputs[i].Focus()
						b.addUserInputs[i].PromptStyle = focusedStyle
						b.addUserInputs[i].TextStyle = focusedStyle
						continue
					}
					// Remove focused state
					b.addUserInputs[i].Blur()
					b.addUserInputs[i].PromptStyle = noStyle
					b.addUserInputs[i].TextStyle = noStyle
				}
			}

		case changeUserPasswordState:
			if msg.Type == tea.KeyEsc {
				b.state = idleState
				b.resetChangeUserPasswordInputs()
				return b, nil
			}
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeyUp || msg.Type == tea.KeyDown {
				s := msg.String()

				if s == "enter" && b.focusIndex == len(b.changeUserPasswordInputs) {
					b.state = idleState
					changeUserPasswordCmd := b.ChangeUserPasswordCmd(b.GetCurrentUser(), b.changeUserPasswordInputs[1].Value())
					// createStardogUserCmd := b.CreateStardogCmd(b.addUserInputs[0].Value(), b.addUserInputs[1].Value())
					getStardogUsersCmd := b.GetStardogUsersCmd()
					b.changeUserPasswordInputs[0].Reset()
					b.changeUserPasswordInputs[1].Reset()
					cmds = append(cmds, tea.Sequentially(changeUserPasswordCmd, getStardogUsersCmd))
				}

				if s == "up" {
					b.focusIndex--
				} else {
					b.focusIndex++
				}

				if b.focusIndex > len(b.changeUserPasswordInputs) {
					b.focusIndex = 0
				} else if b.focusIndex < 0 {
					b.focusIndex = len(b.changeUserPasswordInputs)
				}

				cmds := make([]tea.Cmd, len(b.changeUserPasswordInputs))
				for i := 0; i <= len(b.changeUserPasswordInputs)-1; i++ {
					if i == b.focusIndex {
						// Set focused state
						cmds[i] = b.changeUserPasswordInputs[i].Focus()
						b.changeUserPasswordInputs[i].PromptStyle = focusedStyle
						b.changeUserPasswordInputs[i].TextStyle = focusedStyle
						continue
					}
					// Remove focused state
					b.changeUserPasswordInputs[i].Blur()
					b.changeUserPasswordInputs[i].PromptStyle = noStyle
					b.changeUserPasswordInputs[i].TextStyle = noStyle
				}
			}
		}

		if msg.Type == tea.KeyCtrlA {
			username := b.addUserInputs[0]
			username.Focus()
			username.PromptStyle = focusedStyle
			b.state = addUserState
			return b, textinput.Blink
		}

		if msg.Type == tea.KeyCtrlP {
			password := b.changeUserPasswordInputs[0]
			password.Focus()
			password.PromptStyle = focusedStyle
			b.state = changeUserPasswordState
			return b, textinput.Blink
		}

		if msg.Type == tea.KeyCtrlD {
			if b.state == idleState {
				deleteStardogUserCmd := b.DeleteStardogUserCmd(b.GetCurrentUser())
				getStardogUsers := b.GetStardogUsersCmd()
				cmds = append(cmds, tea.Sequentially(deleteStardogUserCmd, getStardogUsers))

			}
		}

	}

	if b.active {
		switch b.state {
		case idleState:
			b.list, cmd = b.list.Update(msg)
			cmds = append(cmds, cmd)
		case addUserState:
			cmd := b.updateAddUserInputs(msg)
			cmds = append(cmds, cmd)
		case changeUserPasswordState:
			cmd := b.updateChangeUserPasswordInputs(msg)
			cmds = append(cmds, cmd)
		}
	}

	return b, tea.Batch(cmds...)
}

func (b *Bubble) updateAddUserInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(b.addUserInputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range b.addUserInputs {
		b.addUserInputs[i], cmds[i] = b.addUserInputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (b *Bubble) updateChangeUserPasswordInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(b.changeUserPasswordInputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range b.changeUserPasswordInputs {
		b.changeUserPasswordInputs[i], cmds[i] = b.changeUserPasswordInputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (b Bubble) View() string {
	if b.active {
		bubbleStyle = bubbleStyle.BorderForeground(lipgloss.Color("33"))
	} else {
		bubbleStyle = bubbleStyle.BorderForeground(lipgloss.Color("#FAFAFA"))
	}

	var inputView string
	var builder strings.Builder

	switch b.state {
	case idleState:
		inputView = "\n\n"
	case addUserState:

		for i := range b.addUserInputs {
			builder.WriteString(b.addUserInputs[i].View())
			if i < len(b.addUserInputs)-1 {
				builder.WriteRune('\n')

			}
		}
		button := &blurredButton
		if b.focusIndex == len(b.addUserInputs) {
			button = &focusedButton
		}
		builder.WriteString("\n" + *button)

		inputView = builder.String()
	case changeUserPasswordState:
		for i := range b.changeUserPasswordInputs {
			builder.WriteString(b.changeUserPasswordInputs[i].View())
			if i < len(b.changeUserPasswordInputs)-1 {
				builder.WriteRune('\n')

			}
		}
		button := &blurredButton
		if b.focusIndex == len(b.changeUserPasswordInputs) {
			button = &focusedButton
		}
		builder.WriteString("\n" + *button)
		inputView = builder.String()
	}

	return bubbleStyle.Render(lipgloss.JoinVertical(lipgloss.Top, b.list.View(), addUserInputStyle.Render(inputView)))
}
