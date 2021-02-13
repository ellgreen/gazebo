package g

import (
	"fmt"
	"math"
	"reflect"
)

func _protocolmethods() []string {
	var methods []string

	value := reflect.ValueOf(Protocols)

	for i := 0; i < value.NumField(); i++ {
		methods = append(methods, value.Field(i).String())
	}

	return methods
}

// Builtins returns all of the builtin objects in gazebo
func Builtins() map[string]Object {
	methodcall := func(name string) Object {
		return NewObject(func(args Args) Object {
			self, args := args.SelfWithArgs()
			return self.Call(name, args)
		})
	}

	wrapmethods := func(builtins map[string]Object) map[string]Object {
		for _, method := range _protocolmethods() {
			builtins[method] = methodcall(method)
		}

		return builtins
	}

	return wrapmethods(map[string]Object{
		"nil": NewObject(nil),

		"false": NewObject(false),

		"true": NewObject(true),

		"!": NewObject(func(args Args) Object {
			return NewObject(!IsTruthy(args.Self()))
		}),

		"%": NewObject(func(args Args) Object {
			args.Expects(2)

			result := math.Mod(ToFloat(args[0]), ToFloat(args[1]))

			return NewObject(result)
		}),

		"call": NewObject(func(args Args) Object {
			name, args := args.SelfWithArgs()
			self, args := args.SelfWithArgs()
			return self.Call(name.Value().(string), args)
		}),

		"nil?": NewObject(func(args Args) Object {
			return NewObject(args.Self().Type() == TypeNil)
		}),

		"println": NewObject(func(args Args) Object {
			fmt.Println(args.Values()...)
			return NewObject(nil)
		}),

		"printf": NewObject(func(args Args) Object {
			format, args := args.SelfWithArgs()
			fmt.Printf(format.Value().(string), args.Values()...)
			return NewObject(nil)
		}),
	})
}
