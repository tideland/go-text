// Tideland Go Text - Generic JSON Processing - Value
//
// Copyright (C) 2019-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package gjp // import "tideland.dev/go/text/gjp"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"reflect"
	"strconv"
)

//--------------------
// VALUE
//--------------------

// Value contains one JSON value.
type Value struct {
	raw interface{}
	err error
}

// IsUndefined returns true if this value is undefined.
func (v *Value) IsUndefined() bool {
	return v.raw == nil
}

// AsString returns the value as string.
func (v *Value) AsString(dv string) string {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.raw.(type) {
	case string:
		return tv
	case int:
		return strconv.Itoa(tv)
	case float64:
		return strconv.FormatFloat(tv, 'f', -1, 64)
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
	switch tv := v.raw.(type) {
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
	switch tv := v.raw.(type) {
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
	switch tv := v.raw.(type) {
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

// Equals compares a value with the passed one.
func (v *Value) Equals(to *Value) bool {
	return reflect.DeepEqual(v.raw, to.raw)
}

// String implements fmt.Stringer.
func (v *Value) String() string {
	if v.IsUndefined() {
		return "null"
	}
	return fmt.Sprintf("%v", v.raw)
}

// EOF
