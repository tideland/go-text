// Tideland Go Text - Dynamic JSON
//
// Copyright (C) 2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dj // import "tideland.dev/go/text/dj"

//--------------------
// VALUE
//--------------------

// Value contains a document value which can be a string, int, float64, or bool.
// Based on the creation it also can be a structure or list and so allows to
// navigate deeper.
type Value struct {
	data interface{}
}

// newValue creates a values based on the passed data.
func newValue(data interface{}) *Value {
	return &Value{
		data: data,
	}
}

// Set allows to set a value to nil or one of string, int, float64, or boo.
func (v *Value) Set(data interface{}) {
	switch d := data.(type) {
	case string, int, float64, bool:
		v.data = d
		return
	default:
		if data == nil {
			v.data = nil
			return
		}
	}
	panic("invalid type for value setting")
}

// AsString returns a potential string value or panics.
func (v *Value) AsString() string {
	if v.data == nil {
		return ""
	}
	s, ok := v.data.(string)
	if !ok {
		panic("value is no string")
	}
	return s
}

// AsInt returns a potential int value or panics.
func (v *Value) AsInt() int {
	if v.data == nil {
		return 0
	}
	i, ok := v.data.(int)
	if !ok {
		panic("value is no int")
	}
	return i
}

// AsFloat64 returns a potential float64 value or panics.
func (v *Value) AsFloat64() float64 {
	if v.data == nil {
		return 0.0
	}
	f, ok := v.data.(float64)
	if !ok {
		panic("value is no float64")
	}
	return f
}

// AsBool returns a potential bool value or panics.
func (v *Value) AsBool() bool {
	if v.data == nil {
		return false
	}
	b, ok := v.data.(bool)
	if !ok {
		panic("value is no bool")
	}
	return b
}

// IsNil provides the check for a nil value.
func (v *Value) IsNil() bool {
	return v.data == nil
}

// IsNode provides the check for structures and lists.
func (v *Value) IsNode() bool {
	switch v.data.(type) {
	case []interface{}, map[string]interface{}:
		return true
	default:
		return false
	}
}

// Do performs a function on all elements of the value
// if it is a node.
func (v *Value) Do(f func(key string, nv *Value) error) error {
	return nodeDo(v.data, f)
}

// EOF