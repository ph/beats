package module

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	f := func() {}

	t.Run("namespace and plugin doesn't exist", func(t *testing.T) {
		r := newRegistry()
		err := r.Register("processor", "foo", f)
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, 1, r.Size())
	})

	t.Run("namespace exists and plugin doesn't exist", func(t *testing.T) {
		r := newRegistry()
		r.Register("processor", "foo", f)
		err := r.Register("processor", "bar", f)
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, 2, r.Size())
	})

	t.Run("namespace exists and plugin exists", func(t *testing.T) {
		r := newRegistry()
		r.Register("processor", "foo", f)
		err := r.Register("processor", "foo", f)
		if !assert.Error(t, err) {
			return
		}
		assert.Equal(t, 1, r.Size())
	})
}

func TestModule(t *testing.T) {
	f := func() {}

	r := newRegistry()
	r.Register("processor", "foo", f)

	t.Run("when namespace and module are present", func(t *testing.T) {
		factory, err := r.Module("processor", "foo")
		if !assert.NotNil(t, factory) {
			return
		}
		assert.NoError(t, err)
	})

	t.Run("when namespace doesn't exist", func(t *testing.T) {
		_, err := r.Module("hello", "foo")
		if !assert.Error(t, err) {
			return
		}
	})
}

func TestModules(t *testing.T) {
	f := func() {}

	r := newRegistry()
	r.Register("processor", "foo", f)
	r.Register("processor", "foo2", f)

	t.Run("when namespace and module are present", func(t *testing.T) {
		modules, err := r.Modules("processor")
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, 2, len(modules))
	})

	t.Run("when namespace is not present", func(t *testing.T) {
		_, err := r.Modules("foobar")
		if !assert.Error(t, err) {
			return
		}
	})
}

func TestUnregister(t *testing.T) {
	f := func() {}

	t.Run("when the namespace and the module exist", func(t *testing.T) {
		r := newRegistry()
		r.Register("processor", "foo", f)
		assert.Equal(t, 1, r.Size())
		err := r.Unregister("processor", "foo")
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, 0, r.Size())
	})

	t.Run("when the namespace exist and the module doesn't", func(t *testing.T) {
		r := newRegistry()
		r.Register("processor", "foo", f)
		assert.Equal(t, 1, r.Size())
		err := r.Unregister("processor", "bar")
		if assert.Error(t, err) {
			return
		}
		assert.Equal(t, 0, r.Size())
	})

	t.Run("when the namespace doesn't exists", func(t *testing.T) {
		r := newRegistry()
		r.Register("processor", "foo", f)
		assert.Equal(t, 1, r.Size())
		err := r.Unregister("outputs", "bar")
		if assert.Error(t, err) {
			return
		}
		assert.Equal(t, 0, r.Size())
	})
}
