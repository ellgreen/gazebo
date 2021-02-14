package g

// IsTruthy determines if the provided Object is truthy
func IsTruthy(object Object) bool {
	return object.Call(Protocols.ToBool, nil).Value().(bool)
}

// ToString returns an Object's string value
func ToString(object Object) string {
	return object.Call(Protocols.ToString, nil).Value().(string)
}

// ToFloat returns an Object's float value
func ToFloat(object Object) float64 {
	return object.Call(Protocols.ToNumber, nil).Value().(float64)
}

// ToInt returns an Object's int value
func ToInt(object Object) int {
	return int(ToFloat(object))
}

// Invoke calls an Object
func Invoke(object Object, args Args) Object {
	return object.Call(Protocols.Invoke, args)
}
