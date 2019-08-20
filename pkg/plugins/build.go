package plugins

import (
	"bytes"
	"fmt"
	"github.com/elek/klepif/pkg/client"
	"github.com/elek/klepif/pkg/persistence"
	"github.com/elek/klepif/pkg/source"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"html/template"
	"os/exec"
	"strings"
)

type BuildPlugin struct {
	DefaultPlugin
	GithubClient *client.GithubClient
	Persistence  persistence.Persistence
	Command      string
	Job          string
	DryRun       bool
	Label        string
	Rerun        bool
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (trigger *BuildPlugin) build(context *ClientContext, change *source.GithubPrChange) error {
	if !stringInSlice(change.Pr.User.GetLogin(), context.Contributors) {
		logrus.Warnf("PR '%s' is created by %s who is not on the contributor list. Build is skipped", change.Pr.GetTitle(), change.Pr.User.GetLogin());
		return nil
	}
	parameters := make(map[string]string)
	parameters["ref"] = change.Pr.GetHead().GetRef()
	parameters["repo"] = change.Pr.GetHead().Repo.GetName()
	parameters["org"] = change.Pr.GetHead().Repo.Owner.GetLogin()

	funcs := template.FuncMap{}
	funcs["tolower"] = strings.ToLower
	command, err := template.New("command").Funcs(funcs).Parse(trigger.Command)
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString("")
	err = command.Execute(buf, parameters)
	if err != nil {
		return err
	}
	if ! trigger.DryRun {
		cmdParts := strings.Split(buf.String(), " ");
		fmt.Printf("Executing command: %s\n", buf.String())
		out, err := exec.Command(cmdParts[0], cmdParts[1:]...).CombinedOutput()

		if err != nil {
			output := string(out[:])
			fmt.Println(output)
			return err
		}

		fmt.Println("Command Successfully Executed")
		output := string(out[:])
		fmt.Println(output)
	} else {
		fmt.Printf("Build command would be scheduled %s\n", buf.String())
	}
	return trigger.Persistence.Write(trigger.createPersistenceKey(change), *change.LastCommit.SHA)
}

func (trigger *BuildPlugin) createPersistenceKey(change *source.GithubPrChange) string {
	return fmt.Sprintf("github/%s/%s/PR-%d/last_checked_commit",
		trigger.GithubClient.Org,
		trigger.GithubClient.Repo,
		*change.Pr.Number);
}

func (trigger *BuildPlugin) HandlePrEvent(ctx *ClientContext, change *source.GithubPrChange) error {
	if len(trigger.Label) == 0 || HasLabel(change.Pr, trigger.Label) {
		lastCheckedCommit, err := trigger.Persistence.Read(trigger.createPersistenceKey(change))
		if err != nil {
			panic(err)
		}

		if trigger.Rerun || *change.LastCommit.SHA != lastCheckedCommit {
			return trigger.build(ctx, change)
		}

		if found, _ := client.GetCommand("/retest", change.Comments); found {
			return trigger.build(ctx, change)
		}
	}
	return nil
}

func init() {
	Plugins.RegisterPluginType("build", func(config *viper.Viper, ctx *ClientContext) (plugin Plugin, e error) {
		return &BuildPlugin{
			GithubClient: &ctx.GithubClient,
			Persistence:  ctx.Persistence,
			Job:          config.GetString("job"),
			DryRun:       config.GetBool("dryrun"),
			Label:        config.GetString("label"),
			Rerun:        config.GetBool("rerun"),
			Command:      config.GetString("command"),
		}, nil
	});
}
