package interpreter

import (
	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/gvalue"
	"github.com/johnfrankmorgan/gazebo/op"
)

type Interpreter struct{}

var _the *Interpreter

func Create() {
	_the = &Interpreter{}
}

func The() *Interpreter {
	assert.NotNil(_the)

	return _the
}

func (m *Interpreter) Eval(code []op.Instruction) gvalue.Instance {
	var (
		pc    int
		stack Stack
	)

	env := NewEnv(nil, nil)
	LoadBuiltins(env)

	for pc < len(code) {
		ins := code[pc]
		pc++

		switch ins.Opcode {
		case op.LoadConst:
			stack.Push(ins.Arg)

		case op.StoreName:
			name := ins.Arg.ToString()
			f := env.Define

			if env.defined(name) {
				f = env.Assign
			}

			f(name, stack.Pop())

		case op.LoadName:
			stack.Push(env.Lookup(ins.Arg.ToString()))

		case op.CallFunc:
			argc := int(ins.Arg.(*gvalue.Number).Value)
			args := make([]gvalue.Instance, argc)

			for i := 0; i < argc; i++ {
				args[argc-i-1] = stack.Pop()
			}

			fun := stack.Pop()
			stack.Push(fun.(gvalue.Func).Call(args))

		default:
			assert.Unreached()
		}
	}

	if stack.Size() > 0 {
		return stack.Pop()
	}

	return nil
}
