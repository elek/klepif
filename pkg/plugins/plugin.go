package plugins

import (
	"errors"
	"github.com/elek/klepif/pkg/source"
	"github.com/spf13/viper"
)

type Plugin interface {
	HandlePrEvent(ctx *ClientContext, client *source.GithubPrChange) error
}

type DefaultPlugin struct {
}

type PluginFactory = func(config *viper.Viper, context *ClientContext) (Plugin, error)

type PluginRegistry struct {
	PluginTypes map[string]PluginFactory
	Instances   []Plugin
}

func (registry *PluginRegistry) RegisterPluginType(name string, factory PluginFactory) {
	registry.PluginTypes[name] = factory
}

func (registry *PluginRegistry) Initialize(pluginType string, pluginConfig *viper.Viper, clientContext *ClientContext) error {
	if factory, ok := registry.PluginTypes[pluginType]; ok {
		pluginInstance, err := factory(pluginConfig, clientContext)
		if err != nil {
			return err
		}
		registry.Instances = append(registry.Instances, pluginInstance)
		return nil
	} else {
		return errors.New(pluginType + " plugin is not registered")
	}
}

var Plugins = &PluginRegistry{
	PluginTypes: make(map[string]PluginFactory),
	Instances:   make([]Plugin, 0),
}