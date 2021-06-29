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
	"bytes"
	"errors"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/text/dj"
)

//--------------------
// VARIABLES
//--------------------

var emptyPath []string = []string{}

//--------------------
// TESTS
//--------------------

// TestValueAccess verifies valid and invalid access to values.
func TestValueAccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// String.
	v := dj.NewValue(emptyPath, "test", nil)
	assert.Equal(v.String(), "test")
	v = dj.NewValue(emptyPath, nil, nil)
	assert.Equal(v.String(), "")
	v = dj.NewValue(emptyPath, true, nil)
	assert.PanicsWith(func() {
		_ = v.String()
	}, "value is no string")

	// Int.
	v = dj.NewValue(emptyPath, 12345, nil)
	assert.Equal(v.Int(), 12345)
	v = dj.NewValue(emptyPath, nil, nil)
	assert.Equal(v.Int(), 0)
	v = dj.NewValue(emptyPath, true, nil)
	assert.PanicsWith(func() {
		v.Int()
	}, "value is no int")

	// Float64.
	v = dj.NewValue(emptyPath, 123.45, nil)
	assert.Equal(v.Float64(), 123.45)
	v = dj.NewValue(emptyPath, nil, nil)
	assert.Equal(v.Float64(), 0.0)
	v = dj.NewValue(emptyPath, true, nil)
	assert.PanicsWith(func() {
		v.Float64()
	}, "value is no float64")

	// Bool.
	v = dj.NewValue(emptyPath, true, nil)
	assert.Equal(v.Bool(), true)
	v = dj.NewValue(emptyPath, nil, nil)
	assert.Equal(v.Bool(), false)
	v = dj.NewValue(emptyPath, "true", nil)
	assert.PanicsWith(func() {
		v.Bool()
	}, "value is no bool")
}

// TestValueSetting verifies the setting of values.
func TestValueSetting(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	v := dj.NewValue(emptyPath, nil, nil)

	v.Set("test")
	assert.Equal(v.String(), "test")
	v.Set(12345)
	assert.Equal(v.Int(), 12345)
	v.Set(123.45)
	assert.Equal(v.Float64(), 123.45)
	v.Set(true)
	assert.Equal(v.Bool(), true)

	v.Set(struct{}{})
	assert.True(v.IsError())
	assert.ErrorContains(v.Error(), "invalid type")
}

// TestValueTests verifies testing of values.
func TestValueTests(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	v := dj.NewValue(emptyPath, nil, nil)
	assert.True(v.IsNil())

	v.Set("test")
	assert.False(v.IsNil())

	v = dj.NewValue(emptyPath, []interface{}{}, nil)
	assert.True(v.IsNode())
	v = dj.NewValue(emptyPath, map[string]interface{}{}, nil)
	assert.True(v.IsNode())

	v = dj.NewValue(emptyPath, map[string]int{}, nil)
	assert.False(v.IsNode())
}

// TestValueIteration verifies the iteration over value data.
func TestValueIteration(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	inA := []interface{}{1, 2, 3, 4, 5}
	outA := map[string]int{}
	v := dj.NewValue(emptyPath, inA, nil)
	err := v.Do(func(k string, fv *dj.Value) error {
		outA[k] = fv.Int()
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
	v = dj.NewValue(emptyPath, inB, nil)
	err = v.Do(func(k string, fv *dj.Value) error {
		outB[fv.String()] = k
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

	v = dj.NewValue(emptyPath, "test", nil)
	s := "a "
	err = v.Do(func(k string, fv *dj.Value) error {
		s += k + fv.String()
		return nil
	})
	assert.NoError(err)
	assert.Equal(s, "a test")

	err = dj.NewValue(emptyPath, false, nil).Do(func(k string, fv *dj.Value) error {
		return errors.New("ouch")
	})
	assert.ErrorContains(err, "ouch")
}

// TestValueAt verifies the navigation to a value of a value.
func TestValueAt(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	tests := []struct {
		name  string
		in    string
		path  []string
		value string
		err   string
	}{
		{
			"single string value",
			`{"root":["test"]}`,
			[]string{"#0"},
			"test",
			"",
		}, {
			"key/value document",
			`{"root":{"test": "12345"}}`,
			[]string{"test"},
			"12345",
			"",
		}, {
			"list document",
			`{"root":["1", "2", "3", "4", "5"]}`,
			[]string{"#2"},
			"3",
			"",
		}, {
			"nested document",
			`{"root":{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}}`,
			[]string{"o", "a", "#3"},
			"4",
			"",
		}, {
			"not existing path",
			`{"root":{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}}`,
			[]string{"o", "oops", "a"},
			"",
			"path does not exist",
		}, {
			"path too long",
			`{"root":{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}}`,
			[]string{"o", "x", "oops"},
			"",
			"path too long",
		}, {
			"invalid index / no number",
			`{"root":{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}}`,
			[]string{"o", "a", "oops"},
			"",
			"no index",
		}, {
			"invalid index / invalid number",
			`{"root":{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}}`,
			[]string{"o", "a", "#999"},
			"",
			"invalid array index",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			defer assert.SetFailable(t)()
			b := bytes.NewBufferString(test.in)
			doc, err := dj.Parse(b)
			assert.NoError(err, "in:", test.in)
			root := doc.At("root")
			value := root.At(test.path...)
			if test.err != "" {
				assert.ErrorContains(value.Error(), test.err)
			} else {
				assert.Equal(value.String(), test.value)
			}
		})
	}
}

// EOF
