package module

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	f := func() {}

	t.Run("namespace and plugin doesn't exist", func(t *testing.T) {
		r := NewRegistry()
		err := r.Register("processor", "foo", f)
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, 1, r.Size())
	})

	t.Run("namespace exists and plugin doesn't exist", func(t *testing.T) {
		r := NewRegistry()
		r.Register("processor", "foo", f)
		err := r.Register("processor", "bar", f)
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, 2, r.Size())
	})

	t.Run("namespace exists and plugin exists", func(t *testing.T) {
		r := NewRegistry()
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

	r := NewRegistry()
	r.Register("processor", "foo", f)

	t.Run("when namespace and module are present", func(t *testing.T) {
		factory, err := r.Module("processor", "foo")
		if !assert.NotNil(t, factory) {
			return
		}
		assert.NoError(t, err)
	})
}
