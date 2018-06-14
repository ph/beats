package outputs

import "github.com/elastic/beats/libbeat/module"

// Module expose a new output module to the registry.
func Module(name string, f Factory) module.Module {
	return module.NewModule(pluginKey, name, f)
}
