package gazebo

import (
	"fmt"

	"github.com/johnfrankmorgan/gazebo/assert"
)

var gbuiltins = map[string]*GObject{
	"nil":   NewGObjectInferred(nil),
	"true":  NewGObjectInferred(true),
	"false": NewGObjectInferred(false),
	"printf": NewGObjectInferred(func(ctx *GFuncArgCtx) *GObject {
		ctx.ExpectsAtLeast(1)
		assert.True(ctx.Args[0].Type == gtypes.String)

		ifaces := make([]interface{}, len(ctx.Args)-1)
		for i, arg := range ctx.Args[1:] {
			ifaces[i] = arg.Value
		}

		fmt.Printf(ctx.Args[0].Value.(string), ifaces...)
		return NewGObjectInferred(nil)
	}),
}
