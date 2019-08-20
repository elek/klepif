package plugins

import (
	"github.com/elek/klepif/pkg/client"
	"github.com/elek/klepif/pkg/persistence"
)

type ClientContext struct {
	GithubClient client.GithubClient
	Persistence  persistence.Persistence
	Contributors []string
}
