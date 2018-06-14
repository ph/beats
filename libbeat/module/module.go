package module

import "fmt"

// Registry is the global plugin registry.
var Registry = newRegistry()

type Module struct {
	namespace string
	name      string
	factory   interface{}
}

func NewModule(namespace, name string, factory interface{}) Module {
	return Module{namespace: namespace, name: name, factory: factory}
}

// Bundle takes a list of of module and return them grouped by namespace
func Bundle(bundles ...interface{}) ([]Module, error) {
	var list []Module
	for _, bundle := range bundles {
		switch v := bundle.(type) {
		case []Module:
			list = append(list, v...)
		case Module:
			list = append(list, v)
		default:
			return nil, fmt.Errorf("incomptible type for bundle, expected 'Module' received '%T'", v)
		}
	}
	return list, nil
}

// MustBundle takes a list of modules panic if any errors is detected when bundling the modules.
func MustBundle(bundles ...interface{}) []Module {
	modules, err := Bundle(bundles)
	if err != nil {
		panic(err)
	}
	return modules
}
