package g

import (
	"strconv"

	"github.com/johnfrankmorgan/gazebo/errors"
)

func initstring() {
	TypeString = &Type{
		Name:   "String",
		Parent: TypeBase,
		Methods: Methods{
			Protocols.ToBool: Method(func(self Object, _ Args) Object {
				return NewObjectBool(EnsureString(self).Len() > 0)
			}),

			Protocols.ToNumber: Method(func(self Object, _ Args) Object {
				value, err := strconv.ParseFloat(EnsureString(self).String(), 64)
				errors.ErrRuntime.ExpectNil(err, err.Error())

				return NewObjectNumber(value)
			}),

			Protocols.Len: Method(func(self Object, _ Args) Object {
				return NewObjectNumber(float64(EnsureString(self).Len()))
			}),

			Protocols.Index: Method(func(self Object, args Args) Object {
				index := ToInt(args.Self())
				return NewObjectString(EnsureString(self).String()[index : index+1])
			}),
		},
	}
}
