package plugins

import (
	"fmt"
	"github.com/elek/klepif/pkg/source"
	"github.com/spf13/viper"
)

type StdOut struct {
	DefaultPlugin
}

func (*StdOut) HandlePrEvent(ctx *ClientContext, change *source.GithubPrChange) error {
	fmt.Printf("%d new commit is detected\n", len(change.Commits))
	fmt.Println("")
	for _, commit := range change.Commits {
		fmt.Printf("[PR %d] %s %s\n", *change.Pr.Number, *commit.SHA, *commit.Commit.Message);
	}
	fmt.Println()
	for _, comment := range change.Comments {
		fmt.Printf("[PR %d] %s %s\n", *change.Pr.Number, *comment.AuthorAssociation, *comment.Body);
	}
	return nil;
}

func init() {
	Plugins.RegisterPluginType("print", func(config *viper.Viper, context *ClientContext) (plugin Plugin, e error) {
		return &StdOut{}, nil
	});
}