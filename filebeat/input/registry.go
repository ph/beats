package input

import (
	"fmt"

	"github.com/elastic/beats/filebeat/channel"
	"github.com/elastic/beats/filebeat/input/file"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/feature"
)

type Context struct {
	States        []file.State
	Done          chan struct{}
	BeatDone      chan struct{}
	DynamicFields *common.MapStrPointer
	Meta          map[string]string
}

// Factory is used to register functions creating new Input instances.
type Factory = func(config *common.Config, connector channel.Connector, context Context) (Input, error)

func Register(name string, factory Factory) error {
	return feature.Registry.Register(pluginKey, name, factory)
}

// Feature expose a new input module to the registry.
func Feature(name string, f Factory) feature.Module {
	return feature.NewModule(pluginKey, name, f)
}

func GetFactory(name string) (Factory, error) {
	f, err := feature.Registry.Module(pluginKey, name)
	if err != nil {
		return nil, err
	}

	factory, ok := f.(Factory)
	if !ok {
		return nil, fmt.Errorf("plugin '%s' does not match the input type", name)
	}

	return factory, nil
}
