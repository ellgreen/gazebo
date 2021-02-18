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
				return NewObject(true)
			}),

			Protocols.ToString: Method(func(self Object, _ Args) Object {
				return NewObject(fmt.Sprintf("%v", self.Value()))
			}),

			Protocols.ToNumber: Method(func(self Object, _ Args) Object {
				return NewObject(0)
			}),

			Protocols.Inspect: Method(func(self Object, _ Args) Object {
				inspection := fmt.Sprintf(
					"<gtypes.%s %p>(%v %p)",
					self.Type().Name,
					self.Type(),
					self.Value(),
					self,
				)

				return NewObject(inspection)
			}),

			Protocols.Equal: Method(func(self Object, args Args) Object {
				for _, arg := range args {
					if !reflect.DeepEqual(self.Value(), arg.Value()) {
						return NewObject(false)
					}
				}

				return NewObject(true)
			}),

			Protocols.HasAttr: Method(func(self Object, args Args) Object {
				var name string

				args.Parse(&name)

				return NewObject(self.Attributes().Has(name))
			}),

			Protocols.GetAttr: Method(func(self Object, args Args) Object {
				var name string

				args.Parse(&name)

				return self.Attributes().Get(name)
			}),

			Protocols.SetAttr: Method(func(self Object, args Args) Object {
				var (
					name  string
					value Object
				)

				args.Parse(&name, &value)
				self.Attributes().Set(name, value)

				return NewObject(nil)
			}),

			Protocols.DelAttr: Method(func(self Object, args Args) Object {
				var name string

				args.Parse(&name)
				self.Attributes().Delete(name)

				return NewObject(nil)
			}),
		},
	}
}
