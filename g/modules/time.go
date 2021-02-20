package modules

import (
	"time"

	"github.com/johnfrankmorgan/gazebo/errors"
	"github.com/johnfrankmorgan/gazebo/g"
)

// TypeTime is the type of Time objects
var TypeTime *g.Type

// ObjectTime is the underlying value of Time objects in gazebo
type ObjectTime struct {
	g.PartialObject
	value time.Time
}

// NewObjectTime creates a new ObjectTime
func NewObjectTime(value time.Time) *ObjectTime {
	object := &ObjectTime{value: value}
	object.SetType(TypeTime)
	return object
}

// Value satisfies the g.Object interface
func (m *ObjectTime) Value() interface{} {
	return m.value
}

// Call satisfies the g.Object interface
func (m *ObjectTime) Call(method string, args g.Args) g.Object {
	return m.CallMethod(m, method, args)
}

// Time returns the ObjectTime's value
func (m *ObjectTime) Time() time.Time {
	return m.value
}

// EnsureTime asserts a value is a time object
func EnsureTime(value g.Object) *ObjectTime {
	errors.ErrRuntime.Expect(
		value.Type() == TypeTime,
		"expected type Time got %s",
		value.Type().Name,
	)

	return value.(*ObjectTime)
}

// Time holds the definitions for the time module
var Time = &Module{
	Name: "time",
	Init: func() {
		TypeTime = &g.Type{
			Name:   "Time",
			Parent: g.TypeBase,
			Methods: g.Methods{
				g.Protocols.ToNumber: func(self g.Object, _ g.Args) g.Object {
					return g.NewObjectNumber(float64(EnsureTime(self).Time().Unix()))
				},

				g.Protocols.ToString: func(self g.Object, _ g.Args) g.Object {
					return g.NewObjectString(EnsureTime(self).Time().String())
				},
			},
		}
	},
	Values: map[string]g.Object{
		"new": g.NewObjectInternalFunc(func(args g.Args) g.Object {
			str := g.EnsureString(args.Self())

			value, err := time.Parse("2006-01-02 15:04:05", str.String())
			errors.ErrRuntime.ExpectNil(err, "%v", err)

			return NewObjectTime(value)
		}),

		"now": g.NewObjectInternalFunc(func(_ g.Args) g.Object {
			return NewObjectTime(time.Now())
		}),

		"sleep": g.NewObjectInternalFunc(func(args g.Args) g.Object {
			duration := time.Duration(g.EnsureNumber(args.Self()).Int())

			time.Sleep(duration * time.Millisecond)

			return g.NewObjectNil()
		}),
	},
}
