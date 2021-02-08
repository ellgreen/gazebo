package gazebo

import (
	"fmt"
	"math"
	"reflect"
)

var gbuiltins = map[string]*GObject{
	"nil": NewGObjectInferred(nil),

	"true": NewGObjectInferred(true),

	"false": NewGObjectInferred(false),

	"printf": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
		var format string

		ctx.Parse(&format)

		fmt.Printf(format, ctx.Interfaces()[1:]...)
		return NewGObjectInferred(nil)
	}),

	"=": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
		ctx.ExpectsAtLeast(2)

		args := ctx.Interfaces()

		for _, arg := range args {
			if !reflect.DeepEqual(args[0], arg) {
				return NewGObjectInferred(false)
			}
		}

		return NewGObjectInferred(true)
	}),

	"?": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
		ctx.Expects(1)

		return NewGObjectInferred(ctx.Self().IsTruthy())
	}),

	"!": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
		ctx.Expects(1)

		return NewGObjectInferred(!ctx.Self().IsTruthy())
	}),

	">": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
		var (
			val1 float64
			val2 float64
		)

		ctx.Parse(&val1, &val2)

		return NewGObjectInferred(val1 > val2)
	}),

	"<": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
		var (
			val1 float64
			val2 float64
		)

		ctx.Parse(&val1, &val2)

		return NewGObjectInferred(val1 < val2)
	}),

	"+": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
		var (
			val1 float64
			val2 float64
		)

		ctx.Parse(&val1, &val2)

		return NewGObjectInferred(val1 + val2)
	}),

	"-": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
		var (
			val1 float64
			val2 float64
		)

		ctx.Parse(&val1, &val2)

		return NewGObjectInferred(val1 - val2)
	}),

	"*": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
		var (
			val1 float64
			val2 float64
		)

		ctx.Parse(&val1, &val2)

		return NewGObjectInferred(val1 * val2)

	}),

	"/": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
		var (
			val1 float64
			val2 float64
		)

		ctx.Parse(&val1, &val2)

		return NewGObjectInferred(val1 / val2)
	}),

	"%": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
		var (
			val1 float64
			val2 float64
		)

		ctx.Parse(&val1, &val2)

		return NewGObjectInferred(math.Mod(val1, val2))
	}),
}
