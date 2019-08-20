package plugins

import (
	"github.com/elek/klepif/pkg/client"
	"github.com/elek/klepif/pkg/source"
	"github.com/spf13/viper"
)

type Label struct {
	DefaultPlugin
}

func (*Label) HandlePrEvent(ctx *ClientContext, change *source.GithubPrChange) error {
	if found, label := client.GetCommand("/label", change.Comments); found {
		err := ctx.GithubClient.AddLabel(*change.Pr.Number, label)
		if err != nil {
			return err
		}
	}
	return nil;
}

func init() {
	Plugins.RegisterPluginType("label", func(config *viper.Viper, context *ClientContext) (plugin Plugin, e error) {
		return &Label{}, nil
	});
}
