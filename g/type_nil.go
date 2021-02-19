package g

func initnil() {
	TypeNil = &Type{
		Name:   "Nil",
		Parent: TypeBase,
		Methods: Methods{
			Protocols.ToBool: Method(func(_ Object, _ Args) Object {
				return NewObjectBool(false)
			}),
		},
	}
}
