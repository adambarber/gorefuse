package gorefuse

import (
	"reflect"
)

// From builds a Because[R] from whatever the describable can say. From is the
// hardened chokepoint: untyped nil, typed nil (pointer/map/slice/func), or a
// Describe reporting no description all resolve to no description available.
func From[R any](reason R, d Describable) Because[R] {
	b := Because[R]{Reason: reason, Cause: nil, What: Refusal{File: "", Line: 0, Subject: "", Why: ""}}

	if isNilDescribable(d) {
		return b
	}

	if what, ok := d.Describe(); ok {
		b.What = what
	}

	return b
}

// isNilDescribable reports whether the describable holds nothing to call:
// the interface itself nil, or a nil pointer/map/slice/func inside it.
func isNilDescribable(d Describable) bool {
	if d == nil {
		return true
	}

	v := reflect.ValueOf(d)
	switch v.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func:
		return v.IsNil()
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan,
		reflect.Interface, reflect.String, reflect.Struct, reflect.UnsafePointer:
		return false
	default:
		return false
	}
}
