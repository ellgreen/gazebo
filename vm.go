package gazebo

import (
	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/compiler"
	"github.com/johnfrankmorgan/gazebo/compiler/op"
	"github.com/johnfrankmorgan/gazebo/g"
)

type stack struct {
	values []g.Object
}

func (m *stack) push(value g.Object) {
	m.values = append(m.values, value)
}

func (m *stack) pop() g.Object {
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
	values map[string]g.Object
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

func (m *env) lookup(name string) g.Object {
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

func (m *env) define(name string, value g.Object) {
	m.values[name] = value
}

func (m *env) assign(name string, value g.Object) {
	env := m.resolve(name)
	if env != nil {
		env.define(name, value)
		return
	}

	assert.Unreached("undefined name %q", name)
}

// VM is the structure responsible for running code and keeping track of state
type VM struct {
	stack *stack
	env   *env
}

// NewVM creates a new VM
func NewVM() *VM {
	env := &env{values: map[string]g.Object{}}

	for name, builtin := range g.Builtins() {
		env.define(name, builtin)
	}

	return &VM{
		stack: &stack{},
		env:   env,
	}
}

// Run runs the provided code
func (m *VM) Run(code compiler.Code) g.Object {
	var pc int

	for pc < len(code) {
		ins := code[pc]
		pc++

		switch ins.Opcode {
		case op.LoadConst:
			m.stack.push(g.NewObject(ins.Arg))

		case op.StoreName:
			name := ins.Arg.(string)
			if m.env.defined(name) {
				m.env.assign(name, m.stack.pop())
			} else {
				m.env.define(name, m.stack.pop())
			}

		case op.LoadName:
			name := ins.Arg.(string)
			m.stack.push(m.env.lookup(name))

		case op.CallFunc:
			argc := ins.Arg.(int)
			args := make(g.Args, argc)

			for i := 0; i < argc; i++ {
				args[argc-i-1] = m.stack.pop()
			}

			fun := m.stack.pop()

			switch fun.Type() {
			case g.TypeBuiltinFunc:
				m.stack.push(fun.Call(g.Protocols.Invoke, args))

			case g.TypeFunc:
				desc := fun.Value().(g.FuncDescription)
				vmenv := m.env
				env := &env{
					values: map[string]g.Object{},
					parent: desc.Env.(*env),
				}

				for i, param := range desc.Params {
					env.define(param, args[i])
				}

				m.env = env
				m.stack.push(m.Run(desc.Body))
				m.env = vmenv

			default:
				assert.Unreached("unexpected type called as function: gtypes.%s", fun.Type().Name)
			}

		case op.RelJump:
			pc += ins.Arg.(int)

		case op.RelJumpIfTrue:
			condition := m.stack.pop()
			if g.IsTruthy(condition) {
				pc += ins.Arg.(int)
			}

		case op.RelJumpIfFalse:
			condition := m.stack.pop()
			if !g.IsTruthy(condition) {
				pc += ins.Arg.(int)
			}

		case op.PushValue:
			m.stack.push(g.NewInternalObject(ins.Arg))

		case op.MakeFunc:
			body := m.stack.pop().Value().(compiler.Code)
			params := m.stack.pop().Value().([]string)

			m.stack.push(g.NewObject(g.FuncDescription{
				Params: params,
				Body:   body,
				Env:    m.env,
			}))

		default:
			assert.Unreached("unknown instruction: %v", ins)
		}
	}

	if m.stack.size() > 0 {
		return m.stack.pop()
	}

	return nil
}
