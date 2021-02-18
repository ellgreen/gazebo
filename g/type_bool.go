package g

func initbool() {
	TypeBool = &Type{
		Name:   "Bool",
		Parent: TypeBase,
		Methods: Methods{
			Protocols.ToBool: Method(func(self Object, _ Args) Object {
				return NewObject(self.Value())
			}),

			Protocols.ToNumber: Method(func(self Object, _ Args) Object {
				if self.Value().(bool) {
					return NewObject(1)
				}

				return NewObject(0)
			}),
		},
	}
}
