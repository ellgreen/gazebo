package g

func initlist() {
	TypeList = &Type{
		Name:   "List",
		Parent: TypeBase,
		Methods: Methods{
			Protocols.ToBool: Method(func(self Object, _ Args) Object {
				return NewObjectBool(EnsureList(self).Len() > 0)
			}),

			Protocols.Len: Method(func(self Object, _ Args) Object {
				return NewObjectNumber(float64(EnsureList(self).Len()))
			}),

			Protocols.Index: Method(func(self Object, args Args) Object {
				index := ToInt(args.Self())

				return EnsureList(self).Index(index)
			}),
		},
	}
}
