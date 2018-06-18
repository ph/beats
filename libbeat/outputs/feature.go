package outputs

import (
	"errors"
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/cfgwarn"
	"github.com/elastic/beats/libbeat/feature"
)

// pluginKey is the namespace in the registry for all the outputs feature.
var pluginKey = "libbeat.outputs"

// Factory is used by output plugins to build an output instance
type Factory func(
	beat beat.Info,
	stats Observer,
	cfg *common.Config) (Group, error)

// Feature expose a new output module to the registry.
func Feature(name string, f Factory) feature.Module {
	return feature.NewModule(pluginKey, name, f)
}

// TODO: implements backward compatibility
func Plugin(name string, f Factory) map[string][]interface{} {
	return nil
}

// Group configures and combines multiple clients into load-balanced group of clients
// being managed by the publisher pipeline.
type Group struct {
	Clients   []Client
	BatchSize int
	Retry     int
}

// RegisterType registers a new output type.
func RegisterType(name string, f Factory) {
	err := feature.Registry.Register(pluginKey, name, f)
	if err != nil {
		panic(err)
	}
}

// Load creates and configures a output Group using a configuration object.
func Load(info beat.Info, stats Observer, name string, config *common.Config) (Group, error) {
	factory, err := feature.Registry.Module(pluginKey, name)
	if err != nil {
		return Group{}, fmt.Errorf("output type %v undefined", name)
	}

	b, ok := factory.(Factory)
	if !ok {
		return Group{}, errors.New("plugin does not match output plugin type")
	}

	if err := cfgwarn.CheckRemoved5xSetting(config, "flush_interval"); err != nil {
		return Fail(err)
	}

	if stats == nil {
		stats = NewNilObserver()
	}

	return b(info, stats, config)
}
