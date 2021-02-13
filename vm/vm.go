package vm

import (
	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/compiler"
	"github.com/johnfrankmorgan/gazebo/compiler/op"
	"github.com/johnfrankmorgan/gazebo/g"
)

// VM is the structure responsible for running code and keeping track of state
type VM struct {
	stack *stack
	env   *env
}

// New creates a new VM
func New() *VM {
	env := new(env)

	for name, builtin := range g.Builtins() {
		env.define(name, builtin)
	}

	return &VM{
		stack: new(stack),
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