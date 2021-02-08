package gazebo

import (
	"github.com/johnfrankmorgan/gazebo/assert"
)

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
	env   *env
}

func NewVM() *VM {
	env := &env{values: map[string]*GObject{}}

	for name, builtin := range gbuiltins {
		env.define(name, builtin)
	}

	return &VM{
		stack: &stack{},
		env:   env,
	}
}

func (m *VM) Run(code Code) *GObject {
	var pc int

	for pc < len(code) {
		ins := code[pc]
		pc++

		switch ins.Opcode {
		case OpLoadConst:
			m.stack.push(NewGObjectInferred(ins.Arg))

		case OpStoreName:
			name := ins.Arg.(string)
			if m.env.defined(name) {
				m.env.assign(name, m.stack.pop())
			} else {
				m.env.define(name, m.stack.pop())
			}

		case OpLoadName:
			name := ins.Arg.(string)
			m.stack.push(m.env.lookup(name))

		case OpCallFunc:
			argc := ins.Arg.(int)
			args := make([]*GObject, argc)

			for i := 0; i < argc; i++ {
				args[argc-i-1] = m.stack.pop()
			}

			fun := m.stack.pop()
			ctx := &GFuncCtx{VM: m, Args: args}

			assert.True(fun.Type == gtypes.Func)

			m.stack.push(fun.Value.(GFunc)(ctx))

		case OpRelJump:
			pc += ins.Arg.(int)

		case OpRelJumpIfTrue:
			condition := m.stack.pop()
			if condition.IsTruthy() {
				pc += ins.Arg.(int)
			}

		case OpRelJumpIfFalse:
			condition := m.stack.pop()
			if !condition.IsTruthy() {
				pc += ins.Arg.(int)
			}

		default:
			assert.Unreached("unknown instruction: %v", ins)
		}
	}

	if m.stack.size() > 0 {
		return m.stack.pop()
	}

	return nil
}
