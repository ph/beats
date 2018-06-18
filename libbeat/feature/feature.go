package feature

// Registry is the global plugin registry.
var Registry = newRegistry()

// Module
// TODO: Make it an interface and use it in the registry.
type Module struct {
	namespace string
	name      string
	factory   interface{}
}

// NewModule creates a new module.
func NewModule(namespace, name string, factory interface{}) Module {
	return Module{namespace: namespace, name: name, factory: factory}
}

// Bundle merge multiple bundle declaration into one.
func Bundle(bundles ...Module) ([]Module, error) {
	var list []Module
	for _, bundle := range bundles {
		list = append(list, bundle)
	}
	return list, nil
}

// MustBundle takes a list of modules panic if any errors is detected when bundling the modules.
func MustBundle(bundles ...Module) []Module {
	modules, err := Bundle(bundles...)
	if err != nil {
		panic(err)
	}
	return modules
}

// Register registers a list of modules
func Register(bundles []Module) error {
	for _, module := range bundles {
		Registry.Register(module.namespace, module.name, module.factory)
	}
	return nil
}
