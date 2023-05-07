package statusbar

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/noahgorstein/go-stardog/stardog"
	"github.com/noahgorstein/guard-dog/internal/mode"
)

var logo = fmt.Sprintf("%s%s%s",
	lipgloss.NewStyle().Foreground(nord12).Bold(true).Render("GUARD"),
	lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: string(nord10),
		Dark:  string(nord8),
	}).Bold(true).Render("ˁ˚ᴥ˚ˀ"),
	lipgloss.NewStyle().Foreground(nord12).Bold(true).Render("DOG"))

type Bubble struct {
	mode                  mode.ActiveMode
	stardogClient         stardog.Client
	width                 int
	help                  help.Model
	keys                  KeyMap
	Styles                Styles
	loggedInUser          string
	stardogServerEndpoint string

	StatusMessageLifetime time.Duration

	statusMessage      string
	statusMessageTimer *time.Timer
}

func New(stardogClient stardog.Client, loggedInUser, serverEndpoint string) Bubble {
	styles := DefaultStyles()
	help := help.New()
	help.Styles.ShortKey = styles.HelpKeyStyle
	help.Styles.ShortDesc = styles.HelpTextStyle
	help.Styles.ShortSeparator = styles.HelpTextStyle

	return Bubble{
		stardogClient:         stardogClient,
		Styles:                styles,
		help:                  help,
		keys:                  Keys,
		statusMessage:         lipgloss.NewStyle().Bold(true).Render(""),
		loggedInUser:          loggedInUser,
		stardogServerEndpoint: serverEndpoint,
	}
}

type statusMessageTimeoutMsg struct{}

func (b *Bubble) NewStatusMessage(s string, success bool) tea.Cmd {

	if success {
		b.statusMessage = b.Styles.SuccessMessageStyle.Render(s)
	} else {
		b.statusMessage = b.Styles.ErrorMessageStyle.Render(s)
	}

	if b.statusMessageTimer != nil {
		b.statusMessageTimer.Stop()
	}

	b.statusMessageTimer = time.NewTimer(b.StatusMessageLifetime)

	// Wait for timeout
	return func() tea.Msg {
		<-b.statusMessageTimer.C
		return statusMessageTimeoutMsg{}
	}
}

func (b Bubble) collectHelpBindings() []key.Binding {
	k := b.keys
	bindings := []key.Binding{}

	switch b.mode {
	case mode.UserListMode:
		bindings = append(bindings, k.UserDetailsView, k.Roles)
	case mode.UserDetailsMode:
		bindings = append(bindings, k.UserListView, k.Roles)
	case mode.RoleListMode:
		bindings = append(bindings, k.RoleDetailsView, k.Users)
	case mode.RoleDetailsMode:
		bindings = append(bindings, k.RoleListView, k.Users)
	}

	bindings = append(bindings, k.Help, k.Quit)
	return bindings
}

func (b Bubble) Init() tea.Cmd {
	return nil
}

func (b *Bubble) SetWidth(width int) {
	b.width = width - b.Styles.StatusBarStyle.GetHorizontalFrameSize() - 1
	b.Styles.StatusBarStyle.Width(b.width)
}

func (b *Bubble) hideStatusMessage() {
	b.statusMessage = ""
	if b.statusMessageTimer != nil {
		b.statusMessageTimer.Stop()
	}
}

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case statusMessageTimeoutMsg:
		b.hideStatusMessage()
	case tea.WindowSizeMsg:
		b.SetWidth(msg.Width)
	}
	return b, tea.Batch(cmd)
}

func (b *Bubble) UpdateMode(mode mode.ActiveMode) {
	b.mode = mode
}

func (b Bubble) View() string {

	help := lipgloss.NewStyle().
		PaddingLeft(1).
		Align(lipgloss.Left).
		Width(int(float64(b.width) * 0.5)).
		Render(logo + " | " + b.help.ShortHelpView(b.collectHelpBindings()))

	endpoint := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(int(float64(b.width) * 0.5)).
		Render(b.Styles.EndpointStyle.Render(
			fmt.Sprintf("%s@%s", b.loggedInUser, b.stardogServerEndpoint)))

	return b.Styles.StatusBarStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().PaddingLeft(1).Render(b.statusMessage),
			lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false, false).
				BorderForeground(grey).
				Render(lipgloss.JoinHorizontal(
					lipgloss.Top,
					help,
					endpoint))))
}
