package run

import (
	"github.com/elek/klepif/pkg"
	"github.com/elek/klepif/pkg/client"
	"github.com/elek/klepif/pkg/persistence"
	"github.com/elek/klepif/pkg/plugins"
	"github.com/elek/klepif/pkg/source"
	"github.com/spf13/viper"
	"strings"
)

func Run(prNumber int) error {

	ctx, ghSource, e := Init()
	if e != nil {
		return e
	}

	events, err := ghSource.GetEventsOfPr(prNumber);
	if err != nil {
		panic(err)
	}
	for _, plg := range plugins.Plugins.Instances {
		for _, event := range events {
			err = plg.HandlePrEvent(&ctx, event);
			if err != nil {
				return err;
			}
		}

	}

	return nil;
}

func Init() (plugins.ClientContext, source.GithubPr, error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("KLEPIF")
	viper.AutomaticEnv()

	err := pkg.ReadConfig()
	if err != nil {
		return plugins.ClientContext{}, source.GithubPr{}, err
	}
	github_config := pkg.GithubConfig{}
	err = viper.Sub("github").Unmarshal(&github_config)
	if err != nil {
		return plugins.ClientContext{}, source.GithubPr{}, err
	}
	if viper.GetString("github.token") != "" {
		github_config.Token = viper.GetString("github.token")
	}
	ctx := plugins.ClientContext{
		GithubClient: client.CreateGithubClient(&github_config),
		Persistence: &persistence.Dir{
			Path: viper.GetString("persistence.path"),
		},
		Contributors: viper.GetStringSlice("contributors.user"),
	}
	err = ctx.Persistence.Init()
	if err != nil {
		return plugins.ClientContext{}, source.GithubPr{}, err
	}
	ghSource := source.GithubPr{
		Client:      &ctx.GithubClient,
		Persistence: ctx.Persistence,
		Org:         github_config.Org,
		Repo:        github_config.Repo,
	}
	for key, _ := range viper.GetStringMap("actions") {
		pluginConfig := viper.Sub("actions." + key);
		err := plugins.Plugins.Initialize(key, pluginConfig, &ctx);
		if err != nil {
			return plugins.ClientContext{}, source.GithubPr{}, err
		}
	}
	return ctx, ghSource, nil
}
