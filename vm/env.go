package vm

import (
	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/g"
)

type env struct {
	parent *env
	values g.Attributes
}

func (m *env) resolve(name string) *env {
	if m.values.Has(name) {
		return m
	}

	if m.parent != nil {
		return m.parent.resolve(name)
	}

	return nil
}

func (m *env) lookup(name string) g.Object {
	if env := m.resolve(name); env != nil {
		return env.values.Get(name)
	}

	assert.Unreached("undefined name: %s", name)
	return nil
}

func (m *env) defined(name string) bool {
	return m.resolve(name) != nil
}

func (m *env) define(name string, value g.Object) {
	m.values.Set(name, value)
}

func (m *env) assign(name string, value g.Object) {
	if env := m.resolve(name); env != nil {
		env.values.Set(name, value)
		return
	}

	assert.Unreached("undefined name: %s", name)
	return
}
