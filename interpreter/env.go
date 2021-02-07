package interpreter

import (
	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/gvalue"
)

type Env struct {
	parent *Env
	values map[string]gvalue.Instance
}

func NewEnv(values map[string]gvalue.Instance, parent *Env) *Env {
	if len(values) == 0 {
		values = make(map[string]gvalue.Instance)
	}

	return &Env{parent: parent, values: values}
}

func (m *Env) resolve(name string) *Env {
	if _, ok := m.values[name]; ok {
		return m
	}

	if m.parent != nil {
		return m.parent.resolve(name)
	}

	assert.Unreached("undefined variable: %q", name)
	return nil
}

func (m *Env) defined(name string) bool {
	if _, ok := m.values[name]; ok {
		return true
	}

	if m.parent != nil {
		return m.parent.defined(name)
	}

	return false
}

func (m *Env) Define(name string, value gvalue.Instance) {
	m.values[name] = value
}

func (m *Env) Assign(name string, value gvalue.Instance) {
	m.resolve(name).Define(name, value)
}

func (m *Env) Lookup(name string) gvalue.Instance {
	return m.resolve(name).values[name]
}
