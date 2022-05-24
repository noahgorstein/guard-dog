package userlist

import (
	"strconv"

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

type sessionState int

const (
	idleState sessionState = iota
	createUserState
)

type Bubble struct {
	state      sessionState
	list       list.Model
	active     bool
	width      int
	height     int
	connection stardog.ConnectionDetails
	input      textinput.Model
	inputs     []textinput.Model
}

type item struct {
	title, desc string
}

func (i item) Title() string { return i.title }

func (i item) Description() string { return i.desc }

func (i item) FilterValue() string { return i.title }

func New(config config.Config) Bubble {
	connectionDetails := stardog.NewConnectionDetails(config.Endpoint, config.Username, config.Password)
	// items := []list.Item{}
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

	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.Placeholder = "Enter name of new user"
	input.CharLimit = 250
	input.Width = 20

	b := Bubble{
		state:      idleState,
		list:       userList,
		connection: *connectionDetails,
		input:      input,
	}

	var t textinput.Model
	for i := range b.inputs {
		t = textinput.New()
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Username"
			t.Focus()
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		b.inputs[i] = t
	}

	return b

}

func (b Bubble) Init() tea.Cmd {
	// return GetStardogUsersCmd(b.connection)
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
	b.list.SetSize(
		width-horizontal-vertical,
		height-vertical-lipgloss.Height(b.input.View())-inputStyle.GetVerticalPadding(),
	)
}

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
	case GetStardogUsersMsg:
		f.WriteString("get stardog users msg caught \n")
		f.WriteString("msg != nil -> " + strconv.FormatBool(msg != nil) + "\n")
		if msg != nil {
			cmd = b.list.SetItems(msg)
			cmds = append(cmds, cmd)
		}
	case tea.KeyMsg:

		if msg.Type == tea.KeyCtrlC {
			return b, tea.Quit
		}
		if msg.Type == tea.KeyCtrlA {
			if !b.input.Focused() {
				b.input.Focus()
				b.input.Placeholder = "Enter name of new user"
				b.state = createUserState
				return b, textinput.Blink
			}
		}

		if msg.Type == tea.KeyEnter {
			if b.state == createUserState {
				createStardogUserCmd := b.CreateStardogCmd(b.input.Value(), b.input.Value())
				getStardogUsersCmd := GetStardogUsersCmd(b.connection)
				cmds = append(cmds, tea.Sequentially(createStardogUserCmd, getStardogUsersCmd))
			}
			b.state = idleState
			b.input.Blur()
			b.input.Reset()
		}

		if msg.Type == tea.KeyCtrlD {
			if b.state == idleState {
				deleteStardogUserCmd := b.DeleteStardogUserCmd(b.GetCurrentUser())
				getStardogUsers := GetStardogUsersCmd(b.connection)
				cmds = append(cmds, tea.Sequentially(deleteStardogUserCmd, getStardogUsers))

			}
		}

	}

	if b.active {
		switch b.state {
		case idleState:
			f.WriteString("gonna update the list in idle state \n")
			b.list, cmd = b.list.Update(msg)
			cmds = append(cmds, cmd)
		case createUserState:
			f.WriteString("gonna update the input in createUser state \n")

			b.input, cmd = b.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	// b.list, cmd = b.list.Update(msg)
	// cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

var bubbleStyle = lipgloss.NewStyle().
	PaddingLeft(1).
	PaddingRight(1).
	BorderStyle(lipgloss.RoundedBorder())

var inputStyle = lipgloss.NewStyle().PaddingTop(1)

func (b Bubble) View() string {
	if b.active {
		bubbleStyle = bubbleStyle.BorderForeground(lipgloss.Color("33"))
	} else {
		bubbleStyle = bubbleStyle.BorderForeground(lipgloss.Color("#FAFAFA"))
	}

	var inputView string

	switch b.state {
	case idleState:
		inputView = ""
	case createUserState:
		inputView = b.input.View()
	default:
		inputView = ""
	}

	return bubbleStyle.Render(lipgloss.JoinVertical(lipgloss.Top, b.list.View(), inputStyle.Render(inputView)))
}
