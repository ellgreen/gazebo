package vm

import (
	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/compiler"
	"github.com/johnfrankmorgan/gazebo/compiler/op"
	"github.com/johnfrankmorgan/gazebo/g"
	"github.com/johnfrankmorgan/gazebo/g/modules"
)

// VM is the structure responsible for running code and keeping track of state
type VM struct {
	stack   *stack
	env     *env
	modules map[string]*modules.Module
}

// New creates a new VM
func New(argv ...string) *VM {
	env := new(env)

	for name, builtin := range g.Builtins() {
		env.define(name, builtin)
	}

	gargv := g.NewObjectList(nil)

	for _, arg := range argv {
		gargv.Append(g.NewObjectString(arg))
	}

	env.define("argv", gargv)

	return &VM{
		stack: new(stack),
		env:   env,
		modules: map[string]*modules.Module{
			"str":  modules.Str,
			"http": modules.HTTP,
		},
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
			case g.TypeInternalFunc:
				m.stack.push(g.Invoke(fun, args))

			case g.TypeFunc:
				fun := g.EnsureFunc(fun)

				vmenv := m.env
				env := &env{
					parent: fun.Env().(*env),
				}

				for i, param := range fun.Params() {
					env.define(param, args[i])
				}

				m.env = env
				m.stack.push(m.Run(fun.Code()))
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
			m.stack.push(g.NewObjectInternal(ins.Arg))

		case op.MakeFunc:
			body := m.stack.pop().Value().(compiler.Code)
			params := m.stack.pop().Value().([]string)
			m.stack.push(g.NewObjectFunc(params, body, m.env))

		case op.LoadModule:
			name := ins.Arg.(string)
			module, ok := m.modules[name]

			assert.True(ok, "undefined module: %s", name)

			module.Load(&m.env.values)

		case op.MakeList:
			length := ins.Arg.(int)
			values := make([]g.Object, length)

			for i := 0; i < length; i++ {
				values[length-i-1] = m.stack.pop()
			}

			m.stack.push(g.NewObjectList(values))

		case op.IndexGet:
			index := m.stack.pop()
			object := m.stack.pop()
			m.stack.push(object.Call(g.Protocols.Index, g.Args{index}))

		default:
			assert.Unreached("unknown instruction: 0x%02x (%s) %#v", int(ins.Opcode), ins.Opcode.Name(), ins)
		}
	}

	if m.stack.size() > 0 {
		return m.stack.pop()
	}

	return nil
}
