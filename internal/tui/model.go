package tui

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	help "github.com/knipferrc/teacup/help"
	"github.com/noahgorstein/go-stardog/stardog"
	roledetails "github.com/noahgorstein/guard-dog/internal/bubbles/roleDetails"
	rolelist "github.com/noahgorstein/guard-dog/internal/bubbles/roleList"
	"github.com/noahgorstein/guard-dog/internal/bubbles/statusbar"
	userdetails "github.com/noahgorstein/guard-dog/internal/bubbles/userDetails"
	userlist "github.com/noahgorstein/guard-dog/internal/bubbles/userList"
	"github.com/noahgorstein/guard-dog/internal/mode"
)

type Bubble struct {
	userList    userlist.Bubble
	roleList    rolelist.Bubble
	userDetails userdetails.Bubble
	roleDetails roledetails.Bubble
	help        help.Bubble
	helpEnabled bool
	statusBar   statusbar.Bubble
	activeView  mode.ActiveMode
	width       int
	height      int
}

func getUserListViewHelpEntries() []help.Entry {
	return []help.Entry{
		{Key: "↑/↓/←/→", Description: "Navigate list"},
		{Key: "g", Description: "Jump to top of list"},
		{Key: "G", Description: "Jump to end of list"},
		{Key: "/", Description: "Filter user list"},
		{Key: "esc", Description: "Reset input field/clear filter"},
		{Key: "", Description: ""},
		{Key: "ctrl+n", Description: "Create new user"},
		{Key: "ctrl+e", Description: "Enable/disable user"},
		{Key: "ctrl+p", Description: "Change selected user's password"},
		{Key: "ctrl+d", Description: "Delete selected user"},
		{Key: "", Description: ""},
		{Key: "tab", Description: "Toggle active view"},
		{Key: "?", Description: "Close this help menu"},
		{Key: "ctrl+c", Description: "Exit guard-dog"},
	}
}

func getUserDetailsViewHelpEntries() []help.Entry {
	return []help.Entry{
		{Key: "T", Description: "Table mode"},
		{Key: "↑/↓/←/→", Description: "Scroll viewport/navigate table"},
		{Key: "/", Description: "Filter table"},
		{Key: "esc", Description: "Reset input field/clear filter/back"},
		{Key: "", Description: ""},
		{Key: "ctrl+g", Description: "Grant user permission"},
		{Key: "ctrl+r", Description: "Assign role to user"},
		{Key: "ctrl+u", Description: "Unassign role from user"},
		{Key: "", Description: ""},
		{Key: "tab", Description: "Toggle active view"},
		{Key: "?", Description: "Close this help menu"},
		{Key: "ctrl+c", Description: "Exit guard-dog"},
	}
}

func getRoleListViewHelpEntries() []help.Entry {
	return []help.Entry{
		{Key: "↑/↓/←/→", Description: "Navigate list"},
		{Key: "g", Description: "Jump to top of list"},
		{Key: "G", Description: "Jump to end of list"},
		{Key: "/", Description: "filter role list"},
		{Key: "esc", Description: "Reset input field/clear filter/back"},
		{Key: "", Description: ""},
		{Key: "ctrl+n", Description: "Create new role"},
		{Key: "ctrl+d", Description: "Delete selected role"},
		{Key: "ctrl+f", Description: "Forcefully delete selected role"},
		{Key: "", Description: ""},
		{Key: "tab", Description: "Toggle active view"},
		{Key: "?", Description: "Close this help menu"},
		{Key: "ctrl+c", Description: "Exit guard-dog"},
	}
}

func getRoleDetailsViewHelpEntries() []help.Entry {
	return []help.Entry{
		{Key: "T", Description: "Table mode"},
		{Key: "↑/↓/←/→", Description: "Scroll viewport/navigate table"},
		{Key: "/", Description: "Filter permission table"},
		{Key: "esc", Description: "Reset input field/clear filter"},
		{Key: "ctrl+g", Description: "Grant role permission"},
		{Key: "", Description: ""},
		{Key: "tab", Description: "Toggle active view"},
		{Key: "?", Description: "Close this help menu"},
		{Key: "ctrl+c", Description: "Exit guard-dog"},
	}
}

func New(stardogClient stardog.Client, loggedInUser string) Bubble {

	userlist := userlist.New(stardogClient)
	userlist.SetIsActive(true)
	userlist.SetListTitleBackgroundColor(lipgloss.AdaptiveColor{
		Light: string(nord12),
		Dark:  string(nord12),
	})
	userlist.SetListTitleForegroundColor(lipgloss.AdaptiveColor{
		Light: string(nord6),
		Dark:  string(nord6),
	})
	userlist.Styles.InactiveListStyle.BorderForeground(lipgloss.AdaptiveColor{
		Light: string(black),
		Dark:  string(grey),
	})
	userlist.Styles.ActiveListStyle.BorderForeground(lipgloss.Color(nord12))

	rolelist := rolelist.New(stardogClient)
	rolelist.SetIsActive(true)
	rolelist.SetListTitleBackgroundColor(lipgloss.AdaptiveColor{
		Light: string(nord12),
		Dark:  string(nord12),
	})
	rolelist.SetListTitleForegroundColor(lipgloss.AdaptiveColor{
		Light: string(nord6),
		Dark:  string(nord6),
	})
	rolelist.Styles.InactiveListStyle.BorderForeground(lipgloss.AdaptiveColor{
		Light: string(black),
		Dark:  string(grey),
	})
	rolelist.Styles.ActiveListStyle.BorderForeground(lipgloss.Color(nord12))

	roledetails := roledetails.New(stardogClient)
	roledetails.SetIsActive(false)
	roledetails.Styles.ActiveViewportStyle.BorderForeground(lipgloss.Color(nord12))
	roledetails.Styles.InactiveViewportStyle.BorderForeground(lipgloss.AdaptiveColor{
		Light: string(black),
		Dark:  string(grey),
	})

	userdetails := userdetails.New(stardogClient, loggedInUser)
	userdetails.SetIsActive(false)
	userdetails.Styles.ActiveViewportStyle.BorderForeground(lipgloss.Color(nord12))
	userdetails.Styles.InactiveViewportStyle.BorderForeground(lipgloss.AdaptiveColor{
		Light: string(black),
		Dark:  string(grey),
	})

	helpModel := help.New(
		true,
		false,
		"Help",
		help.TitleColor{
			Background: lipgloss.AdaptiveColor{Light: string(nord15), Dark: string(nord15)},
			Foreground: lipgloss.AdaptiveColor{Light: string(nord6), Dark: string(nord6)},
		},
		lipgloss.AdaptiveColor{Light: "16", Dark: "15"},
		getUserListViewHelpEntries(),
	)

	helpModel.Borderless = true

	statusbar := statusbar.New(stardogClient)
	statusbar.StatusMessageLifetime = time.Duration(30 * time.Second)

	return Bubble{
		userList:    userlist,
		roleList:    rolelist,
		userDetails: userdetails,
		roleDetails: roledetails,
		help:        helpModel,
		helpEnabled: false,
		activeView:  mode.UserListMode,
		statusBar:   statusbar,
	}
}
