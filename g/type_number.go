package g

func initnumber() {
	TypeNumber = &Type{
		Name:   "Number",
		Parent: TypeBase,
		Methods: Methods{
			Protocols.ToBool: Method(func(self Object, _ Args) Object {
				return NewObjectBool(EnsureNumber(self).Float() != 0)
			}),

			Protocols.ToNumber: Method(func(self Object, _ Args) Object {
				return NewObjectNumber(EnsureNumber(self).Float())
			}),

			Protocols.Add: Method(func(self Object, args Args) Object {
				result := EnsureNumber(self).Float()

				for _, arg := range args {
					result += ToFloat(arg)
				}

				return NewObjectNumber(result)
			}),

			Protocols.Sub: Method(func(self Object, args Args) Object {
				result := EnsureNumber(self).Float() - ToFloat(args.Self())
				return NewObjectNumber(result)
			}),

			Protocols.Mul: Method(func(self Object, args Args) Object {
				result := EnsureNumber(self).Float() * ToFloat(args.Self())
				return NewObjectNumber(result)
			}),

			Protocols.Div: Method(func(self Object, args Args) Object {
				result := EnsureNumber(self).Float() / ToFloat(args.Self())
				return NewObjectNumber(result)
			}),

			Protocols.LessThan: Method(func(self Object, args Args) Object {
				result := EnsureNumber(self).Float() < ToFloat(args.Self())
				return NewObjectBool(result)
			}),

			Protocols.GreaterThan: Method(func(self Object, args Args) Object {
				result := EnsureNumber(self).Float() > ToFloat(args.Self())
				return NewObjectBool(result)
			}),
		},
	}
}
