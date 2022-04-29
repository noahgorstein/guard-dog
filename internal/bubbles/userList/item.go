package userlist

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	title       string
	description string
	username    string
	superuser   bool
	enabled     bool
	roles       []string
}

var (
	superuser = lipgloss.NewStyle().Foreground(nord14).Render("superuser")
	roleStyle = lipgloss.NewStyle().Background(nord15).Foreground(nord6).Padding(0, 1)
)

func (i Item) Title() string {
	if i.superuser {
		return fmt.Sprintf("%s %s", i.title, superuser)
	}

	if !i.enabled {
		return lipgloss.NewStyle().Foreground(nord11).Render(fmt.Sprintf("%s (Disabled)", i.title))
	}
	return i.title
}

func (i Item) Description() string {
	if len(i.roles) > 0 {
		var sb strings.Builder

		for _, role := range i.roles {
			sb.WriteString(roleStyle.Bold(true).Render(role) + " ")
		}
		return sb.String()
	}
	return i.description
}

func (i Item) FilterValue() string {
	return i.title
}

func (i Item) Username() string {
	return i.username
}

func (i Item) Superuser() bool {
	return i.superuser
}

func (i Item) Enabled() bool {
	return i.enabled
}

func (i Item) Roles() []string {
	return i.roles
}
