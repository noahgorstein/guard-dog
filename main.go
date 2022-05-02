package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/noahgorstein/stardog-go/internal/tui"
)

func main() {

	bubble := tui.New()
	p := tea.NewProgram(bubble, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
