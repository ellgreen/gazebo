package gazebo

import (
	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/compiler"
	"github.com/johnfrankmorgan/gazebo/compiler/op"
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

// VM is the structure responsible for running code and keeping track of state
type VM struct {
	stack *stack
	env   *env
}

// NewVM creates a new VM
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

// Run runs the provided code
func (m *VM) Run(code compiler.Code) *GObject {
	var pc int

	for pc < len(code) {
		ins := code[pc]
		pc++

		switch ins.Opcode {
		case op.LoadConst:
			m.stack.push(NewGObjectInferred(ins.Arg))

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
			args := make([]*GObject, argc)

			for i := 0; i < argc; i++ {
				args[argc-i-1] = m.stack.pop()
			}

			fun := m.stack.pop()

			switch fun.Type {
			case gtypes.Func:
				ctx := &GFuncCtx{VM: m, Args: args}
				m.stack.push(fun.Interface().(GFunc)(ctx))

			case gtypes.UserFunc:
				guserfunc := fun.Value.(GUserFunc)
				vmenv := m.env
				env := &env{values: map[string]*GObject{}, parent: guserfunc.env}
				for i, param := range guserfunc.params {
					env.values[param] = args[i]
				}
				m.env = env
				m.stack.push(m.Run(guserfunc.body))
				m.env = vmenv

			default:
				assert.Unreached("unexpected type called as function: gtypes.%s", fun.Type.Name)
			}

		case op.RelJump:
			pc += ins.Arg.(int)

		case op.RelJumpIfTrue:
			condition := m.stack.pop()
			if condition.IsTruthy() {
				pc += ins.Arg.(int)
			}

		case op.RelJumpIfFalse:
			condition := m.stack.pop()
			if !condition.IsTruthy() {
				pc += ins.Arg.(int)
			}

		case op.PushValue:
			m.stack.push(&GObject{Type: gtypes.Internal, Value: ins.Arg})

		case op.MakeFunc:
			body := m.stack.pop().Interface().(compiler.Code)
			params := m.stack.pop().Interface().([]string)
			m.stack.push(&GObject{
				Type: gtypes.UserFunc,
				Value: GUserFunc{
					params: params,
					body:   body,
					env:    m.env,
				},
			})

		default:
			assert.Unreached("unknown instruction: %v", ins)
		}
	}

	if m.stack.size() > 0 {
		return m.stack.pop()
	}

	return nil
}
