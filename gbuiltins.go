package gazebo

import (
	"fmt"
	"math"
	"reflect"

	"github.com/johnfrankmorgan/gazebo/assert"
)

var gbuiltins map[string]*GObject

func initbuiltins() {
	gbuiltins = map[string]*GObject{
		"nil": NewGObjectInferred(nil),

		"true": NewGObjectInferred(true),

		"false": NewGObjectInferred(false),

		"printf": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
			var format string

			ctx.Parse(&format)

			fmt.Printf(format, ctx.Interfaces()[1:]...)
			return NewGObjectInferred(nil)
		}),

		"println": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
			fmt.Println(ctx.Interfaces()...)
			return NewGObjectInferred(nil)
		}),

		"call": NewGObjectInferred(func(ctx *GFuncCtx) *GObject {
			var (
				method string
				self   *GObject
			)

			ctx.Parse(&method, &self)

			assert.True(self.Type.Implements(method))

			return self.Call(method, &GFuncCtx{
				VM:   ctx.VM,
				Args: ctx.Args[1:],
			})
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
}
