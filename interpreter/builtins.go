package interpreter

import (
	"fmt"

	"github.com/johnfrankmorgan/gazebo/assert"
	"github.com/johnfrankmorgan/gazebo/gvalue"
	"github.com/kr/pretty"
)

var builtins = map[string]gvalue.Instance{
	"nil":   gvalue.New(nil),
	"true":  gvalue.New(true),
	"false": gvalue.New(false),
	"!": gvalue.Builtin("!", func(args []gvalue.Instance) gvalue.Instance {
		assert.True(len(args) == 1)

		return gvalue.New(!args[0].ToBool())
	}),
	"+": gvalue.Builtin("+", func(args []gvalue.Instance) gvalue.Instance {
		var sum float64

		for _, arg := range args {
			sum += arg.(*gvalue.Number).Value
		}

		return gvalue.New(sum)
	}),
	"printf": gvalue.Builtin("printf", func(args []gvalue.Instance) gvalue.Instance {
		assert.True(len(args) > 0)

		ifaces := make([]interface{}, len(args)-1)
		for i, arg := range args[1:] {
			ifaces[i] = arg.Interface()
		}

		fmt.Printf(args[0].ToString(), ifaces...)
		return gvalue.New(nil)
	}),
	"println": gvalue.Builtin("println", func(args []gvalue.Instance) gvalue.Instance {
		ifaces := make([]interface{}, len(args))
		for i, arg := range args {
			ifaces[i] = arg.Interface()
		}

		fmt.Println(ifaces...)
		return gvalue.New(nil)
	}),
	"debugln": gvalue.Builtin("debugln", func(args []gvalue.Instance) gvalue.Instance {
		ifaces := make([]interface{}, len(args))
		for i, arg := range args {
			ifaces[i] = arg
		}

		pretty.Println(ifaces...)
		return gvalue.New(nil)
	}),
}

func LoadBuiltins(env *Env) {
	for name, value := range builtins {
		env.Define(name, value)
	}
}
