package tuiv2

import (
	"github.com/noahgorstein/stardog-go/internal/bubbles/statusbar"
	userdetails "github.com/noahgorstein/stardog-go/internal/bubbles/userDetails"
	userlist "github.com/noahgorstein/stardog-go/internal/bubbles/userList"
	"github.com/noahgorstein/stardog-go/internal/config"
	"github.com/noahgorstein/stardog-go/stardog"
)

type activeView int

const (
	listView activeView = iota
	detailsView
)

type Bubble struct {
	userList          userlist.Bubble
	userDetails       userdetails.Bubble
	statusBar         statusbar.Bubble
	stardogConnection stardog.ConnectionDetails
	activeView        activeView
}

func New(config config.Config) Bubble {

	stardogConnection := stardog.NewConnectionDetails(config.Endpoint, config.Username, config.Password)

	userlist := userlist.New(config)
	userlist.SetIsActive(true)

	userdetails := userdetails.New(config)
	userdetails.SetIsActive(false)

	statusbar := statusbar.New(config)

	return Bubble{
		userList:          userlist,
		userDetails:       userdetails,
		stardogConnection: *stardogConnection,
		activeView:        listView,
		statusBar:         statusbar,
	}
}
