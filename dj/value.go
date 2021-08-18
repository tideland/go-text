// Tideland Go Text - Dynamic JSON
//
// Copyright (C) 2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dj // import "tideland.dev/go/text/dj"

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

//--------------------
// VALUE
//--------------------

// Value contains a document value which can be a string, int, float64, or bool.
// Based on the creation it also can be a structure or list and so allows to
// navigate deeper.
type Value struct {
	path []string
	data interface{}
	err  error
}

// newValue creates a values based on the passed data.
func newValue(path []string, data interface{}, err error) *Value {
	return &Value{
		path: path,
		data: data,
		err:  err,
	}
}

// Set allows to set a value to nil or one of string, int, float64, or boo.
func (v *Value) Set(data interface{}) {
	switch d := data.(type) {
	case string, int, float64, bool:
		v.data = d
		v.err = nil
		return
	default:
		if data == nil {
			v.data = nil
			v.err = nil
			return
		}
		v.data = nil
		v.err = &ValueError{
			Mode: "set",
			Err:  errors.New("invalid type"),
		}
	}
}

// IsUndefined returns true if the value contains no data.
func (v *Value) IsUndefined() bool {
	return v.data == nil
}

// Type returns the JSON type of the value.
func (v *Value) Type() NodeType {
	return nodeType(v.data)
}

// Len return 1 in case of simple value types of the number of
// elements in case of objects or arrays.
func (v *Value) Len() int {
	return nodeLen(v.data)
}

// AsString returns the value as string.
func (v *Value) AsString(dv string) string {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.data.(type) {
	case string:
		return tv
	case int:
		return strconv.Itoa(tv)
	case float64:
		return strconv.FormatFloat(tv, 'g', -1, 64)
	case bool:
		return strconv.FormatBool(tv)
	}
	return dv
}

// AsInt returns the value as int.
func (v *Value) AsInt(dv int) int {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.data.(type) {
	case string:
		i, err := strconv.Atoi(tv)
		if err != nil {
			return dv
		}
		return i
	case int:
		return tv
	case float64:
		return int(tv)
	case bool:
		if tv {
			return 1
		}
		return 0
	}
	return dv
}

// AsFloat64 returns the value as float64.
func (v *Value) AsFloat64(dv float64) float64 {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.data.(type) {
	case string:
		f, err := strconv.ParseFloat(tv, 64)
		if err != nil {
			return dv
		}
		return f
	case int:
		return float64(tv)
	case float64:
		return tv
	case bool:
		if tv {
			return 1.0
		}
		return 0.0
	}
	return dv
}

// AsBool returns the value as bool.
func (v *Value) AsBool(dv bool) bool {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.data.(type) {
	case string:
		b, err := strconv.ParseBool(tv)
		if err != nil {
			return dv
		}
		return b
	case int:
		return tv == 1
	case float64:
		return tv == 1.0
	case bool:
		return tv
	}
	return dv
}

// String implements fmt.Stringer.
func (v *Value) String() string {
	if v.IsUndefined() {
		return "null"
	}
	return fmt.Sprintf("%v", v.data)
}

// DeepEqual compares the value with the given one.
func (v *Value) DeepEqual(to *Value) bool {
	return reflect.DeepEqual(v.data, to.data)
}

// IsError provides the check for an error value.
func (v *Value) IsError() bool {
	return v.err != nil
}

// Error returns a potential error inside in the value.
func (v *Value) Error() error {
	return v.err
}

// At retrieves a value at a given path of keys.
func (v *Value) At(path ...string) *Value {
	jpath := append(v.path, path...)
	data, err := nodeAt(v.data, []string{}, path)
	if err != nil {
		return newValue(jpath, nil, err)
	}
	return newValue(jpath, data, nil)
}

// Do performs a function on all elements of the value
// if it is a node.
func (v *Value) Do(f func(key string, nv *Value) error) error {
	return nodeDo(v.path, v.data, f)
}

// EOF
