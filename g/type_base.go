package g

import (
	"fmt"
	"reflect"
)

func initbase() {
	TypeBase = &Type{
		Name:   "Base",
		Parent: nil,
		Methods: Methods{
			Protocols.ToBool: Method(func(_ Object, _ Args) Object {
				return NewObjectBool(true)
			}),

			Protocols.ToString: Method(func(self Object, _ Args) Object {
				return NewObjectString(fmt.Sprintf("%v", self.Value()))
			}),

			Protocols.ToNumber: Method(func(self Object, _ Args) Object {
				return NewObjectNumber(0)
			}),

			Protocols.Inspect: Method(func(self Object, _ Args) Object {
				inspection := fmt.Sprintf(
					"<gtypes.%s %p>(%v %p)",
					self.Type().Name,
					self.Type(),
					self.Value(),
					self,
				)

				return NewObjectString(inspection)
			}),

			Protocols.Equal: Method(func(self Object, args Args) Object {
				for _, arg := range args {
					if !reflect.DeepEqual(self.Value(), arg.Value()) {
						return NewObjectBool(false)
					}
				}

				return NewObjectBool(true)
			}),

			Protocols.HasAttr: Method(func(self Object, args Args) Object {
				name := EnsureString(args.Self()).String()
				return NewObjectBool(self.Attributes().Has(name))
			}),

			Protocols.GetAttr: Method(func(self Object, args Args) Object {
				name := EnsureString(args.Self()).String()
				return self.Attributes().Get(name)
			}),

			Protocols.SetAttr: Method(func(self Object, args Args) Object {
				args.Expects(2)

				name := EnsureString(args.Self()).String()

				self.Attributes().Set(name, args[1])

				return NewObjectNil()
			}),

			Protocols.DelAttr: Method(func(self Object, args Args) Object {
				name := EnsureString(args.Self()).String()

				self.Attributes().Delete(name)

				return NewObjectNil()
			}),
		},
	}
}
