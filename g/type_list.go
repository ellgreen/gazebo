package g

func initlist() {
	TypeList = &Type{
		Name:   "List",
		Parent: TypeBase,
		Methods: Methods{
			Protocols.ToBool: Method(func(self Object, _ Args) Object {
				return NewObject(len(self.Value().([]Object)) > 0)
			}),

			Protocols.Len: Method(func(self Object, _ Args) Object {
				return NewObject(len(self.Value().([]Object)))
			}),

			Protocols.Index: Method(func(self Object, args Args) Object {
				index := ToInt(args.Self())
				return NewObject(self.Value().([]Object)[index].Value())
			}),
		},
	}
}
