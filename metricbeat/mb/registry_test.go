// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// +build !integration

package mb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/feature"
)

const (
	moduleName    = "mymodule"
	metricSetName = "mymetricset"
)

var fakeModuleFactory = func(b BaseModule) (Module, error) { return nil, nil }
var fakeModuleFactory2 = func(b BaseModule) (Module, error) { return nil, nil }
var fakeMetricSetFactory = func(b BaseMetricSet) (MetricSet, error) { return nil, nil }

func TestAddModuleEmptyName(t *testing.T) {
	registry := NewRegister(feature.NewRegistry())
	err := registry.AddModule("", fakeModuleFactory)
	assert.Error(t, err)
}

func TestAddModuleDuplicateName(t *testing.T) {
	registry := NewRegister(feature.NewRegistry())
	err := registry.AddModule(moduleName, fakeModuleFactory)
	if err != nil {
		t.Fatal(err)
	}

	err = registry.AddModule(moduleName, fakeModuleFactory2)
	assert.Error(t, err)
}

func TestAddMetricSetEmptyModuleName(t *testing.T) {
	registry := NewRegister(feature.NewRegistry())
	err := registry.AddMetricSet("", metricSetName, fakeMetricSetFactory)
	if assert.Error(t, err) {
		assert.Equal(t, "module name is required", err.Error())
	}
}

func TestAddMetricSetEmptyName(t *testing.T) {
	registry := NewRegister(feature.NewRegistry())
	err := registry.AddMetricSet(moduleName, "", fakeMetricSetFactory)
	assert.Error(t, err)
}

func TestAddMetricSetNilFactory(t *testing.T) {
	registry := NewRegister(feature.NewRegistry())
	err := registry.AddMetricSet(moduleName, metricSetName, nil)
	if assert.Error(t, err) {
		assert.Equal(t, "metricset 'mymodule/mymetricset' cannot be registered with a nil factory", err.Error())
	}
}

func TestAddMetricSetDuplicateName(t *testing.T) {
	registry := NewRegister(feature.NewRegistry())
	err := registry.AddMetricSet(moduleName, metricSetName, fakeMetricSetFactory)
	if err != nil {
		t.Fatal(err)
	}

	err = registry.AddMetricSet(moduleName, metricSetName, fakeMetricSetFactory)
	assert.Error(t, err)
}

func TestAddMetricSet(t *testing.T) {
	registry := NewRegister(feature.NewRegistry())
	err := registry.AddMetricSet(moduleName, metricSetName, fakeMetricSetFactory)
	if err != nil {
		t.Fatal(err)
	}
	l := registry.MetricSets(moduleName)
	assert.Equal(t, metricSetName, l[0])
}

func TestModuleFactory(t *testing.T) {
	registry := NewRegister(feature.NewRegistry())
	registry.AddModule(moduleName, fakeModuleFactory)

	module := registry.moduleFactory(moduleName)
	assert.NotNil(t, module)
}

func TestModuleFactoryUnknownModule(t *testing.T) {
	registry := NewRegister(feature.NewRegistry())
	module := registry.moduleFactory("unknown")
	assert.Nil(t, module)
}

func TestMetricSetFactory(t *testing.T) {
	t.Run("without HostParser", func(t *testing.T) {
		registry := NewRegister(feature.NewRegistry())
		err := registry.AddMetricSet(moduleName, metricSetName, fakeMetricSetFactory)
		if err != nil {
			t.Fatal(err)
		}

		reg, err := registry.metricSetRegistration(moduleName, metricSetName)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, metricSetName, reg.Name)
		assert.NotNil(t, reg.Factory)
		assert.Nil(t, reg.HostParser)
		assert.False(t, reg.IsDefault)
		assert.Empty(t, reg.Namespace)
	})

	t.Run("with HostParser", func(t *testing.T) {
		registry := NewRegister(feature.NewRegistry())
		hostParser := func(Module, string) (HostData, error) { return HostData{}, nil }
		err := registry.AddMetricSet(moduleName, metricSetName, fakeMetricSetFactory, hostParser)
		if err != nil {
			t.Fatal(err)
		}

		reg, err := registry.metricSetRegistration(moduleName, metricSetName)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, reg.HostParser) // Can't compare functions in Go so just check for non-nil.
	})

	t.Run("with options HostParser", func(t *testing.T) {
		registry := NewRegister(feature.NewRegistry())
		hostParser := func(Module, string) (HostData, error) { return HostData{}, nil }
		err := registry.addMetricSet(moduleName, metricSetName, fakeMetricSetFactory, WithHostParser(hostParser))
		if err != nil {
			t.Fatal(err)
		}

		reg, err := registry.metricSetRegistration(moduleName, metricSetName)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, reg.HostParser) // Can't compare functions in Go so just check for non-nil.
	})

	t.Run("with namespace", func(t *testing.T) {
		const ns = moduleName + "foo.bar"

		registry := NewRegister(feature.NewRegistry())
		err := registry.addMetricSet(moduleName, metricSetName, fakeMetricSetFactory, WithNamespace(ns))
		if err != nil {
			t.Fatal(err)
		}

		reg, err := registry.metricSetRegistration(moduleName, metricSetName)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, metricSetName, reg.Name)
		assert.NotNil(t, reg.Factory)
		assert.Nil(t, reg.HostParser)
		assert.False(t, reg.IsDefault)
		assert.Equal(t, ns, reg.Namespace)
	})
}

func TestDefaultMetricSet(t *testing.T) {
	registry := NewRegister(feature.NewRegistry())
	err := registry.addMetricSet(moduleName, metricSetName, fakeMetricSetFactory, DefaultMetricSet())
	if err != nil {
		t.Fatal(err)
	}

	names, err := registry.DefaultMetricSets(moduleName)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, names, metricSetName)
}

func TestMetricSetQuery(t *testing.T) {
	registry := NewRegister(feature.NewRegistry())
	err := registry.AddMetricSet(moduleName, metricSetName, fakeMetricSetFactory)
	if err != nil {
		t.Fatal(err)
	}

	metricsets := registry.MetricSets(moduleName)
	assert.Equal(t, len(metricsets), 1)
	assert.Equal(t, metricsets[0], metricSetName)

	metricsets = registry.MetricSets("foo")
	assert.Equal(t, len(metricsets), 0)
}

func TestModuleQuery(t *testing.T) {
	registry := NewRegister(feature.NewRegistry())
	registry.AddModule(moduleName, fakeModuleFactory)

	modules := registry.Modules()
	assert.Equal(t, len(modules), 1)
	assert.Equal(t, modules[0], moduleName)
}
