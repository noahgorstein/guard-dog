package statusbar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/noahgorstein/stardog-go/internal/config"
	"github.com/noahgorstein/stardog-go/stardog"
)

type Bubble struct {
	connection stardog.ConnectionDetails
	width      int
	Style      lipgloss.Style
}

func (b Bubble) Init() tea.Cmd {
	return nil
}

func New(config config.Config) Bubble {
	connectionDetails := stardog.NewConnectionDetails(config.Endpoint, config.Username, config.Password)

	return Bubble{
		connection: *connectionDetails,
		Style:      lipgloss.NewStyle().Background(lipgloss.Color("240")),
	}
}

func (b *Bubble) SetWidth(width int) {
	b.width = width - b.Style.GetHorizontalFrameSize() - 1
	b.Style.Width(b.width)
}

func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.SetWidth(msg.Width)
	}
	return b, nil
}

var endpointStyle = lipgloss.NewStyle().Background(lipgloss.Color("168")).Padding(0, 1).Foreground(lipgloss.Color("#FAFAFA"))

func (b Bubble) View() string {
	return b.Style.Render(endpointStyle.Render(b.connection.Username + "@" + b.connection.Endpoint))
}
