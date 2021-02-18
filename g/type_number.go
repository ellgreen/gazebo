package g

func initnumber() {
	TypeNumber = &Type{
		Name:   "Number",
		Parent: TypeBase,
		Methods: Methods{
			Protocols.ToBool: Method(func(self Object, _ Args) Object {
				return NewObject(self.Value().(float64) != 0)
			}),

			Protocols.ToNumber: Method(func(self Object, _ Args) Object {
				return NewObject(self.Value())
			}),

			Protocols.Add: Method(func(self Object, args Args) Object {
				result := self.Value().(float64)

				for _, arg := range args {
					result += ToFloat(arg)
				}

				return NewObject(result)
			}),

			Protocols.Sub: Method(func(self Object, args Args) Object {
				result := self.Value().(float64) - ToFloat(args.Self())
				return NewObject(result)
			}),

			Protocols.Mul: Method(func(self Object, args Args) Object {
				result := self.Value().(float64) * ToFloat(args.Self())
				return NewObject(result)
			}),

			Protocols.Div: Method(func(self Object, args Args) Object {
				result := self.Value().(float64) / ToFloat(args.Self())
				return NewObject(result)
			}),

			Protocols.LessThan: Method(func(self Object, args Args) Object {
				result := self.Value().(float64) < ToFloat(args.Self())
				return NewObject(result)
			}),

			Protocols.GreaterThan: Method(func(self Object, args Args) Object {
				result := self.Value().(float64) > ToFloat(args.Self())
				return NewObject(result)
			}),
		},
	}
}
