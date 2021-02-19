package g

func initbool() {
	TypeBool = &Type{
		Name:   "Bool",
		Parent: TypeBase,
		Methods: Methods{
			Protocols.ToBool: Method(func(self Object, _ Args) Object {
				return NewObjectBool(EnsureBool(self).Bool())
			}),

			Protocols.ToNumber: Method(func(self Object, _ Args) Object {
				if EnsureBool(self).Bool() {
					return NewObjectNumber(1)
				}

				return NewObjectNumber(0)
			}),
		},
	}
}
