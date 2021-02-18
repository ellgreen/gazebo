package g

import (
	"strconv"

	"github.com/johnfrankmorgan/gazebo/assert"
)

func initstring() {
	TypeString = &Type{
		Name:   "String",
		Parent: TypeBase,
		Methods: Methods{
			Protocols.ToBool: Method(func(self Object, _ Args) Object {
				return NewObject(self.Value().(string) != "")
			}),

			Protocols.ToNumber: Method(func(self Object, _ Args) Object {
				value, err := strconv.ParseFloat(self.Value().(string), 64)
				assert.Nil(err)

				return NewObject(value)
			}),

			Protocols.Len: Method(func(self Object, _ Args) Object {
				return NewObject(len(self.Value().(string)))
			}),

			Protocols.Index: Method(func(self Object, args Args) Object {
				index := ToInt(args.Self())
				return NewObject(self.Value().(string)[index : index+1])
			}),
		},
	}
}
