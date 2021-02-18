package g

import "github.com/johnfrankmorgan/gazebo/assert"

// EnsureNil asserts that an Object is an ObjectNil
func EnsureNil(value Object) *ObjectNil {
	assert.True(value.Type() == TypeNil, "expected type Nil got %s", value.Type().Name)
	return value.(*ObjectNil)
}

// EnsureBool asserts that an Object is an ObjectBool
func EnsureBool(value Object) *ObjectBool {
	assert.True(value.Type() == TypeBool, "expected type Bool got %s", value.Type().Name)
	return value.(*ObjectBool)
}

// EnsureNumber asserts that an Object is an ObjectNumber
func EnsureNumber(value Object) *ObjectNumber {
	assert.True(value.Type() == TypeNumber, "expected type Number got %s", value.Type().Name)
	return value.(*ObjectNumber)
}

// EnsureString asserts that an Object is an ObjectString
func EnsureString(value Object) *ObjectString {
	assert.True(value.Type() == TypeString, "expected type String got %s", value.Type().Name)
	return value.(*ObjectString)
}

// EnsureList asserts that an Object is an ObjectList
func EnsureList(value Object) *ObjectList {
	assert.True(value.Type() == TypeList, "expected type List got %s", value.Type().Name)
	return value.(*ObjectList)
}

// EnsureInternalFunc asserts that an Object is an ObjectInternalFunc
func EnsureInternalFunc(value Object) *ObjectInternalFunc {
	assert.True(
		value.Type() == TypeInternalFunc,
		"expected type InternalFunc got %s",
		value.Type().Name,
	)

	return value.(*ObjectInternalFunc)
}

// IsTruthy determines if the provided Object is truthy
func IsTruthy(object Object) bool {
	return EnsureBool(object.Call(Protocols.ToBool, nil)).Bool()
}

// ToString returns an Object's string value
func ToString(object Object) string {
	return EnsureString(object.Call(Protocols.ToString, nil)).String()
}

// ToFloat returns an Object's float value
func ToFloat(object Object) float64 {
	return EnsureNumber(object.Call(Protocols.ToNumber, nil)).Float()
}

// ToInt returns an Object's int value
func ToInt(object Object) int {
	return EnsureNumber(object.Call(Protocols.ToNumber, nil)).Int()
}

// Invoke calls an Object
func Invoke(object Object, args Args) Object {
	return object.Call(Protocols.Invoke, args)
}
