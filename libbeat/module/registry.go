package module

import (
	"fmt"
	"sync"

	"github.com/elastic/beats/libbeat/logp"
)

type mapper map[string]map[string]interface{}

// Registry contains reference to all the libbeat module, a module is bitsize piece of code that
// only do one thing. This is different from a plugin that can be composed of multiple modules.
//
// Example of modules:
// - Processor
// - Queue
// - Ouput
//
// Specific modules can be implemented by a specific beat.
type Registry struct {
	sync.RWMutex
	namespaces mapper
	log        *logp.Logger
}

// NewRegistry returns a new registry.
func NewRegistry() *Registry {
	return &Registry{
		namespaces: make(mapper),
		log:        logp.NewLogger("registry"),
	}
}

// Register registers a new module into a specific namespace, it will lazy creates the namespace,
// module name must be unique in the registry.
func (r *Registry) Register(namespace, moduleName string, factory interface{}) error {
	r.Lock()
	defer r.Unlock()

	// Lazy create namespaces
	_, found := r.namespaces[namespace]
	if !found {
		r.namespaces[namespace] = make(map[string]interface{})
	}

	_, found = r.namespaces[namespace][moduleName]
	if found {
		return fmt.Errorf(
			"could not register new module '%s' in namespace '%s', module name must be unique",
			moduleName,
			namespace,
		)
	}

	r.log.Debugw("registering new module", "namespace", namespace, "module", moduleName)

	r.namespaces[namespace][moduleName] = factory

	return nil
}

// Unregister(namespace, module)

// Module returns a specific module from a specific namespace.
func (r *Registry) Module(namespace, moduleName string) (interface{}, error) {
	v, found := r.namespaces[namespace]
	if !found {
		return nil, fmt.Errorf("unknown namespace named '%s'", namespace)
	}

	m, found := v[moduleName]
	if !found {
		return nil, fmt.Errorf("unknown module '%s' in namespace '%s'", moduleName, namespace)
	}

	return m, nil
}

// Size returns the number of registered plugins in the registry.
func (r *Registry) Size() int {
	r.RLock()
	defer r.RUnlock()
	c := 0
	for _, namespace := range r.namespaces {
		c += len(namespace)
	}

	return c
}

// Fix the registry
// Add the new method on the output package to use the new registry.
// Wrap the factory of the output the new registry
// Make sure we have some kind of compaibility layer
// Make next PR to move the modules to new registry
// Allow to inject DI to beats instance (backward compaibility?)
