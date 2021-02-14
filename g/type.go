package g

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/johnfrankmorgan/gazebo/assert"
)

// Type represents the type of gazebo values
type Type struct {
	Name    string
	Parent  *Type
	Methods Methods
}

// Method is the type of methods in gazebo
type Method func(Object, Args) Object

// Methods is a map of names to Funcs
type Methods map[string]Method

// Resolve resolves a method on a Type
func (m *Type) Resolve(name string) Method {
	if method, ok := m.Methods[name]; ok {
		return method
	}

	if m.Parent != nil {
		return m.Parent.Resolve(name)
	}

	return nil
}

// Implements checks if a method is implemented by a Type
func (m *Type) Implements(name string) bool {
	return m.Resolve(name) != nil
}

// Builtin types
var (
	TypeBase        *Type
	TypeNil         *Type
	TypeBool        *Type
	TypeNumber      *Type
	TypeString      *Type
	TypeList        *Type
	TypeBuiltinFunc *Type
	TypeFunc        *Type
	TypeInternal    *Type
)

func init() {
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

	TypeNil = &Type{
		Name:   "Nil",
		Parent: TypeBase,
		Methods: Methods{
			Protocols.ToBool: Method(func(_ Object, _ Args) Object {
				return NewObject(false)
			}),
		},
	}

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

	TypeString = &Type{
		Name:   "String",
		Parent: TypeBase,
		Methods: Methods{
			Protocols.ToBool: Method(func(self Object, _ Args) Object {
				return NewObject(self.Value().(string) != "")
			}),

			Protocols.ToNumber: Method(func(self Object, _ Args) Object {
				value, err := strconv.ParseFloat(self.Value().(string), 64)
				assert.Nil(err)

				return NewObject(value)
			}),

			Protocols.Len: Method(func(self Object, _ Args) Object {
				return NewObject(len(self.Value().(string)))
			}),

			Protocols.Index: Method(func(self Object, args Args) Object {
				index := ToInt(args.Self())
				return NewObject(self.Value().(string)[index : index+1])
			}),
		},
	}

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

	TypeBuiltinFunc = &Type{
		Name:   "BuiltinFunc",
		Parent: TypeBase,
		Methods: Methods{
			Protocols.Invoke: Method(func(self Object, args Args) Object {
				return self.Value().(Func)(args)
			}),
		},
	}

	TypeFunc = &Type{
		Name:    "Func",
		Parent:  TypeBase,
		Methods: Methods{},
	}

	TypeInternal = &Type{
		Name:    "Internal",
		Parent:  TypeBase,
		Methods: Methods{},
	}
}
