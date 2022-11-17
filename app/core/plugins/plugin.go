package plugins

import (
	"fmt"

	"github.com/space-fold-technologies/aurora-service/app/core/logging"
)

type PluginCategory string

const (
	REVERSE_PROXY PluginCategory = "reverse-proxy-plugin"
)

type PluginProviderCallback func(p Plugin) error
type Plugin interface {
	//Labels(target map[string]string, network, hostname string, ports []uint) error
	Name() string
	Category() PluginCategory
	OnStartUp() error
	OnShutDown() error
	Call(operation string, request interface{}, response interface{}) error
}

type PluginRegistry interface {
	Put(category PluginCategory, plugin Plugin) error
	Invoke(category PluginCategory, loader PluginProviderCallback) error
	StartUp()
	Shutdown()
}

func NewPluginRegistry() PluginRegistry {
	instance := new(internalPluginRegistry)
	instance.plugins = make(map[PluginCategory]Plugin)
	return instance
}

type internalPluginRegistry struct {
	plugins map[PluginCategory]Plugin
}

func (ir *internalPluginRegistry) Put(category PluginCategory, plugin Plugin) error {
	if _, found := ir.plugins[category]; found {
		return fmt.Errorf("plugin of type: [%s] already registered", category)
	}
	ir.plugins[category] = plugin
	return nil
}

func (ir *internalPluginRegistry) Invoke(category PluginCategory, loader PluginProviderCallback) error {
	if plugin, found := ir.plugins[category]; !found {
		return fmt.Errorf("plugin of category : [%s] not found", category)
	} else {
		return loader(plugin)
	}
}

func (ir *internalPluginRegistry) StartUp() {
	for category, plugin := range ir.plugins {
		logging.GetInstance().Infof("LOADING :[%s] plugin", category)
		if err := plugin.OnStartUp(); err != nil {
			logging.GetInstance().Error(err)
		}
	}
}

func (ir *internalPluginRegistry) Shutdown() {
	for category, plugin := range ir.plugins {
		logging.GetInstance().Infof("SHUTTING DOWN :[%s] plugin", category)
		if err := plugin.OnStartUp(); err != nil {
			logging.GetInstance().Error(err)
		}
	}
}
