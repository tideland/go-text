// Tideland Go Text - Dynamic JSON
//
// Copyright (C) 2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dj // import "tideland.dev/go/text/dj"

import "strconv" //--------------------
// NODE HELPERS
//--------------------

// nodeLen returns the length of the passed node (which can be a single
// value too).
func nodeLen(node interface{}) int {
	if node == nil {
		return 0
	}
	switch n := node.(type) {
	case []interface{}:
		return len(n)
	case map[string]interface{}:
		return len(n)
	case string, int, float64, bool:
		return 1
	}
	return 0
}

// nodeDo performs a function on all elements of the passed node (which
// can be a single value too).
func nodeDo(node interface{}, f func(k string, v *Value) error) error {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case []interface{}:
		for i, d := range n {
			k := "#" + strconv.Itoa(i)
			if err := f(k, newValue(d)); err != nil {
				return err
			}
		}
		return nil
	case map[string]interface{}:
		for k, d := range n {
			if err := f(k, newValue(d)); err != nil {
				return err
			}
		}
		return nil
	case string, int, float64, bool:
		return f("", newValue(node))
	}
	return nil
}

// EOF
