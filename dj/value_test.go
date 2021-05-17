// Tideland Go Text - Dynamic JSON - Testing
//
// Copyright (C) 2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dj_test // import "tideland.dev/go/text/dj"

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/text/dj"
)

//--------------------
// TESTS
//--------------------

// TestValueAccess verifies valid and invalid access to values.
func TestValueAccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// String.
	v := dj.NewValue("test")
	assert.Equal(v.AsString(), "test")
	v = dj.NewValue(nil)
	assert.Equal(v.AsString(), "")
	v = dj.NewValue(true)
	assert.PanicsWith(func() {
		v.AsString()
	}, "value is no string")

	// Int.
	v = dj.NewValue(12345)
	assert.Equal(v.AsInt(), 12345)
	v = dj.NewValue(nil)
	assert.Equal(v.AsInt(), 0)
	v = dj.NewValue(true)
	assert.PanicsWith(func() {
		v.AsInt()
	}, "value is no int")

	// Float64.
	v = dj.NewValue(123.45)
	assert.Equal(v.AsFloat64(), 123.45)
	v = dj.NewValue(nil)
	assert.Equal(v.AsFloat64(), 0.0)
	v = dj.NewValue(true)
	assert.PanicsWith(func() {
		v.AsFloat64()
	}, "value is no float64")

	// Bool.
	v = dj.NewValue(true)
	assert.Equal(v.AsBool(), true)
	v = dj.NewValue(nil)
	assert.Equal(v.AsBool(), false)
	v = dj.NewValue("true")
	assert.PanicsWith(func() {
		v.AsBool()
	}, "value is no bool")
}

// TestValueSetting verifies the setting of values.
func TestValueSetting(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	v := dj.NewValue(nil)

	v.Set("test")
	assert.Equal(v.AsString(), "test")
	v.Set(12345)
	assert.Equal(v.AsInt(), 12345)
	v.Set(123.45)
	assert.Equal(v.AsFloat64(), 123.45)
	v.Set(true)
	assert.Equal(v.AsBool(), true)

	assert.PanicsWith(func() {
		v.Set(struct{}{})
	}, "invalid type for value setting")
}

// TestValueTests verifies testing of values.
func TestValueTests(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	v := dj.NewValue(nil)
	assert.True(v.IsNil())

	v.Set("test")
	assert.False(v.IsNil())

	v = dj.NewValue([]interface{}{})
	assert.True(v.IsNode())
	v = dj.NewValue(map[string]interface{}{})
	assert.True(v.IsNode())

	v = dj.NewValue(map[string]int{})
	assert.False(v.IsNode())
}

// TestValueIteration verifies the iteration over value data.
func TestValueIteration(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	inA := []interface{}{1, 2, 3, 4, 5}
	outA := map[string]int{}
	v := dj.NewValue(inA)
	err := v.Do(func(k string, fv *dj.Value) error {
		outA[k] = fv.AsInt()
		return nil
	})
	assert.NoError(err)
	assert.Equal(outA, map[string]int{
		"#0": 1,
		"#1": 2,
		"#2": 3,
		"#3": 4,
		"#4": 5,
	})

	inB := map[string]interface{}{
		"a": "one",
		"b": "two",
		"c": "three",
		"d": "four",
		"e": "five",
	}
	outB := map[string]string{}
	v = dj.NewValue(inB)
	err = v.Do(func(k string, fv *dj.Value) error {
		outB[fv.AsString()] = k
		return nil
	})
	assert.NoError(err)
	assert.Equal(outB, map[string]string{
		"one":   "a",
		"two":   "b",
		"three": "c",
		"four":  "d",
		"five":  "e",
	})

	v = dj.NewValue("test")
	s := "a "
	err = v.Do(func(k string, fv *dj.Value) error {
		s += k + fv.AsString()
		return nil
	})
	assert.NoError(err)
	assert.Equal(s, "a test")

	err = dj.NewValue(false).Do(func(k string, fv *dj.Value) error {
		return errors.New("ouch")
	})
	assert.ErrorContains(err, "ouch")
}

// EOF
