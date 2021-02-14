package modules

import (
	"strings"

	"github.com/johnfrankmorgan/gazebo/g"
)

// Str holds the definitions for the str module
var Str = &Module{
	Name: "str",
	Values: map[string]g.Object{
		"contains": g.NewObject(func(args g.Args) g.Object {
			args.Expects(2)

			source := g.ToString(args[0])
			substr := g.ToString(args[1])

			return g.NewObject(strings.Contains(source, substr))
		}),
	},
}
