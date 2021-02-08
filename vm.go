package gazebo

import "github.com/johnfrankmorgan/gazebo/assert"

type stack struct {
	values []*GObject
}

func (m *stack) push(value *GObject) {
	m.values = append(m.values, value)
}

func (m *stack) pop() *GObject {
	size := m.size()

	if size > 0 {
		value := m.values[size-1]
		m.values = m.values[:size-1]
		return value
	}

	assert.Unreached("stack empty")
	return nil
}

func (m *stack) size() int {
	return len(m.values)
}

type env struct {
	values map[string]*GObject
	parent *env
}

func (m *env) resolve(name string) *env {
	if _, ok := m.values[name]; ok {
		return m
	}

	if m.parent != nil {
		return m.parent.resolve(name)
	}

	return nil
}

func (m *env) lookup(name string) *GObject {
	env := m.resolve(name)
	if env != nil {
		return env.values[name]
	}

	assert.Unreached("undefined name %q", name)
	return nil
}

func (m *env) defined(name string) bool {
	return m.resolve(name) != nil
}

func (m *env) define(name string, value *GObject) {
	m.values[name] = value
}

func (m *env) assign(name string, value *GObject) {
	env := m.resolve(name)
	if env != nil {
		env.define(name, value)
		return
	}

	assert.Unreached("undefined name %q", name)
}

type VM struct {
	stack *stack
}
